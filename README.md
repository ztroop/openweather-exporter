# OpenWeather Exporter

[![build](https://github.com/ztroop/openweather-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/ztroop/openweather-exporter/actions/workflows/build.yml)
[![pulls](https://img.shields.io/docker/pulls/ztroop/openweather-exporter)](https://hub.docker.com/r/ztroop/openweather-exporter)


Prometheus exporter for [OpenWeather API V2](https://openweathermap.org/api), supporting `OneCall` and `Pollution` endpoints.

## Configuration

OpenWeather exporter can be controlled by both ENV or CLI flags as described below.

| Environment        	       | CLI (`--flag`)              | Default                 	    | Description                                                                                                      |
|----------------------------|-----------------------------|---------------------------- |------------------------------------------------------------------------------------------------------------------|
| `OW_LISTEN_ADDRESS`           | `listen-address`            | `:9091`                     | The port for /metrics to listen on |
| `OW_APIKEY`                   | `apikey`                    | `<REQUIRED>`                | Your Openweather API key |
| `OW_CITY`                     | `city`                      | `Toronto, ON`              | City/Location in which to gather weather metrics. Separate multiple locations with \| for example "Toronto, NY\|Seattle, WA" |
| `OW_DEGREES_UNIT`             | `degrees-unit`              | `C`                         | Unit in which to show metrics (Kelvin, Fahrenheit or Celsius) |
| `OW_LANGUAGE`                 | `language`                  | `EN`                        | Language in which to show metrics |
