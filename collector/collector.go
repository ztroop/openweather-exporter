// Copyright 2020 Billy Wooten
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"strings"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	prom "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/ztroop/openweather-exporter/geo"
	"github.com/ztroop/openweather-exporter/owm"
)

var notFound = ttlcache.ErrNotFound

type OpenweatherCollector struct {
	ApiKey      string
	DegreesUnit string
	Language    string
	Locations   []Location
	Cache       ttlcache.SimpleCache
	temperature *prom.Desc
	humidity    *prom.Desc
	feelslike   *prom.Desc
	pressure    *prom.Desc
	windspeed   *prom.Desc
	winddegree  *prom.Desc
	windgust    *prom.Desc
	rain1h      *prom.Desc
	snow1h      *prom.Desc
	cloudiness  *prom.Desc
	sunrise     *prom.Desc
	sunset      *prom.Desc
	dewpoint    *prom.Desc
	uvi         *prom.Desc
	aqi         *prom.Desc
	co          *prom.Desc
	no          *prom.Desc
	no2         *prom.Desc
	o3          *prom.Desc
	so2         *prom.Desc
	pm25        *prom.Desc
	pm10        *prom.Desc
	nh3         *prom.Desc
}

type Location struct {
	Location  string
	Latitude  float64
	Longitude float64
}

func resolveLocations(locations string) []Location {
	res := []Location{}

	for _, location := range strings.Split(locations, "|") {
		latitude, longitude, err := geo.Get_coords(openstreetmap.Geocoder(), location)
		if err != nil {
			log.Fatal("failed to resolve location:", err)
		}
		res = append(res, Location{Location: location, Latitude: latitude, Longitude: longitude})
	}
	return res
}

func NewOpenweatherCollector(degreesUnit string, language string, apikey string, locations string, cache *ttlcache.SimpleCache) *OpenweatherCollector {
	return &OpenweatherCollector{
		ApiKey:      apikey,
		DegreesUnit: degreesUnit,
		Language:    language,
		Locations:   resolveLocations(locations),
		Cache:       *cache,
		temperature: prom.NewDesc("openweather_temperature",
			"Current temperature in degrees",
			[]string{"location"}, nil,
		),
		humidity: prom.NewDesc("openweather_humidity",
			"Current relative humidity",
			[]string{"location"}, nil,
		),
		feelslike: prom.NewDesc("openweather_feelslike",
			"Current feels_like temperature in degrees",
			[]string{"location"}, nil,
		),
		pressure: prom.NewDesc("openweather_pressure",
			"Current Atmospheric pressure hPa",
			[]string{"location"}, nil,
		),
		rain1h: prom.NewDesc("openweather_rain1h",
			"Rain volume for last hour, in millimeters",
			[]string{"location"}, nil,
		),
		snow1h: prom.NewDesc("openweather_snow1h",
			"Snow volume for last hour, in millimeters",
			[]string{"location"}, nil,
		),
		windspeed: prom.NewDesc("openweather_windspeed",
			"Current Wind Speed in mph or meters/sec if imperial",
			[]string{"location"}, nil,
		),
		winddegree: prom.NewDesc("openweather_winddegree",
			"Wind direction, degrees (meteorological)",
			[]string{"location"}, nil,
		),
		windgust: prom.NewDesc("openweather_windgust",
			"Wind gust",
			[]string{"location"}, nil,
		),
		cloudiness: prom.NewDesc("openweather_cloudiness",
			"Cloudiness percentage",
			[]string{"location"}, nil,
		),
		sunrise: prom.NewDesc("openweather_sunrise",
			"Sunrise time, unix, UTC",
			[]string{"location"}, nil,
		),
		sunset: prom.NewDesc("openweather_sunset",
			"Sunset time, unix, UTC",
			[]string{"location"}, nil,
		),
		dewpoint: prom.NewDesc("openweather_dewpoint",
			"Atmospheric temperature below which water droplets begin to condense",
			[]string{"location"}, nil,
		),
		uvi: prom.NewDesc("openweather_uvi",
			"Current UV index",
			[]string{"location"}, nil,
		),
		aqi: prom.NewDesc("openweather_aqi",
			"Air quality index",
			[]string{"location"}, nil,
		),
		co: prom.NewDesc("openweather_co",
			"Concentration of CO (Carbon monoxide), μg/m3",
			[]string{"location"}, nil,
		),
		no: prom.NewDesc("openweather_no",
			"Concentration of NO (Nitrogen monoxide), μg/m3",
			[]string{"location"}, nil,
		),
		no2: prom.NewDesc("openweather_no2",
			"Concentration of NO2 (Nitrogen dioxide), μg/m3",
			[]string{"location"}, nil,
		),
		o3: prom.NewDesc("openweather_o3",
			"Concentration of O3 (Ozone), μg/m3",
			[]string{"location"}, nil,
		),
		so2: prom.NewDesc("openweather_so2",
			"Concentration of SO2 (Sulphur dioxide), μg/m3",
			[]string{"location"}, nil,
		),
		pm25: prom.NewDesc("openweather_pm2_5",
			"Concentration of PM2.5 (Fine particles matter), μg/m3",
			[]string{"location"}, nil,
		),
		pm10: prom.NewDesc("openweather_pm10",
			"Concentration of PM10 (Coarse particulate matter), μg/m3",
			[]string{"location"}, nil,
		),
		nh3: prom.NewDesc("openweather_nh3",
			"Concentration of NH3 (Ammonia), μg/m3",
			[]string{"location"}, nil,
		),
	}
}

