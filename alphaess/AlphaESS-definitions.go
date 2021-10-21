package alphaess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

//const DefaultExpiry = 1500 // seconds after which to expire a sensor value;
// TODO Unit tests for AlphaESS-definitions; parse from files.

const SERIALRQPATTERN = "{\"SN\""
const CONFIGRSPATTERN = "ZipCode"

type AuthRQ struct {
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	CompanyName string `json:"CompanyName"`
}

//SuccessRS {"Status":"Success"}
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
//	{"Command":"SetConfig","CmdIndex":"49662827","Status":"Success"}
type CommandRQ struct {
	Command  string `json:"Command"`
	CmdIndex int64  `json:"CmdIndex,string"`
	Status   string `json:"Status,omitempty"`
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

// ConfigRS is the main config message shared between the cloud servers and battery system/
type ConfigRS struct {
	Status             string  `json:"Status,omitempty"`
	SN                 string  `json:"SN"`
	Address            string  `json:"Address"`
	ZipCode            string  `json:"ZipCode"`
	Country            string  `json:"Country"`
	PhoneNumber        string  `json:"PhoneNumber"`
	License            string  `json:"License"`
	Popv               float32 `json:"Popv,string"` //"11.50"
	Minv               string  `json:"Minv"`
	Poinv              float32 `json:"Poinv,string"` //"5.00"
	Cobat              float32 `json:"Cobat,string"` //"Cobat":"11.40"
	Mbat               string  `json:"Mbat"`
	Uscapacity         float32 `json:"Uscapacity,string"` //"90.00"
	InstallMeterOption string  `json:"InstallMeterOption"`
	Mmeter             string  `json:"Mmeter"`
	PVMeterMode        string  `json:"PVMeterMode,omitempty"`
	CTRate             int     `json:"CTRate,string"`
	PVMeterCTRate      int     `json:"PVMeterCTRate,string"`
	GridMeterCTE       int     `json:"GridMeterCTE,string"`
	PVMeterCTE         int     `json:"PVMeterCTE,string"`
	BatterySN1         string  `json:"BatterySN1,omitempty"`
	BatterySN2         string  `json:"BatterySN2,omitempty"`
	BatterySN3         string  `json:"BatterySN3,omitempty"`
	BatterySN4         string  `json:"BatterySN4,omitempty"`
	BatterySN5         string  `json:"BatterySN5,omitempty"`
	BatterySN6         string  `json:"BatterySN6,omitempty"`
	BatterySN7         string  `json:"BatterySN7,omitempty"`
	BatterySN8         string  `json:"BatterySN8,omitempty"`
	BatterySN9         string  `json:"BatterySN9,omitempty"`
	BatterySN10        string  `json:"BatterySN10,omitempty"`
	BatterySN11        string  `json:"BatterySN11,omitempty"`
	BatterySN12        string  `json:"BatterySN12,omitempty"`
	BatterySN13        string  `json:"BatterySN13,omitempty"`
	BatterySN14        string  `json:"BatterySN14,omitempty"`
	BatterySN15        string  `json:"BatterySN15,omitempty"`
	BatterySN16        string  `json:"BatterySN16,omitempty"`
	BatterySN17        string  `json:"BatterySN17,omitempty"`
	BatterySN18        string  `json:"BatterySN18,omitempty"`
	BMSVersion         string  `json:"BMSVersion,omitempty"`
	EMSVersion         string  `json:"EMSVersion"`
	InvVersion         string  `json:"InvVersion"`
	InvSN              string  `json:"InvSN"`
	ACDC               string  `json:"ACDC"`
	Generator          bool    `json:"Generator"`
	BackUpBox          bool    `json:"BackUpBox"`
	Fan                bool    `json:"Fan"`
	GridCharge         bool    `json:"GridCharge"` //true
	GridChargeWE       bool    `json:"GridChargeWE"`
	CtrDis             bool    `json:"CtrDis"`
	CtrDisWE           bool    `json:"CtrDisWE"`
	TimeChaF1          int     `json:"TimeChaF1,string"`
	TimeChaE1          int     `json:"TimeChaE1,string"`
	TimeChaF2          int     `json:"TimeChaF2,string"`
	TimeChaE2          int     `json:"TimeChaE2,string"`
	TimeDisF1          int     `json:"TimeDisF1,string"`
	TimeDisE1          int     `json:"TimeDisE1,string"`
	TimeDisF2          int     `json:"TimeDisF2,string"`
	TimeDisE2          int     `json:"TimeDisE2,string"`
	BatHighCap         float32 `json:"BatHighCap,string"`
	BatHighCapWE       float32 `json:"BatHighCapWE,string"`
	BatUseCap          float32 `json:"BatUseCap,string"`
	BatUseCapWE        float32 `json:"BatUseCapWE,string"`
	SetMode            int     `json:"SetMode,string"`
	SetPhase           int     `json:"SetPhase,string"`
	SetFeed            int     `json:"SetFeed,string"`
	BakBoxSN           string  `json:"BakBoxSN"`
	SCBSN              string  `json:"SCBSN,omitempty"`
	BakBoxVer          string  `json:"BakBoxVer,omitempty"`
	SCBVer             string  `json:"SCBVer,omitempty"`
	BMUModel           string  `json:"BMUModel,omitempty"`
	GeneratorMode      string  `json:"GeneratorMode"`
	GCSOCStart         string  `json:"GCSOCStart"`
	GCSOCEnd           string  `json:"GCSOCEnd"`
	GCTimeStart        string  `json:"GCTimeStart"`
	GCTimeEnd          string  `json:"GCTimeEnd"`
	GCOutputMode       string  `json:"GCOutputMode"`
	GCChargePower      string  `json:"GCChargePower"`
	GCRatedPower       string  `json:"GCRatedPower"`
	EmsLanguage        int     `json:"EmsLanguage"`
	L1Priority         int     `json:"L1Priority,string"`
	L2Priority         int     `json:"L2Priority,string"`
	L3Priority         int     `json:"L3Priority,string"`
	L1SocLimit         float32 `json:"L1SocLimit,string"`
	L2SocLimit         float32 `json:"L2SocLimit,string"`
	L3SocLimit         float32 `json:"L3SocLimit,string"`
	FirmwareVersion    string  `json:"FirmwareVersion"`
	OnGridCap          float32 `json:"OnGridCap,string"`
	StorageCap         float32 `json:"StorageCap,string"`
	BatReady           int     `json:"BatReady,string"`
	MeterACNegate      int     `json:"MeterACNegate,string"`
	MeterDCNegate      int     `json:"MeterDCNegate,string"`
	Safe               int     `json:"Safe,string"`
	PowerFact          int     `json:"PowerFact,string"`
	Volt5MinAvg        int     `json:"Volt5MinAvg,string"`
	Volt10MinAvg       int     `json:"Volt10MinAvg,string"`
	TempThreshold      int     `json:"TempThreshold,string"`
	OutCurProtect      int     `json:"OutCurProtect,string"`
	DCI                int     `json:"DCI,string"`
	RCD                int     `json:"RCD,string"`
	PvISO              int     `json:"PvISO,string"`
	ChargeBoostCur     int     `json:"ChargeBoostCur,string"`
	Channel1           int     `json:"Channel1,string"`
	ControlMode1       int     `json:"ControlMode1,string"`
	StartTime1A        string  `json:"StartTime1A"`
	EndTime1A          string  `json:"EndTime1A"`
	StartTime1B        string  `json:"StartTime1B"`
	EndTime1B          string  `json:"EndTime1B"`
	Date1              string  `json:"Date1"`
	ChargeSOC1         string  `json:"ChargeSOC1,omitempty"`
	ChargeMode1        string  `json:"ChargeMode1"`
	ChargeWeekend      int     `json:"ChargeWeekend,string"`
	ChargeWorkDays     int     `json:"ChargeWorkDays,string"`
	UPS1               string  `json:"UPS1"`
	SwitchOn1          string  `json:"SwitchOn1"`
	SwitchOff1         string  `json:"SwitchOff1"`
	Delay1             string  `json:"Delay1"`
	Duration1          string  `json:"Duration1"`
	Pause1             string  `json:"Pause1"`
	Channel2           string  `json:"Channel2"`
	ControlMode2       string  `json:"ControlMode2"`
	StartTime2A        string  `json:"StartTime2A"`
	EndTime2A          string  `json:"EndTime2A"`
	StartTime2B        string  `json:"StartTime2B"`
	EndTime2B          string  `json:"EndTime2B"`
	Date2              string  `json:"Date2"`
	ChargeSOC2         string  `json:"ChargeSOC2,omitempty"`
	ChargeMode2        string  `json:"ChargeMode2"`
	UPS2               string  `json:"UPS2"`
	SwitchOn2          string  `json:"SwitchOn2"`
	SwitchOff2         string  `json:"SwitchOff2"`
	Delay2             string  `json:"Delay2"`
	Duration2          string  `json:"Duration2"`
	Pause2             string  `json:"Pause2"`
	UseCt              string  `json:"use_ct"`
	CtRate             string  `json:"ct_rate"`
	InstallModule      string  `json:"InstallModule"`
	StringAE           string  `json:"StringAE"`
	StringBE           string  `json:"StringBE"`
	StringCE           string  `json:"StringCE"`
	NetType            string  `json:"NetType"`
	WifiSN             string  `json:"WifiSN"`
	WifiSW             string  `json:"WifiSW"`
	WifiHW             string  `json:"WifiHW"`
	SelfUseOrEconomic  int     `json:"SelfUseOrEconomic,string"`
	DG_Cap             int     `json:"DG_Cap,string"`
	FAAEnable          int     `json:"FAAEnable,string"`
	ReliefMode         int     `json:"ReliefMode"`
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
	loc, _ := time.LoadLocation(gTZLocation)
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
	} else if bytes.Index(rawData, []byte(CONFIGRSPATTERN)) >= 0 {
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
			ErrorLog("decoding AuthRQ: " + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("\"Status\"")) >= 0 {
		jsonResult := SuccessRS{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.Status) == 0 {
			ErrorLog("decoding SuccessRS: " + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte(SERIALRQPATTERN)) >= 0 {
		jsonResult := SerialRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.SN) == 0 {
			ErrorLog("decoding SerialRQ: " + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("{\"Command\"")) >= 0 {
		jsonResult := CommandRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if jsonResult.CmdIndex <= 0 {
			ErrorLog("decoding CommandRQ: " + string(rawData))
		} else {
			result = jsonResult
		}
	} else if bytes.Index(rawData, []byte("{\"CmdIndex\"")) >= 0 {
		jsonResult := CommandIndexRQ{}
		err = json.Unmarshal(rawData, &jsonResult)
		if len(jsonResult.CmdIndex) == 0 {
			ErrorLog("decoding CommandIndexRQ: " + string(rawData))
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
