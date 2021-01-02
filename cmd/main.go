package main

import (
	"fmt"
	"github.com/henninge/owm-exporter/metrics"
	"github.com/henninge/owm-exporter/owm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CityMetrics struct {
	Id              int
	Name            string
	ApiToken        string
	IntervalMinutes time.Duration
	Metrics         *metrics.WeatherMetrics
}

func NewCityMetricsFromEnv() (*CityMetrics, error) {
	cityIdString := os.Getenv("METEO_CITY_ID")
	if cityIdString == "" {
		return nil, fmt.Errorf("Missing environment variable: METEO_CITY_ID")
	}
	cityId, err := strconv.Atoi(cityIdString)
	if err != nil {
		return nil, fmt.Errorf("Invalid number for city ID: %s, %v ", cityIdString, err)
	}

	cityName := os.Getenv("METEO_CITY_NAME")
	if cityIdString == "" {
		return nil, fmt.Errorf("Missing environment variable: METEO_CITY_NAME")
	}

	token := os.Getenv("METEO_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("Missing environment variable: METEO_API_TOKEN")
	}

	intervalMinutes := 5
	intervalMinutesString := os.Getenv("METEO_INTERVAL_MINUTES")
	if intervalMinutesString != "" {
		intervalMinutes, err = strconv.Atoi(intervalMinutesString)
		if err != nil {
			return nil, fmt.Errorf("Invalid Number for interval: %s, %v ", intervalMinutesString, err)
		}
	}

	return &CityMetrics{
		Id:              cityId,
		Name:            cityName,
		ApiToken:        token,
		Metrics:         metrics.NewWeatherMetrics(cityId, cityName),
		IntervalMinutes: time.Duration(intervalMinutes) * time.Minute,
	}, nil
}

func (city *CityMetrics) updateMetrics() {
	weatherData, err := owm.GetOwmData(city.Id, city.ApiToken)
	if err != nil {
		log.Errorf("Error during data retrieval for %s (%d): %v\n", city.Name, city.Id, err)
	} else {
		log.Infof("Received weather data for %s (%d), updating.", city.Name, city.Id)
		city.Metrics.Update(weatherData)
	}
}

func (city *CityMetrics) RecordMetrics() {
	ticker := time.NewTicker(city.IntervalMinutes)
	go func() {
		// Initial data retrieval
		city.updateMetrics()
		for {
			select {
			case <-ticker.C:
				city.updateMetrics()
			}
		}
	}()
}

func main() {
	city, err := NewCityMetricsFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	city.RecordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Fatal(err)
	}
}
