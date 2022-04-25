# alphaess-to-mqtt

[**AlphaESS**](https://www.alphaess.com) are a provider of Solar/battery systems. System information is sent to cloud servers enabling their mobile Apps.
This project utilises a proxy [library](https://github.com/230delphi/go-any-proxy) (originally built by [Ryan Chapman](http://blog.rchapman.org/posts/Transparently_proxying_any_tcp_connection/)) to intercept that data and send it to a [Home Assistant](https://www.home-assistant.io/) (HA) instance (via [MQTT discovery](https://www.home-assistant.io/docs/mqtt/discovery/)).

There are currently 2 implementations:
* Read only - simply mirrors the data (sent to the server every 10 seconds) to HA allowing the alphaess.com mobile apps continue as normal.
* Read & Inject - mirrors the same data to HA, but can also modify and/or inject requests into the stream. Currently only enable/disable charging from grid is supported

Current implementation could support further modify the stream - enabling other changes or simply filtering out private data normally sent to the cloud service.
Further work could simply fake the server and disconnect from the internet entirely.

The goal of this implementation is to gather all data from the system and eventually lead to independence from the cloud systems. For many, Charles Gillanders [API polling](https://github.com/CharlesGillanders/alphaess) implementation is more suitable - it is easier to configure and less susceptible breaking due to minor changes by the provider.

**Note:** This is the core library implementation. Simplest installation is via the [HA Addon](https://github.com/230delphi/hassio-addons/tree/main/alphaess-proxy-addon).
Example use can be seen [here](example_use/README.md).

# Overview


**Default setup without proxy:**

    |---- ---- ---- Home ---- ---- ----|
    |(Solar Panels >) AlphaESS > Router| -> Internet -> | alphaess.com |

###1. Direct Proxy
Modify AlphaESS configuration to direct traffic to proxy.

    |---- ---- ---- Home ---- ---- ----|
    |AlphaESS > new proxy mirrors to   |
    |             1. > MQTT > HA       |
    |             2. > Router          | -> Internet -> | alphaess.com |   

**2. Transparent Proxy/Router Redirect**
Modify network routing rules to send traffic to the proxy, which will forward the traffic to the intended destination (AlphaESS cloud).

    |---- ---- ---- Home ---- ---- ----|
    |AlphaESS > Router NAT redirect to |
    |             proxy which mirrors: |
    |             1. > MQTT > HA       |
    |             2. > Router          | -> Internet -> | alphaess.com |   

# Crude Installation Steps
1. Configure proxy instance via alphaESS-proxy.conf & optionally install as service
2. Configure AlphaESS system or router to ensure traffic is routed to proxy
3. Configure HA to convert values to kWh so they can be used in the [Energy dashboard](https://www.home-assistant.io/blog/2021/08/04/home-energy-management/). In configuraton.yaml under (or per your config).

### 1. Configure Proxy Instance
See alphaESS-proxy.conf - For typical deployments as a transparent proxy only the MQTT details need to be configured.
The proxy is build on the work of [Ryan Chapman](http://blog.rchapman.org/posts/Transparently_proxying_any_tcp_connection/) Further details on configuration options can be found on his site.
```javascript
##### Proxy configuration options
#listening ip/port
l=0.0.0.0:7777

#To Enable the Direct proxy you need to configure the AlphaESS destination
#p=52.230.104.147:7777

## Additional logging
v=0
stat=0

#### MQTT instance details
MQTTAddress=<tcp://127.0.0.1:1883>
MQTTUser=<username>
MQTTPassword=<password>
    
## Other MQTT details - defaults should be fine
MQTTSendTimeout=5
MQTTTopicBase=homeassistant/sensor/
AlphaESSID=alphaess1
    
#### Other 
## Optional additional logging
#MSGLogging=GenericRQ,CommandIndexRQ,CommandRQ,ConfigRS,StatusRQ
    
## Alphaess servers are generally set to UTC but use local time for schedule configuration. this setting ensures the correct local time.
TZLocation=Europe/Dublin
    
## There are currently 2 implementations; Read only or Read & inject
# proxyConnection=MQTTReadProxyConnection
proxyConnection=MQTTInjectProxyConnection
```
### 2. Configure Network Routing to Proxy
This depends on your setup.
* **Direct Proxy:** Configure your AlphaESS system with a new destination IP - your proxy server
* **Transparent Proxy:** Configure your gateway to forward traffic bound for AlphaESS Cloud (52.230.104.147:7777) to go to your proxy.

### 3. Configure Home Assistant Values for Energy Dashboard
```javascript
    sensor:
     - platform: integration
        source: sensor.alphaess_totalsolar
        name: int_alphaess_totalsolar
        unit_prefix: k
        round: 3
        method: left
      - platform: integration
        source: sensor.alphaess_grid_consumption
        name: int_alphaess_grid_consumption
        unit_prefix: k
        round: 3
        method: left    
      - platform: integration
        source: sensor.alphaess_grid_return
        name: int_alphaess_grid_return
        unit_prefix: k
        round: 3
        method: left
      - platform: integration
        source: sensor.alphaess_battery_load
        name: int_alphaess_battery_load
        unit_prefix: k
        round: 3
        method: left
      - platform: integration
        source: sensor.alphaess_battery_charge
        name: int_alphaess_battery_charge
        unit_prefix: k
        round: 3
        method: left
      - platform: integration
        source: sensor.alphaess_total_load
        name: int_alphaess_total_load
        unit_prefix: k
        round: 3
        method: left
```

### 4. Optionally configure Daily Meters to compare to the Alpha App:

```javascript
utility_meter:
  daily_alphaess_totalsolar:
    source: sensor.int_alphaess_totalsolar
    cycle: daily
  daily_alphaess_grid_consumption:
    source: sensor.int_alphaess_grid_consumption
    cycle: daily
  daily_alphaess_grid_return:
    source: sensor.int_alphaess_grid_return
    cycle: daily
  daily_alphaess_battery_load:
    source: sensor.int_alphaess_battery_load
    cycle: daily
  daily_alphaess_battery_charge:
    source: sensor.int_alphaess_battery_charge
    cycle: daily
  daily_alphaess_total_load:
    source: sensor.int_alphaess_total_load
    cycle: daily
 ```

### 5. Example use in Home Assistant Script
Service can be configured with a JSON Payload with 4 optional values as below. 
Any *non-empty* or *corrupt / invalid json* payload will toggle on last known charging state with default values.
* ```"GridCharge":true,``` sets the config to charge from Grid or not. - Default: Toggle last known state; proxy start state is false. (ie: first request will enable) 
* ```"StartHour":1,``` sets the starting hour for charging. if omitted, or <0 it will execute now. (charging if ```GridCharge``` is true.) - default "-1" (now)
* ```"MinimumDuration":60,``` Sets the minimum duration for charge. as scheduling is based on hour, this simply considers cross hour and midnight boundaries.
 For example: enable now at 6.45pm for 10 mins, will enable until 7pm. set at 30 mins, it will enable until 8pm. - default 10 minutes
* ```"BatHighCap":40``` Sets the high level charging point; alphaESS will stop charging at this point (but will not revert to battery until end of the schedule) - default 50%

 ```javascript
# Note: Configuration is based on schedule
#  charging will be re-enabled on the following day unless explicity disabled or configuration is reset from cloud.
#  these schedules will not be reflected in the could app.
    
# schedule the battery to run immediatly
enable_battery_charge_now:
  alias: Enable Battery Charging now
  sequence:
    - service: mqtt.publish
      data_template:
        topic: "homeassistant/sensor/alphaess1/action/chargebattery"
        payload: >
          {
            "GridCharge":"true",
            "MinimumDuration":5,
            "BatHighCap":40
          }
          
# this will also disable any schedule
disable_battery_charge_now:
  alias: Disable Battery Charge
  sequence:
    - service: mqtt.publish
      data_template:
        topic: "homeassistant/sensor/alphaess1/action/chargebattery"
        payload: >
          {
            "GridCharge":"false"
          }
          
# this will effectivly disable any current charging and reset any sechdule to 1am
# note, setting to MinimumDuration "60" will actually schedule for 2 hours (59min will set for 1 hour); 
# for finer granularity just enable/disable the schedule on demand from HA. 
enable_battery_charge_at1:
  alias: Enable Battery Charge At 1am
  sequence:
    - service: mqtt.publish
      data_template:
        topic: "homeassistant/sensor/alphaess1/action/chargebattery"
        payload: >
          {
            "GridCharge":"true",
            "StartHour":1,
            "MinimumDuration":60,
            "BatHighCap":40
          }
 ```