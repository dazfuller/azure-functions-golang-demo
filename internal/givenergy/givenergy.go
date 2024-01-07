package givenergy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Account struct {
	Data struct {
		Id              int    `json:"id"`
		Name            string `json:"name"`
		Role            string `json:"role"`
		Email           string `json:"email"`
		Address         string `json:"address"`
		Postcode        string `json:"postcode"`
		Country         string `json:"country"`
		TelephoneNumber string `json:"telephone_number"`
		Timezone        string `json:"timezone"`
	} `json:"data"`
}

type CommunicationDevices struct {
	Data []struct {
		SerialNumber   string    `json:"serial_number"`
		Type           string    `json:"type"`
		CommissionDate time.Time `json:"commission_date"`
		Inverter       struct {
			Serial         string    `json:"serial"`
			Status         string    `json:"status"`
			LastOnline     time.Time `json:"last_online"`
			LastUpdated    time.Time `json:"last_updated"`
			CommissionDate time.Time `json:"commission_date"`
			Info           struct {
				BatteryType string `json:"battery_type"`
				Battery     struct {
					NominalCapacity int     `json:"nominal_capacity"`
					NominalVoltage  float64 `json:"nominal_voltage"`
				} `json:"battery"`
				Model         string `json:"model"`
				MaxChargeRate int    `json:"max_charge_rate"`
			} `json:"info"`
			Warranty struct {
				Type       string    `json:"type"`
				ExpiryDate time.Time `json:"expiry_date"`
			} `json:"warranty"`
			FirmwareVersion struct {
				ARM int `json:"ARM"`
				DSP int `json:"DSP"`
			} `json:"firmware_version"`
			Connections struct {
				Batteries []struct {
					ModuleNumber    int    `json:"module_number"`
					Serial          string `json:"serial"`
					FirmwareVersion string `json:"firmware_version"`
					Capacity        struct {
						Full   float64 `json:"full"`
						Design int     `json:"design"`
					} `json:"capacity"`
					CellCount int  `json:"cell_count"`
					HasUsb    bool `json:"has_usb"`
				} `json:"batteries"`
				Meters []interface{} `json:"meters"`
			} `json:"connections"`
			Flags []interface{} `json:"flags"`
		} `json:"inverter"`
	} `json:"data"`
	Links struct {
		First string      `json:"first"`
		Last  string      `json:"last"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
	} `json:"links"`
	Meta struct {
		CurrentPage int    `json:"current_page"`
		From        int    `json:"from"`
		LastPage    int    `json:"last_page"`
		Path        string `json:"path"`
		PerPage     int    `json:"per_page"`
		To          int    `json:"to"`
		Total       int    `json:"total"`
	} `json:"meta"`
}

type InverterMeterData struct {
	Data struct {
		Time   time.Time `json:"time"`
		Status string    `json:"status"`
		Solar  struct {
			Power  int `json:"power"`
			Arrays []struct {
				Array   int     `json:"array"`
				Voltage float64 `json:"voltage"`
				Current float64 `json:"current"`
				Power   float64 `json:"power"`
			} `json:"arrays"`
		} `json:"solar"`
		Grid struct {
			Voltage   float64 `json:"voltage"`
			Current   float64 `json:"current"`
			Power     float64 `json:"power"`
			Frequency float64 `json:"frequency"`
		} `json:"grid"`
		Battery struct {
			Percent     float64 `json:"percent"`
			Power       float64 `json:"power"`
			Temperature float64 `json:"temperature"`
		} `json:"battery"`
		Inverter struct {
			Temperature     float64 `json:"temperature"`
			Power           float64 `json:"power"`
			OutputVoltage   float64 `json:"output_voltage"`
			OutputFrequency float64 `json:"output_frequency"`
			EpsPower        float64 `json:"eps_power"`
		} `json:"inverter"`
		Consumption float64 `json:"consumption"`
	} `json:"data"`
}

type SummarisedMeterData struct {
	Time              time.Time
	Consumption       float64
	SolarGeneration   float64
	BatteryPercentage float64
	BatteryPower      float64
	GridPower         float64
}

func (imd *InverterMeterData) Summarise() SummarisedMeterData {
	var solarPower float64 = 0
	for i := range imd.Data.Solar.Arrays {
		solarPower += imd.Data.Solar.Arrays[i].Power
	}

	return SummarisedMeterData{
		Time:              imd.Data.Time,
		Consumption:       imd.Data.Consumption,
		SolarGeneration:   solarPower,
		BatteryPercentage: imd.Data.Battery.Percent,
		BatteryPower:      imd.Data.Battery.Power,
		GridPower:         imd.Data.Grid.Power,
	}
}

type Api struct {
	AccountEndpoint              string
	CommunicationDevicesEndpoint string
	InverterDataEndpoint         string
	ApiKey                       string
}

func NewAccountApi(apiKey string) Api {
	return Api{
		AccountEndpoint:              "v1/account",
		CommunicationDevicesEndpoint: "v1/communication-device",
		InverterDataEndpoint:         "v1/inverter/%s/system-data/latest",
		ApiKey:                       apiKey,
	}
}

func (a *Api) GetAccountDetails() Account {
	acctResponse := Account{}
	err := a.makeApiRequest(a.AccountEndpoint, &acctResponse)
	if err != nil {
		log.Fatal("Unable to perform account request ", err)
	}
	return acctResponse
}

func (a *Api) GetCommunicationDevices() CommunicationDevices {
	devices := CommunicationDevices{}
	err := a.makeApiRequest(a.CommunicationDevicesEndpoint, &devices)
	if err != nil {
		log.Fatal("Unable to perform communication devices request ", err)
	}
	return devices
}

func (a *Api) GetLatestInverterData(serialNumber string) InverterMeterData {
	path := fmt.Sprintf(a.InverterDataEndpoint, serialNumber)
	data := InverterMeterData{}
	err := a.makeApiRequest(path, &data)
	if err != nil {
		log.Fatal("Unable to perform communication devices request ", err)
	}
	return data
}

func (a *Api) makeApiRequest(path string, responseTarget interface{}) error {
	requestUrl := fmt.Sprintf("https://api.givenergy.cloud/%s", path)
	fmt.Println("Making request to: ", requestUrl)
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %s", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.ApiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s", err)
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(responseTarget)
}
