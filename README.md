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
| `OW_CITY`                     | `city`                      | `Kitchener, CA`              | City/Location in which to gather weather metrics. Separate multiple locations with \| for example "Kitchener, CA\|Seattle, WA" |
| `OW_DEGREES_UNIT`             | `degrees-unit`              | `C`                         | Unit in which to show metrics (Kelvin, Fahrenheit or Celsius) |
| `OW_LANGUAGE`                 | `language`                  | `EN`                        | Language in which to show metrics |

## Quickstart

The fastest way to get started is to run `docker-compose.yml` which is included in this repository. Change the `CHANGEME` to your OWM token and run `docker-compose up --build` to get things up and running! The metric will be accessible on the `/metrics` endpoint.

```yaml
version: '3.2'
services:
  owm:
    container_name: openweather-exporter
    build: .
    ports:
    - "9091:9091"
    environment:
    - OW_LISTEN_ADDRESS=0.0.0.0:9091
    - OW_CITY=Kitchener, CA
    - OW_APIKEY=CHANGEME
    - OW_DEGREES_UNIT=C
    restart: unless-stopped
```