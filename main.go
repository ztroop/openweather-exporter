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

package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ztroop/openweather-exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("openweather-exporter", "Openweather Exporter for Openweather API").Author("Zackary Troop")
	addr        = app.Flag("listen-address", "HTTP port to listen on").Envar("OW_LISTEN_ADDRESS").Default(":9091").String()
	apiKey      = app.Flag("apikey", "Openweather API Key").Envar("OW_APIKEY").Required().String()
	city        = app.Flag("city", "City for Openweather to gather metrics from.").Envar("OW_CITY").Default("Toronto, ON").String()
	degreesUnit = app.Flag("degrees-unit", "The base unit for temperature output. Fahrenheit or Celsius").Envar("OW_DEGREES_UNIT").Default("C").String()
	language    = app.Flag("language", "The language for metric output").Envar("OW_LANGUAGE").Default("EN").String()
	cacheTTL    = app.Flag("cache-ttl", "Cache time-to-live in seconds").Envar("OW_CACHE_TTL").Default("300").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	cache := ttlcache.NewCache()
	ttl, err := strconv.ParseUint(*cacheTTL, 10, 64)
	if err != nil {
		log.Fatal("Invalid TTL value: ", err)
	}
	cache.SetTTL(time.Duration(ttl) * time.Second)
	cache.SkipTTLExtensionOnHit(true)

	weatherCollector := collector.NewOpenweatherCollector(*degreesUnit, *language, *apiKey, *city, cache)
	prometheus.MustRegister(weatherCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port " + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
