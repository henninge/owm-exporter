package metrics

import (
	"github.com/henninge/owm-exporter/owm"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type WeatherMetrics struct {
	IsRegistered bool
	CityId       int
	CityName     string
	Temp         prometheus.Gauge
	FeelsLike    prometheus.Gauge
	TempMin      prometheus.Gauge
	TempMax      prometheus.Gauge
	Pressure     prometheus.Gauge
	Humidity     prometheus.Gauge
	Visibility   prometheus.Gauge
	Cloudiness   prometheus.Gauge
	WindSpeed    prometheus.Gauge
	WindGusts    prometheus.Gauge
	WindBft      prometheus.Gauge
	WindGustsBft prometheus.Gauge
	WindDir      prometheus.Gauge
	Rain1h       prometheus.Gauge
	Rain3h       prometheus.Gauge
	Snow1h       prometheus.Gauge
	Snow3h       prometheus.Gauge
}

func NewWeatherMetrics(cityId int, cityName string) *WeatherMetrics {
	constLabels := prometheus.Labels{"city_id": strconv.Itoa(cityId), "city_name": cityName}
	return &WeatherMetrics{
		IsRegistered: false,
		CityId:       cityId,
		CityName:     cityName,
		Temp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_temp_c",
			Help:        "Temperature in degrees Celsius.",
			ConstLabels: constLabels,
		}),
		FeelsLike: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_feelslike_c",
			Help:        "Feels-like temperature in degrees Celsius.",
			ConstLabels: constLabels,
		}),
		TempMin: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_tempmin_c",
			Help:        "Minimum temperature in degrees Celsius.",
			ConstLabels: constLabels,
		}),
		TempMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_tempmax_c",
			Help:        "Maximum temperature in degrees Celsius.",
			ConstLabels: constLabels,
		}),
		Pressure: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_pressure_hpa",
			Help:        "Pressure in hectopascal.",
			ConstLabels: constLabels,
		}),
		Humidity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_humidity_percent",
			Help:        "Humidity in percent.",
			ConstLabels: constLabels,
		}),
		Visibility: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_visibilty_km",
			Help:        "Visibility in kilometer.",
			ConstLabels: constLabels,
		}),
		Cloudiness: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_cloudiness_percent",
			Help:        "Cloudiness in percent.",
			ConstLabels: constLabels,
		}),
		WindSpeed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_wind_speed_kn",
			Help:        "Wind speed in knots.",
			ConstLabels: constLabels,
		}),
		WindGusts: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_wind_gusts_kn",
			Help:        "Wind gusts in knots.",
			ConstLabels: constLabels,
		}),
		WindDir: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_wind_direction_deg",
			Help:        "Wind direction in degrees.",
			ConstLabels: constLabels,
		}),
		WindBft: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_wind_speed_bft",
			Help:        "Wind in Beaufort.",
			ConstLabels: constLabels,
		}),
		WindGustsBft: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_wind_gusts_bft",
			Help:        "Wind gusts in Beaufort.",
			ConstLabels: constLabels,
		}),
		Rain1h: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_rain_1h_mm",
			Help:        "Rain volume last 1 hour in mm.",
			ConstLabels: constLabels,
		}),
		Rain3h: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_rain_3h_mm",
			Help:        "Rain volume last 3 hour in mm.",
			ConstLabels: constLabels,
		}),
		Snow1h: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_snow_1h_mm",
			Help:        "Snow volume last 1 hour in mm.",
			ConstLabels: constLabels,
		}),
		Snow3h: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "meteo_snow_3h_mm",
			Help:        "Snow volume last 3 hour in mm.",
			ConstLabels: constLabels,
		}),
	}
}

var bftSteps = []float64{
	1.0, 4.0, 7.0, 11.0, 17.0, 22.0, 28.0, 34.0, 41.0, 48.0, 56.0, 64.0,
}

func knotsToBft(knots float64) float64 {
	for bft, step := range bftSteps {
		if knots < step {
			return float64(bft)
		}
	}
	return 12.0
}

func (wm *WeatherMetrics) Register() {
	if !wm.IsRegistered {
		prometheus.MustRegister(
			wm.Temp,
			wm.FeelsLike,
			wm.TempMin,
			wm.TempMax,
			wm.Pressure,
			wm.Humidity,
			wm.Visibility,
			wm.Cloudiness,
			wm.WindSpeed,
			wm.WindGusts,
			wm.WindDir,
			wm.WindBft,
			wm.WindGustsBft,
			wm.Rain1h,
			wm.Rain3h,
			wm.Snow1h,
			wm.Snow3h,
		)
		wm.IsRegistered = true
	}
}

func (wm *WeatherMetrics) Update(data *owm.WeatherData) {
	wm.Temp.Set(data.Main.Temp)
	wm.FeelsLike.Set(data.Main.FeelsLike)
	wm.TempMin.Set(data.Main.TempMin)
	wm.TempMax.Set(data.Main.TempMax)
	wm.Pressure.Set(data.Main.Pressure)
	wm.Humidity.Set(data.Main.Humidity)
	wm.Visibility.Set(data.Visibility / 1000.0) // Convert m->km
	wm.Cloudiness.Set(data.Clouds.Percentage)

	wind_knots := data.Wind.Speed * 1.943844 // Convert m/s->kn
	gusts_knots := data.Wind.Gust * 1.943844 // Convert m/s->kn
	wm.WindSpeed.Set(wind_knots)
	wm.WindGusts.Set(gusts_knots)
	wm.WindDir.Set(data.Wind.Direction)
	wm.WindBft.Set(knotsToBft(wind_knots))
	wm.WindGustsBft.Set(knotsToBft(gusts_knots))

	wm.Rain1h.Set(data.Rain.Last1h)
	wm.Rain3h.Set(data.Rain.Last3h)
	wm.Snow1h.Set(data.Snow.Last1h)
	wm.Snow3h.Set(data.Snow.Last3h)

	// Registration is delayed until the first data is available
	// to avoid returning zero values.
	wm.Register()
}