func (collector *OpenweatherCollector) Describe(ch chan<- *prom.Desc) {
	ch <- collector.temperature
	ch <- collector.humidity
	ch <- collector.feelslike
	ch <- collector.pressure
	ch <- collector.windspeed
	ch <- collector.winddegree
	ch <- collector.windgust
	ch <- collector.rain1h
	ch <- collector.snow1h
	ch <- collector.cloudiness
	ch <- collector.sunrise
	ch <- collector.sunset
	ch <- collector.dewpoint
	ch <- collector.uvi
	ch <- collector.aqi
	ch <- collector.co
	ch <- collector.no
	ch <- collector.no2
	ch <- collector.o3
	ch <- collector.so2
	ch <- collector.pm25
	ch <- collector.pm10
	ch <- collector.nh3
}

func (c *OpenweatherCollector) Collect(ch chan<- prom.Metric) {
	for _, l := range c.Locations {
		var v owm.OWMHandler
		if val, err := c.Cache.Get("OWM"); err != notFound || val != nil {
			v = val.(owm.OWMHandler)
		} else {
			v = *owm.NewOWMHandler(c.ApiKey, l.Longitude, l.Latitude)
			v.SetUnit(c.DegreesUnit)
			v.FetchData()
			c.Cache.Set("OWM", v)
		}

		ch <- prom.MustNewConstMetric(c.temperature, prom.GaugeValue, v.Current.Values.Temperature, l.Location)
		ch <- prom.MustNewConstMetric(c.humidity, prom.GaugeValue, v.Current.Values.Humidity, l.Location)
		ch <- prom.MustNewConstMetric(c.feelslike, prom.GaugeValue, v.Current.Values.FeelsLike, l.Location)
		ch <- prom.MustNewConstMetric(c.pressure, prom.GaugeValue, v.Current.Values.Pressure, l.Location)
		ch <- prom.MustNewConstMetric(c.windspeed, prom.GaugeValue, v.Current.Values.WindSpeed, l.Location)
		ch <- prom.MustNewConstMetric(c.winddegree, prom.GaugeValue, v.Current.Values.WindDegree, l.Location)
		ch <- prom.MustNewConstMetric(c.windgust, prom.GaugeValue, v.Current.Values.WindGust, l.Location)
		ch <- prom.MustNewConstMetric(c.rain1h, prom.GaugeValue, v.Current.Values.Rain1Hour, l.Location)
		ch <- prom.MustNewConstMetric(c.snow1h, prom.GaugeValue, v.Current.Values.Snow1Hour, l.Location)
		ch <- prom.MustNewConstMetric(c.cloudiness, prom.GaugeValue, float64(v.Current.Values.Clouds), l.Location)
		ch <- prom.MustNewConstMetric(c.sunrise, prom.GaugeValue, float64(v.Current.Values.Sunrise), l.Location)
		ch <- prom.MustNewConstMetric(c.sunset, prom.GaugeValue, float64(v.Current.Values.Sunset), l.Location)
		ch <- prom.MustNewConstMetric(c.dewpoint, prom.GaugeValue, v.Current.Values.DewPoint, l.Location)
		ch <- prom.MustNewConstMetric(c.uvi, prom.GaugeValue, v.Current.Values.UVI, l.Location)
		if len(v.Pollution.List) > 0 {
			ch <- prom.MustNewConstMetric(c.aqi, prom.GaugeValue, float64(v.Pollution.List[0].Main.AQI), l.Location)
			ch <- prom.MustNewConstMetric(c.co, prom.GaugeValue, v.Pollution.List[0].Components.CO, l.Location)
			ch <- prom.MustNewConstMetric(c.no, prom.GaugeValue, v.Pollution.List[0].Components.NO, l.Location)
			ch <- prom.MustNewConstMetric(c.no2, prom.GaugeValue, v.Pollution.List[0].Components.NO2, l.Location)
			ch <- prom.MustNewConstMetric(c.o3, prom.GaugeValue, v.Pollution.List[0].Components.O3, l.Location)
			ch <- prom.MustNewConstMetric(c.so2, prom.GaugeValue, v.Pollution.List[0].Components.SO2, l.Location)
			ch <- prom.MustNewConstMetric(c.pm25, prom.GaugeValue, v.Pollution.List[0].Components.PM25, l.Location)
			ch <- prom.MustNewConstMetric(c.pm10, prom.GaugeValue, v.Pollution.List[0].Components.PM10, l.Location)
			ch <- prom.MustNewConstMetric(c.nh3, prom.GaugeValue, v.Pollution.List[0].Components.NH3, l.Location)
		}
	}
}
