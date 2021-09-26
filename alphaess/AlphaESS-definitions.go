package alphaess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

//const DefaultExpiry = 1500 // seconds after which to expire a sensor value;

type AuthRQ struct {
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	CompanyName string `json:"CompanyName"`
}

type SuccessRS struct {
	Status string `json:"Status"`
}

type GenericRQ struct {
	MsgType     string `json:"MsgType"`
	MsgContent  string `json:"MsgContent"`
	Description string `json:"Description"`
}

type SerialRQ struct {
	SN string `json:"SN"`
}
type StatusRQ struct {
	Time     Timestamp `json:"Time"`
	SN       string    `json:"SN"`
	Ppv1     int       `json:"Ppv1,string"`
	Ppv2     int       `json:"Ppv2,string"`
	PrealL1  float32   `json:"PrealL1,string"`
	PrealL2  float32   `json:"PrealL2,string"`
	PrealL3  float32   `json:"PrealL3,string"`
	PmeterL1 int       `type:"integer" json:"PmeterL1,string"`
	PmeterL2 int       `type:"integer" json:"PmeterL2,string"`
	PmeterL3 int       `json:"PmeterL3,string"`
	PmeterDC int       `json:"PmeterDC,string"`
	Pbat     float32   `json:"Pbat,string"`
	Sva      int       `json:"Sva,string"`
	VarAC    int       `json:"VarAC,string"`
	VarDC    int       `json:"VarDC,string"`
	SOC      float32   `type:"float32" json:"SOC,string"`
}

// CommandRQ {"Command":"SetConfig","CmdIndex":"35235904"}
type CommandRQ struct {
	Command  string `json:"Command"`
	CmdIndex string `json:"CmdIndex"`
}

// CommandIndexRQ {"CmdIndex":"80867000","Command":"Resume","Parameter1":"2021/08/27 23:16:32","Parameter2":"10"}
type CommandIndexRQ struct {
	CmdIndex   string `json:"CmdIndex"`
	Command    string `json:"Command"`
	Parameter1 string `json:"Parameter1"`
	Parameter2 string `json:"Parameter2"`
}

