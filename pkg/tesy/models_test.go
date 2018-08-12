package tesy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStatusURL(t *testing.T) {
	tesy := Tesy{Address: "10.0.1.1"}

	url := tesy.StatusURL()
	if url != "http://10.0.1.1/status" {
		t.Errorf("StatusURL didn't return the expected url")
	}
}

func TestDevStatsURL(t *testing.T) {
	tesy := Tesy{Address: "10.0.1.1"}

	url := tesy.DevStatsURL()
	if url != "http://10.0.1.1/devstat" {
		t.Errorf("DevStatsURL didn't return the expected url")
	}
}

func TestCalcResURL(t *testing.T) {
	tesy := Tesy{Address: "10.0.1.1"}

	url := tesy.CalcResURL()
	if url != "http://10.0.1.1/calcRes" {
		t.Errorf("CalcResURL didn't return the expected url")
	}
}

func TestResetCalcURL(t *testing.T) {
	tesy := Tesy{Address: "10.0.1.1"}

	url := tesy.ResetCalcURL()
	if url != "http://10.0.1.1/resetPow" {
		t.Errorf("ResetCalcURL didn't return the expected url")
	}
}

func TestSetTempURL(t *testing.T) {
	tesy := Tesy{Address: "10.0.1.1"}

	url := tesy.SetTempURL()
	if url != "http://10.0.1.1/setTemp" {
		t.Errorf("SetTempURL didn't return the expected url")
	}
}

func TestResetPowerCalc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"err": "0"}`))
	}))

	defer server.Close()
	url, _ := url.Parse(server.URL)
	tesy := Tesy{Address: url.Host}
	err := tesy.ResetPowerCalc()
	if err != nil {
		t.Errorf("ResetPowerCalc did not reply as expected: %s", err)
	}
}

func TestSetTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("val") != "12" {
			t.Errorf("SetTemperature tried to set the wrong temperature")
		}
		rw.Write([]byte(`{"stat": "ok"}`))
	}))

	defer server.Close()
	url, _ := url.Parse(server.URL)
	tesy := Tesy{Address: url.Host}
	err := tesy.SetTemperature(12)
	if err != nil {
		t.Errorf("SetTemperature did not reply as expected: %s", err)
	}
}

func TestTimeSinceReset(t *testing.T) {
	tesy := Tesy{
		PowerUsage: PowerUsage{
			ResetDate: "2018-08-12 19:00:00",
		},
		Status: TesyStatus{
			Date: "2018-08-12 20:00",
			Tz:   "CEST",
		},
	}

	hours, _ := tesy.TimeSinceReset()
	if hours != 1 {
		t.Errorf("TimeSinceReset should have returned 1, returned: %f", hours)
	}
}

func TestKWH(t *testing.T) {
	pu := PowerUsage{
		Sum:  1500,
		Watt: 2400,
	}

	kwh := pu.KWH()

	if kwh != 1 {
		t.Errorf("KWH() should have returned 1 but returned: %f", kwh)
	}
}

func TestGetStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"heater_state":"READY","gradus":"61.0","ref_gradus": "70","volume":"80"}`))
	}))

	defer server.Close()
	url, _ := url.Parse(server.URL)
	tesy := Tesy{Address: url.Host}
	tesyStatus, _ := tesy.GetStatus()

	if tesyStatus.Gradus != 61.0 {
		t.Errorf("Gradus should be 61.0 but is: %f", tesyStatus.Gradus)
	}
}

func TestGetCalcRes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"sum":"1500","resetDate":"2018-08-12 19:00:00","watt":"2400"}`))
	}))

	defer server.Close()
	url, _ := url.Parse(server.URL)
	tesy := Tesy{Address: url.Host}
	tesyPower, _ := tesy.GetCalcRes()

	if tesyPower.Sum != 1500 {
		t.Errorf("Sum should be 1500 but is: %d", tesyPower.Sum)
	}
}

func TestGetDevStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"devid":"foo FW15.2O","macaddr":"3:1:33:7","inetdev":"WIFI"}`))
	}))

	defer server.Close()
	url, _ := url.Parse(server.URL)
	tesy := Tesy{Address: url.Host}
	_ = tesy.GetDevStatus()

	if tesy.DevID != "foo FW15.2O" {
		t.Errorf("DevID should be %s but got: %s", "foo FW15.2O", tesy.DevID)
	}
}
