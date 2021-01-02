package owm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func buildUrl(cityId int, token string, language string) string {
	return fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&lang=%s&units=metric", cityId, token, language)
}

type MainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	Humidity  float64 `json:"humidity"`
}

type WindData struct {
	Speed     float64 `json:"speed"`
	Gust      float64 `json:"gust,omitempty"`
	Direction float64 `json:"deg"`
}

type CloudsData struct {
	Percentage float64 `json:"all"`
}

type RainData struct {
	Last1h float64 `json:"1h"`
	Last3h float64 `json:"3h"`
}

type SnowData struct {
	Last1h float64 `json:"1h"`
	Last3h float64 `json:"3h"`
}

type WeatherData struct {
	Main       MainData   `json:"main"`
	Visibility float64    `json:"visibility"`
	Wind       WindData   `json:"wind"`
	Clouds     CloudsData `json:"clouds"`
	Rain       RainData   `json:"rain,omitempty"`
	Snow       SnowData   `json:"snow,omitempty"`
}

func GetOwmData(cityId int, token string) (*WeatherData, error) {
	var weatherData WeatherData

	res, err := http.Get(buildUrl(cityId, token, "de"))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &weatherData)
	if err != nil {
		return nil, err
	}

	return &weatherData, nil
}