type BatteryRQ struct {
	Time           time.Time `json:"Time"`
	SN             string    `json:"SN"`
	Ppv1           string    `json:"Ppv1"`
	Ppv2           string    `json:"Ppv2"`
	Upv1           string    `json:"Upv1"`
	Upv2           string    `json:"Upv2"`
	Ua             string    `json:"Ua"`
	Ub             string    `json:"Ub"`
	Uc             string    `json:"Uc"`
	Fac            string    `json:"Fac"`
	Ubus           string    `json:"Ubus"`
	PrealL1        string    `json:"PrealL1"`
	PrealL2        string    `json:"PrealL2"`
	PrealL3        string    `json:"PrealL3"`
	Tinv           string    `json:"Tinv"`
	PacL1          string    `json:"PacL1"`
	PacL2          string    `json:"PacL2"`
	PacL3          string    `json:"PacL3"`
	InvWorkMode    string    `json:"InvWorkMode"`
	EpvTotal       string    `json:"EpvTotal"`
	Einput         string    `json:"Einput"`
	Eoutput        string    `json:"Eoutput"`
	Echarge        string    `json:"Echarge"`
	PmeterL1       string    `json:"PmeterL1"`
	PmeterL2       string    `json:"PmeterL2"`
	PmeterL3       string    `json:"PmeterL3"`
	PmeterDC       string    `json:"PmeterDC"`
	Pbat           string    `json:"Pbat"`
	SOC            string    `json:"SOC"`
	BatV           string    `json:"BatV"`
	BatC           string    `json:"BatC"`
	FlagBms        string    `json:"FlagBms"`
	BmsWork        string    `json:"BmsWork"`
	Pcharge        string    `json:"Pcharge"`
	Pdischarge     string    `json:"Pdischarge"`
	BmsRelay       string    `json:"BmsRelay"`
	BmsNum         string    `json:"BmsNum"`
	VcellLow       string    `json:"VcellLow"`
	VcellHigh      string    `json:"VcellHigh"`
	TcellLow       string    `json:"TcellLow"`
	TcellHigh      string    `json:"TcellHigh"`
	IdTempelover   string    `json:"IdTempelover"`
	IdTempEover    string    `json:"IdTempEover"`
	IdTempediffe   string    `json:"IdTempediffe"`
	IdChargcurre   string    `json:"IdChargcurre"`
	IdDischcurre   string    `json:"IdDischcurre"`
	IdCellvolover  string    `json:"IdCellvolover"`
	IdCellvollower string    `json:"IdCellvollower"`
	IdSoclower     string    `json:"IdSoclower"`
	IdCellvoldiffe string    `json:"IdCellvoldiffe"`
	BatC1          string    `json:"BatC1"`
	BatC2          string    `json:"BatC2"`
	BatC3          string    `json:"BatC3"`
	BatC4          string    `json:"BatC4"`
	BatC5          string    `json:"BatC5"`
	BatC6          string    `json:"BatC6"`
	SOC1           string    `json:"SOC1"`
	SOC2           string    `json:"SOC2"`
	SOC3           string    `json:"SOC3"`
	SOC4           string    `json:"SOC4"`
	SOC5           string    `json:"SOC5"`
	SOC6           string    `json:"SOC6"`
	ErrInv         string    `json:"ErrInv"`
	WarInv         string    `json:"WarInv"`
	ErrEms         string    `json:"ErrEms"`
	ErrBms         string    `json:"ErrBms"`
	ErrMeter       string    `json:"ErrMeter"`
	ErrBackupBox   string    `json:"ErrBackupBox"`
	EGridCharge    string    `json:"EGridCharge"`
	EDischarge     string    `json:"EDischarge"`
	EmsStatus      string    `json:"EmsStatus"`
	InvBatV        string    `json:"InvBatV"`
	BmsShutdown    string    `json:"BmsShutdown"`
	BmuRelay       string    `json:"BmuRelay"`
	BmsHardVer1    string    `json:"BmsHardVer1"`
	BmsHardVer2    string    `json:"BmsHardVer2"`
	BmsHardVer3    string    `json:"BmsHardVer3"`
	DispatchSwitch string    `json:"DispatchSwitch"`
	Pdispatch      string    `json:"Pdispatch"`
	DispatchSoc    string    `json:"DispatchSoc"`
	DispatchMode   string    `json:"DispatchMode"`
	PMeterDCL1     string    `json:"PMeterDCL1"`
	PMeterDCL2     string    `json:"PMeterDCL2"`
	PMeterDCL3     string    `json:"PMeterDCL3"`
	MeterDCUa      string    `json:"MeterDCUa"`
	MeterDCUb      string    `json:"MeterDCUb"`
	MeterDCUc      string    `json:"MeterDCUc"`
	Meter1Actpower string    `json:"meter1.actpower"`
	GridF          string    `json:"GridF"`
	GridVolt       string    `json:"GridVolt"`
	DSPDebug1      string    `json:"DSPDebug1"`
	DSPDebug2      string    `json:"DSPDebug2"`
	DSPDebug3      string    `json:"DSPDebug3"`
	DSPDebug4      string    `json:"DSPDebug4"`
	DSPDebug5      string    `json:"DSPDebug5"`
	DSPDebug6      string    `json:"DSPDebug6"`
	DSPDebug7      string    `json:"DSPDebug7"`
	DSPDebug8      string    `json:"DSPDebug8"`
	DSPDebug9      string    `json:"DSPDebug9"`
	DSPDebug10     string    `json:"DSPDebug10"`
	DSPDebug11     string    `json:"DSPDebug11"`
	DSPDebug12     string    `json:"DSPDebug12"`
	DSPDebug13     string    `json:"DSPDebug13"`
	DSPDebug14     string    `json:"DSPDebug14"`
	DSPDebugChg1   string    `json:"DSPDebugChg1"`
	DSPDebugChg2   string    `json:"DSPDebugChg2"`
	DSPDebugChg3   string    `json:"DSPDebugChg3"`
	DSPDebugChg4   string    `json:"DSPDebugChg4"`
	DSPDebugChg5   string    `json:"DSPDebugChg5"`
	DSPDebugChg6   string    `json:"DSPDebugChg6"`
	DSPDebugChg7   string    `json:"DSPDebugChg7"`
	DSPDebugChg8   string    `json:"DSPDebugChg8"`
	DSPDebugChg9   string    `json:"DSPDebugChg9"`
	DSPDebugChg10  string    `json:"DSPDebugChg10"`
	OVP1Threshold  string    `json:"OVP1Threshold"`
	OVP1TripValue  string    `json:"OVP1TripValue"`
	OVP1TripTime   string    `json:"OVP1TripTime"`
	OVP2Threshold  string    `json:"OVP2Threshold"`
	OVP2TripValue  string    `json:"OVP2TripValue"`
	OVP2TripTime   string    `json:"OVP2TripTime"`
	UVP1Threshold  string    `json:"UVP1Threshold"`
	UVP1TripValue  string    `json:"UVP1TripValue"`
	UVP1TripTime   string    `json:"UVP1TripTime"`
	UVP2Threshold  string    `json:"UVP2Threshold"`
	UVP2TripValue  string    `json:"UVP2TripValue"`
	UVP2TripTime   string    `json:"UVP2TripTime"`
	OFPThreshold   string    `json:"OFPThreshold"`
	OFPTripValue   string    `json:"OFPTripValue"`
	OFPTripTime    string    `json:"OFPTripTime"`
	UFPThreshold   string    `json:"UFPThreshold"`
	UFPTripValue   string    `json:"UFPTripValue"`
	UFPTripTime    string    `json:"UFPTripTime"`
	PowerFactor    string    `json:"PowerFactor"`
	Eirp           string    `json:"Eirp"`
	CSQ            string    `json:"CSQ"`
}

