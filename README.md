# alphaess-to-mqtt

[**AlphaESS**](https://www.alphaess.com) are a provider of Solar/battery systems. System information is sent to cloud servers enabling their mobile Apps.
This project utilises a proxy library to intercept that data and inject it into a [Home Assistant](https://www.home-assistant.io/)(HA) instance (via [MQTT discovery](https://www.home-assistant.io/docs/mqtt/discovery/)).
The initial implementation is read only, and simply mirrors the data to HA allowing the alphaess.com mobile apps continue as normal.
Future implementations could modify the stream - allowing for some control via HA, or act as a server and block all data from leaving the network.

# Overview


**Default setup without proxy:**

    |---- ---- ---- Home ---- ---- ----|
    |(Solar Panels >) AlphaESS > Router| -> Internet -> | alphaess.com |

**1. Modify destintation in config**

    |---- ---- ---- Home ---- ---- ----|
    |AlphaESS > new proxy mirrors to   |
    |             1. > MQTT > HA       |
    |             2. > Router          | -> Internet -> | alphaess.com |   

**2. Router Redirect**

    |---- ---- ---- Home ---- ---- ----|
    |AlphaESS > Router NAT redirect to |
    |             proxy which mirrors: |
    |             1. > MQTT > HA       |
    |             2. > Router          | -> Internet -> | alphaess.com |   

# Crude Installation Notes
1. Configure proxy instance via alphaESS-proxy.conf & optionally install as service
2. Configure AlphaESS system or router to ensure traffic is routed to proxy
3. Configure HA to convert values to kWh so they can be used in the [Energy dashboard](https://www.home-assistant.io/blog/2021/08/04/home-energy-management/). In configuraton.yaml under (or per your config).

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

4. Optionally configure Daily Meters to compare to the Alpha App:

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

