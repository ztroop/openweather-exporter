// Copyright 2021 Zackary Troop
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package owm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type OWMHandler struct {
	ApiKey    string
	Latitude  float64
	Longitude float64
	Unit      string
	Current   OneCallValues
	Pollution PollutionValues
}

type CurrentValues struct {
	DateTime    int64   `json:"dt"`
	Sunrise     int64   `json:"sunrise"`
	Sunset      int64   `json:"sunset"`
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
	DewPoint    float64 `json:"dew_point"`
	UVI         float64 `json:"uvi"`
	Clouds      int8    `json:"clouds"`
	Visibility  float64 `json:"visibility"`
	WindSpeed   float64 `json:"wind_speed"`
	WindDegree  float64 `json:"wind_deg"`
	WindGust    float64 `json:"wind_gust"`
	Rain1Hour   float64 `json:"rain.1h"`
	Snow1Hour   float64 `json:"snow.1h"`
}

type OneCallValues struct {
	Latitude       float64       `json:"lat"`
	Longitude      float64       `json:"lon"`
	Timezone       string        `json:"timezone"`
	TimezoneOffset float64       `json:"timezone_offset"`
	Values         CurrentValues `json:"current"`
}

type PollutionComponents struct {
	CO   float64 `json:"co"`
	NO   float64 `json:"no"`
	NO2  float64 `json:"no2"`
	O3   float64 `json:"o3"`
	SO2  float64 `json:"so2"`
	PM25 float64 `json:"pm2_5"`
	PM10 float64 `json:"pm10"`
	NH3  float64 `json:"nh3"`
}

type PollutionMain struct {
	AQI int8 `json:"aqi"`
}

type PollutionItem struct {
	DateTime   int64 `json:"dt"`
	Main       PollutionMain
	Components PollutionComponents
}

type PollutionValues struct {
	List []PollutionItem
}

const oneCallUrl = "https://api.openweathermap.org/data/2.5/onecall?lat=%v&lon=%v&exclude=minutely,hourly,daily&units=%s&appid=%s"
const pollutionUrl = "https://api.openweathermap.org/data/2.5/air_pollution?lat=%v&lon=%v&appid=%s"

func (c *OWMHandler) FetchData() {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get(fmt.Sprintf(oneCallUrl, c.Latitude, c.Longitude, c.Unit, c.ApiKey))
	if err != nil {
		log.Error("Error requesting OneCall API:", err)
	} else {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Error("Error reading OneCall data:", err)
		} else {
			json.Unmarshal(body, &c.Current)
		}
	}

	resp, err = http.Get(fmt.Sprintf(pollutionUrl, c.Latitude, c.Longitude, c.ApiKey))
	if err != nil {
		log.Error("Error requesting Pollution API:", err)
	} else {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Error("Error reading Pollution data:", err)
		} else {
			json.Unmarshal(body, &c.Pollution)
		}
	}
}

func (c *OWMHandler) SetUnit(degreeUnit string) {
	d := strings.ToUpper(degreeUnit)
	if d == "C" {
		c.Unit = "metric"
	} else if d == "F" {
		c.Unit = "imperial"
	} else {
		c.Unit = "standard"
	}
}

func NewOWMHandler(key string, lon float64, lat float64) *OWMHandler {
	return &OWMHandler{
		ApiKey:    key,
		Longitude: lon,
		Latitude:  lat,
	}
}