type ConfigRS struct {
	SN                 string `json:"SN"`
	Address            string `json:"Address"`
	ZipCode            string `json:"ZipCode"`
	Country            string `json:"Country"`
	PhoneNumber        string `json:"PhoneNumber"`
	License            string `json:"License"`
	Popv               string `json:"Popv"`
	Minv               string `json:"Minv"`
	Poinv              string `json:"Poinv"`
	Cobat              string `json:"Cobat"`
	Mbat               string `json:"Mbat"`
	Uscapacity         string `json:"Uscapacity"`
	InstallMeterOption string `json:"InstallMeterOption"`
	Mmeter             string `json:"Mmeter"`
	PVMeterMode        string `json:"PVMeterMode"`
	CTRate             string `json:"CTRate"`
	PVMeterCTRate      string `json:"PVMeterCTRate"`
	GridMeterCTE       string `json:"GridMeterCTE"`
	PVMeterCTE         string `json:"PVMeterCTE"`
	BatterySN1         string `json:"BatterySN1"`
	BatterySN2         string `json:"BatterySN2"`
	BatterySN3         string `json:"BatterySN3"`
	BatterySN4         string `json:"BatterySN4"`
	BatterySN5         string `json:"BatterySN5"`
	BatterySN6         string `json:"BatterySN6"`
	BatterySN7         string `json:"BatterySN7"`
	BatterySN8         string `json:"BatterySN8"`
	BatterySN9         string `json:"BatterySN9"`
	BatterySN10        string `json:"BatterySN10"`
	BatterySN11        string `json:"BatterySN11"`
	BatterySN12        string `json:"BatterySN12"`
	BatterySN13        string `json:"BatterySN13"`
	BatterySN14        string `json:"BatterySN14"`
	BatterySN15        string `json:"BatterySN15"`
	BatterySN16        string `json:"BatterySN16"`
	BatterySN17        string `json:"BatterySN17"`
	BatterySN18        string `json:"BatterySN18"`
	BMSVersion         string `json:"BMSVersion"`
	EMSVersion         string `json:"EMSVersion"`
	InvVersion         string `json:"InvVersion"`
	InvSN              string `json:"InvSN"`
	ACDC               string `json:"ACDC"`
	Generator          string `json:"Generator"`
	BackUpBox          string `json:"BackUpBox"`
	Fan                string `json:"Fan"`
	GridCharge         string `json:"GridCharge"`
	CtrDis             string `json:"CtrDis"`
	TimeChaF1          string `json:"TimeChaF1"`
	TimeChaE1          string `json:"TimeChaE1"`
	TimeChaF2          string `json:"TimeChaF2"`
	TimeChaE2          string `json:"TimeChaE2"`
	TimeDisF1          string `json:"TimeDisF1"`
	TimeDisE1          string `json:"TimeDisE1"`
	TimeDisF2          string `json:"TimeDisF2"`
	TimeDisE2          string `json:"TimeDisE2"`
	BatHighCap         string `json:"BatHighCap"`
	BatUseCap          string `json:"BatUseCap"`
	SetMode            string `json:"SetMode"`
	SetPhase           string `json:"SetPhase"`
	SetFeed            string `json:"SetFeed"`
	BakBoxSN           string `json:"BakBoxSN"`
	SCBSN              string `json:"SCBSN"`
	BakBoxVer          string `json:"BakBoxVer"`
	SCBVer             string `json:"SCBVer"`
	BMUModel           string `json:"BMUModel"`
	GeneratorMode      string `json:"GeneratorMode"`
	GCSOCStart         string `json:"GCSOCStart"`
	GCSOCEnd           string `json:"GCSOCEnd"`
	GCTimeStart        string `json:"GCTimeStart"`
	GCTimeEnd          string `json:"GCTimeEnd"`
	GCOutputMode       string `json:"GCOutputMode"`
	GCChargePower      string `json:"GCChargePower"`
	GCRatedPower       string `json:"GCRatedPower"`
	EmsLanguage        string `json:"EmsLanguage"`
	L1Priority         string `json:"L1Priority"`
	L2Priority         string `json:"L2Priority"`
	L3Priority         string `json:"L3Priority"`
	L1SocLimit         string `json:"L1SocLimit"`
	L2SocLimit         string `json:"L2SocLimit"`
	L3SocLimit         string `json:"L3SocLimit"`
	FirmwareVersion    string `json:"FirmwareVersion"`
	OnGridCap          string `json:"OnGridCap"`
	StorageCap         string `json:"StorageCap"`
	BatReady           string `json:"BatReady"`
	MeterACNegate      string `json:"MeterACNegate"`
	MeterDCNegate      string `json:"MeterDCNegate"`
	Safe               string `json:"Safe"`
	PowerFact          string `json:"PowerFact"`
	Volt5MinAvg        string `json:"Volt5MinAvg"`
	Volt10MinAvg       string `json:"Volt10MinAvg"`
	TempThreshold      string `json:"TempThreshold"`
	OutCurProtect      string `json:"OutCurProtect"`
	DCI                string `json:"DCI"`
	RCD                string `json:"RCD"`
	PvISO              string `json:"PvISO"`
	ChargeBoostCur     string `json:"ChargeBoostCur"`
	Channel1           string `json:"Channel1"`
	ControlMode1       string `json:"ControlMode1"`
	StartTime1A        string `json:"StartTime1A"`
	EndTime1A          string `json:"EndTime1A"`
	StartTime1B        string `json:"StartTime1B"`
	EndTime1B          string `json:"EndTime1B"`
	Date1              string `json:"Date1"`
	ChargeSOC1         string `json:"ChargeSOC1"`
	ChargeMode1        string `json:"ChargeMode1"`
	UPS1               string `json:"UPS1"`
	SwitchOn1          string `json:"SwitchOn1"`
	SwitchOff1         string `json:"SwitchOff1"`
	Delay1             string `json:"Delay1"`
	Duration1          string `json:"Duration1"`
	Pause1             string `json:"Pause1"`
	Channel2           string `json:"Channel2"`
	ControlMode2       string `json:"ControlMode2"`
	StartTime2A        string `json:"StartTime2A"`
	EndTime2A          string `json:"EndTime2A"`
	StartTime2B        string `json:"StartTime2B"`
	EndTime2B          string `json:"EndTime2B"`
	Date2              string `json:"Date2"`
	ChargeSOC2         string `json:"ChargeSOC2"`
	ChargeMode2        string `json:"ChargeMode2"`
	UPS2               string `json:"UPS2"`
	SwitchOn2          string `json:"SwitchOn2"`
	SwitchOff2         string `json:"SwitchOff2"`
	Delay2             string `json:"Delay2"`
	Duration2          string `json:"Duration2"`
	Pause2             string `json:"Pause2"`
	UseCt              string `json:"use_ct"`
	CtRate             string `json:"ct_rate"`
	InstallModule      string `json:"InstallModule"`
	StringAE           string `json:"StringAE"`
	StringBE           string `json:"StringBE"`
	StringCE           string `json:"StringCE"`
	NetType            string `json:"NetType"`
	WifiSN             string `json:"WifiSN"`
	WifiSW             string `json:"WifiSW"`
	WifiHW             string `json:"WifiHW"`
}

