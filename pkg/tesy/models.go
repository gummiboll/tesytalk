package tesy

import (
	"encoding/json"
	"fmt"
	"time"
)

// SetTempMeessage represents a response from /setTemp
type SetTempMeessage struct {
	Stat string `json:"stat"`
}

// PowerCalcResetMessage represents reponse from /resetPow
type PowerCalcResetMessage struct {
	Err int `json:"err,string"`
}

// PowerUsage represents json response from /calcRes
type PowerUsage struct {
	Sum       int    `json:"sum,string"`
	Watt      int    `json:"watt,string"`
	ResetDate string `json:"resetDate"`
}

// TesyStatus represents json response from /status
type TesyStatus struct {
	HeaterState string  `json:"heater_state"`
	Gradus      float64 `json:"gradus,string"`
	RefGradus   float64 `json:"ref_gradus,string"`
	Volume      int     `json:"volume,string"`
	Watts       int     `json:"watts,string"`
	Date        string  `json:"date"`
	Tz          string  `json:"tz"`
}

// Tesy represents a Tesy
type Tesy struct {
	Address    string
	DevID      string `json:"devid"`
	MacAddr    string `json:"macaddr"`
	Status     TesyStatus
	PowerUsage PowerUsage
}

// StatusURL returns the url for the tesys status endpoint
func (t Tesy) StatusURL() string {
	return fmt.Sprintf("http://%s/status", t.Address)
}

// DevStatsURL returns the url for the tesys devstat endpoint
func (t Tesy) DevStatsURL() string {
	return fmt.Sprintf("http://%s/devstat", t.Address)
}

// CalcResURL returns the url for the tesys calcRes endpoint
func (t Tesy) CalcResURL() string {
	return fmt.Sprintf("http://%s/calcRes", t.Address)
}

// ResetCalcURL returns the url for the tesys resetPow endpoint
func (t Tesy) ResetCalcURL() string {
	return fmt.Sprintf("http://%s/resetPow", t.Address)
}

// SetTempURL returns the url for the tesys setTemp endpoint
func (t Tesy) SetTempURL() string {
	return fmt.Sprintf("http://%s/setTemp", t.Address)
}

// ResetPowerCalc resets the power calculator
func (t Tesy) ResetPowerCalc() error {
	resetMessage := PowerCalcResetMessage{}

	body, err := tesyClient(t.ResetCalcURL())
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &resetMessage)
	if err != nil {
		return err
	}

	if resetMessage.Err != 0 {
		return fmt.Errorf("Failed to reset power calculator")
	}

	return nil

}

// SetTemperature sets the temperature
func (t Tesy) SetTemperature(temp int) error {
	setTempMessage := SetTempMeessage{}

	body, err := tesyClient(fmt.Sprintf("%s?val=%d", t.SetTempURL(), temp))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &setTempMessage)
	if err != nil {
		return err
	}

	if setTempMessage.Stat != "ok" {
		return fmt.Errorf("Failed to set temp, reply was: %s", setTempMessage.Stat)
	}

	return nil
}

// TimeSinceReset calculates time since the power calculation last got reset and prettifies it
func (t Tesy) TimeSinceReset() (float64, error) {
	timeFormat := "2006-01-02 15:04:05 MST"
	resetDate, err := time.Parse(timeFormat, fmt.Sprintf("%s %s", t.PowerUsage.ResetDate, t.Status.Tz))
	if err != nil {
		return 0, err
	}

	tesyDate, err := time.Parse(timeFormat, fmt.Sprintf("%s:00 %s", t.Status.Date, t.Status.Tz))
	if err != nil {
		return 0, err
	}

	return tesyDate.Sub(resetDate).Hours(), nil

}

// KWH calculates kWh
func (p PowerUsage) KWH() float64 {
	return float64(p.Sum) / 3600000 * float64(p.Watt)
}

// GetStatus returns the body of the status endpoint
func (t Tesy) GetStatus() (TesyStatus, error) {
	tesyStatus := TesyStatus{}

	body, err := tesyClient(t.StatusURL())
	if err != nil {
		return tesyStatus, err
	}

	err = json.Unmarshal(body, &tesyStatus)
	if err != nil {
		return tesyStatus, err
	}

	return tesyStatus, nil
}

// GetCalcRes returns the body of the calcRes endpoint
func (t Tesy) GetCalcRes() (PowerUsage, error) {
	tesyPower := PowerUsage{}

	body, err := tesyClient(t.CalcResURL())
	if err != nil {
		return tesyPower, err
	}

	err = json.Unmarshal(body, &tesyPower)
	if err != nil {
		return tesyPower, err
	}

	return tesyPower, nil
}

// GetDevStatus returns the body of the  devstat endpoint
func (t *Tesy) GetDevStatus() error {
	body, err := tesyClient(t.DevStatsURL())
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return err
	}

	return nil
}
