# OWM Exporter

Prometheus exporter for current weather data from [openweathermap.org](https://openweathermap.org/current).

## Building

    go build -o owm_exporter cmd/main.go

## Configuration

Configuration is done via environment variables

 * `METEO_CITY_ID` is the OWM id of the city for which to export weather data.
 * `METEO_CITY_NAME` is the name of the city.
 * `METEO_API_TOKEN` is the OWM API token to use for retrieval requests.
 * `METEO_INTERVAL_MINUTES` (optional) is the retrieval interval (default: 5).

`METEO_CITY_ID` and `METEO_CITY_NAME` are exported as `city_id` and `city_name` labels on each metric.
