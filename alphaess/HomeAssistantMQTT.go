package alphaess

import (
	"encoding/json"
)

type HasMQTTConfig struct {
	DeviceClass       string `json:"device_class"`
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string `json:"value_template"`
	ExpireAfter       int    `json:"expire_after,omitempty"`
	Icon              string `json:"icon,omitempty"`
	//TODO look at adding HAS MQTT support for other fields
	//state_class		  string `json:"state_class,omitempty"`
	//icon_template		  strimg `json:"icon_template,omitempty"`
}

func PublishHASEntityConfig() {
	// discovery definition: https://www.home-assistant.io/docs/mqtt/discovery/
	// device class & Unit of measurements from: https://github.com/home-assistant/core/blob/d7ac4bd65379e11461c7ce0893d3533d8d8b8cbf/homeassistant/const.py#L379
	// device class descriptions from: https://www.home-assistant.io/integrations/sensor/

	var mqClient = gClient
	var myHASConfig HasMQTTConfig
	myHASConfig.StateTopic = gMQTTTopic + "/state"
	myHASConfig.Name = "AlphaESS - Last Updated"
	myHASConfig.DeviceClass = "timestamp"
	myHASConfig.UnitOfMeasurement = ""
	myHASConfig.ValueTemplate = "{{ value_json.Time}}"
	res, _ := json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/LastUpdateTime/config", string(res))
	//SN    	 string `json:"SN"`				//	"SN":"AL2002321010043",
	//PrealL1  float32 `json:"PrealL1,string"`	//	"PrealL1":"756",

	myHASConfig.Name = "AlphaESS - PrealL1"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.ValueTemplate = "{{ value_json.PrealL1}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PrealL1/config", string(res))
	//PrealL2  float32 `json:"PrealL2,string"`
	//PrealL3  float32 `json:"PrealL3,string"`

	//PmeterL1 int `type:"integer" json:"PmeterL1,string"`
	myHASConfig.Name = "AlphaESS - FeedIn/Grid Power In"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.Icon = "mdi:transmission-tower"
	myHASConfig.ValueTemplate = "{{ value_json.PmeterL1}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PmeterL1/config", string(res))
	// distinct grid usage values
	myHASConfig.Name = "AlphaESS - Grid Consumption"
	myHASConfig.ValueTemplate = "{%set mylist = states(\"sensor.alphaess_feedin_grid_power_in\")|float, 0|float,%}{{ mylist|max|float }}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PmeterGridIn/config", string(res))
	myHASConfig.Name = "AlphaESS - Grid Return"
	myHASConfig.ValueTemplate = "{%set mylist = states(\"sensor.alphaess_feedin_grid_power_in\")|float * -1, 0|float,%}{{ mylist|max|float }}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PmeterGridOut/config", string(res))

	myHASConfig.Icon = ""
	myHASConfig.Name = "AlphaESS - PowerMeter2"
	myHASConfig.ValueTemplate = "{{ value_json.PmeterL2}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PmeterL2/config", string(res))

	//PmeterL2 int `type:"integer" json:"PmeterL2,string"`
	//PmeterL3 int `json:"PmeterL3,string"`
	//PmeterDC int `json:"PmeterDC,string"`

	//Pbat     float32 `json:"Pbat,string"`	//	"Pbat":"387.4500",
	myHASConfig.Name = "AlphaESS - BatteryRQ Load(Out)"
	myHASConfig.ValueTemplate = "{{ value_json.Pbat}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/Pbat/config", string(res))
	// distinct battery values
	myHASConfig.Name = "AlphaESS - Battery Load"
	myHASConfig.ValueTemplate = "{%set mylist = states(\"sensor.alphaess_batteryrq_load_out\")|float, 0|float,%}{{ mylist|max|float }}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PBatLoad/config", string(res))
	myHASConfig.Name = "AlphaESS - Battery Charge"
	myHASConfig.ValueTemplate = "{%set mylist = states(\"sensor.alphaess_batteryrq_load_out\")|float * -1, 0|float,%}{{ mylist|max|float }}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PBatCharging/config", string(res))

	//Sva      int `json:"Sva,string"`		//	"Sva":"826",
	myHASConfig.Name = "AlphaESS-Sva"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.ValueTemplate = "{{ value_json.Sva}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/Sva/config", string(res))

	//VarAC    int `json:"VarAC,string"`	//	"VarAC":"-541",
	myHASConfig.Name = "AlphaESS-VarAC"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.ValueTemplate = "{{ value_json.VarAC}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/VarAC/config", string(res))
	//VarDC    int `json:"VarDC,string"`	//	"VarDC":"0",
	myHASConfig.Name = "AlphaESS-VarDC"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.ValueTemplate = "{{ value_json.VarDC}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/VarDC/config", string(res))

	//SOC      float32 `type:"float32" json:"SOC,string"` //"SOC":"24.0"}
	myHASConfig.Name = "AlphaESS - State of Charge"
	myHASConfig.DeviceClass = "battery"
	myHASConfig.UnitOfMeasurement = "%"
	myHASConfig.ValueTemplate = "{{ value_json.SOC}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/SOC/config", string(res))

	//Ppv1     int `json:"Ppv1,string"`			//	"Ppv1":"160",
	myHASConfig.Name = "AlphaESS - Power from PV1"
	myHASConfig.DeviceClass = "power"
	myHASConfig.UnitOfMeasurement = "W"
	myHASConfig.ValueTemplate = "{{ value_json.Ppv1}}"
	myHASConfig.Icon = "mdi:solar-panel"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/Ppv1/config", string(res))

	//Ppv2     int `json:"Ppv2,string"`			//	"Ppv2":"273",
	myHASConfig.Name = "AlphaESS - Power from PV2"
	myHASConfig.ValueTemplate = "{{ value_json.Ppv2}}"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/Ppv2/config", string(res))

	// TEMPLATE: mdi:solar-power AlphaESS-TotalSolar
	myHASConfig.Name = "AlphaESS - TotalSolar"
	myHASConfig.ValueTemplate = "{{states('sensor.alphaess_power_from_pv1')|int + states('sensor.alphaess_power_from_pv2')|int}}"
	myHASConfig.Icon = "mdi:solar-power"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/PTotal/config", string(res))

	// composite load value
	myHASConfig.Name = "AlphaESS - Total Load"
	myHASConfig.ValueTemplate = "{{states('sensor.alphaess_totalsolar')|int + states('sensor.alphaess_feedin_grid_power_in')|int}}"
	myHASConfig.Icon = "mdi:power-socket-uk"
	res, _ = json.Marshal(myHASConfig)
	publishMQTT(mqClient, gMQTTTopic+"/LoadTotal/config", string(res))
}