type Response interface{}

type Timestamp struct {
	time.Time
}

// UnmarshalJSON decodes an int64 timestamp into a time.Time object
func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	// 1. Decode the bytes into an int64
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}
	// 2 - Parse the unix timestamp "2021/08/15 00:27:34"
	loc, _ := time.LoadLocation(gLocation)
	*&p.Time, err = time.ParseInLocation(
		"2006/01/02 15:04:05",
		raw,
		loc,
	)
	return nil
}

func UnmarshalJSON(rawData []byte) (result Response, err error) {
	result = nil
	if bytes.Index(rawData, []byte("\"Time\"")) >= 0 {
		var jsonResult StatusRQ
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.SN) < 5 {
			ErrorLog("decoding StatusRQ " + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("time")) >= 0 {
		jsonResult := BatteryRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.SN) < 5 {
			ErrorLog("decoding BatteryRQ" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("ZipCode")) >= 0 {
		jsonResult := ConfigRS{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.ZipCode) == 0 {
			ErrorLog("decoding ConfigRS" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("UserName")) >= 0 {
		jsonResult := AuthRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.UserName) == 0 {
			ErrorLog("decoding AuthRQ" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("\"Status\"")) >= 0 {
		jsonResult := SuccessRS{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.Status) == 0 {
			ErrorLog("decoding SuccessRS" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("{\"SN\"")) >= 0 {
		jsonResult := SerialRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.SN) == 0 {
			ErrorLog("decoding SerialRQ" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("{\"Command\"")) >= 0 {
		jsonResult := CommandRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.CmdIndex) == 0 {
			ErrorLog("decoding CommandRQ" + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("{\"CmdIndex\"")) >= 0 {
		jsonResult := CommandIndexRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.CmdIndex) == 0 {
			ErrorLog("decoding CommandIndexRQ" + string(rawData))
		} else {
			result = jsonResult
		}
	} else {
		jsonResult := GenericRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if err == nil {
			if len(jsonResult.MsgType) == 0 {
				ErrorLog("decoding GenericRQ: " + string(rawData))
			} else {
				result = jsonResult
			}
		} else {
			DebugLog("unknown message type, trying GenericRQ: " + err.Error())
		}
	}
	return
}
