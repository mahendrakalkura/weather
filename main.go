package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Current struct {
	ApparentTemperature float64 `json:"apparent_temperature"`
	Precipitation       float64 `json:"precipitation"`
	RelativeHumidity    int     `json:"relative_humidity_2m"`
	Temperature         float64 `json:"temperature_2m"`
	Time                string  `json:"time"`
	WeatherCode         int     `json:"weather_code"`
}

type Daily struct {
	PrecipitationProbabilityMax []int     `json:"precipitation_probability_max"`
	PrecipitationSum            []float64 `json:"precipitation_sum"`
	Sunrise                     []string  `json:"sunrise"`
	Sunset                      []string  `json:"sunset"`
	TemperatureMax              []float64 `json:"temperature_2m_max"`
	TemperatureMean             []float64 `json:"temperature_2m_mean"`
	TemperatureMin              []float64 `json:"temperature_2m_min"`
	Time                        []string  `json:"time"`
	WeatherCode                 []int     `json:"weather_code"`
}

type Forecast struct {
	Current Current `json:"current"`
	Daily   Daily   `json:"daily"`
	Hourly  Hourly  `json:"hourly"`
}

type Hourly struct {
	ApparentTemperature      []float64 `json:"apparent_temperature"`
	Precipitation            []float64 `json:"precipitation"`
	PrecipitationProbability []int     `json:"precipitation_probability"`
	RelativeHumidity         []int     `json:"relative_humidity_2m"`
	Temperature              []float64 `json:"temperature_2m"`
	Time                     []string  `json:"time"`
	WeatherCode              []int     `json:"weather_code"`
}

const (
	apiURL     = "https://api.open-meteo.com/v1/forecast"
	latitude   = "13.3409"
	longitude  = "74.7421"
	timeLayout = "2006-01-02T15:04"
)

var weatherCodes = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

func clock(value string) string {
	parsed, err := time.Parse(timeLayout, value)
	if err != nil {
		return value
	}
	return parsed.Format("15:04")
}

func fetchForecast() (Forecast, error) {
	forecast := Forecast{}
	query := fmt.Sprintf("%s?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,weather_code&hourly=temperature_2m,apparent_temperature,relative_humidity_2m,precipitation,precipitation_probability,weather_code&daily=weather_code,temperature_2m_max,temperature_2m_min,temperature_2m_mean,precipitation_probability_max,precipitation_sum,sunrise,sunset&timezone=auto&forecast_days=7", apiURL, latitude, longitude)

	response, err := http.Get(query)
	if err != nil {
		return forecast, err
	}
	defer func() {
		e := response.Body.Close()
		if e != nil {
			fmt.Fprintln(os.Stderr, "close error:", e)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return forecast, fmt.Errorf("api returned status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return forecast, err
	}

	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return forecast, err
	}
	return forecast, nil
}

func hourlyStart(hourly Hourly, currentTime string) int {
	now, err := time.Parse(timeLayout, currentTime)
	if err != nil {
		return 0
	}
	for i, value := range hourly.Time {
		parsed, e := time.Parse(timeLayout, value)
		if e != nil {
			continue
		}
		if !parsed.Before(now) {
			return i
		}
	}
	return 0
}

func main() {
	forecast, err := fetchForecast()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	rain := 0
	index := hourlyStart(forecast.Hourly, forecast.Current.Time)
	if index < len(forecast.Hourly.PrecipitationProbability) {
		rain = forecast.Hourly.PrecipitationProbability[index]
	}

	printCurrent(forecast.Current, rain)
	printHourly(forecast.Hourly, forecast.Current.Time)
	printDaily(forecast.Daily)
}

func printCurrent(current Current, rain int) {
	rows := [][2]string{
		{"Condition", weather(current.WeatherCode)},
		{"Temperature", fmt.Sprintf("%.1f C", current.Temperature)},
		{"Feels like", fmt.Sprintf("%.1f C", current.ApparentTemperature)},
		{"Humidity", fmt.Sprintf("%d %%", current.RelativeHumidity)},
		{"Rain", fmt.Sprintf("%d %%", rain)},
		{"Precipitation", fmt.Sprintf("%.1f mm", current.Precipitation)},
	}
	fieldWidth := len("Field")
	valueWidth := len("Value")
	for _, row := range rows {
		fieldWidth = max(fieldWidth, len(row[0]))
		valueWidth = max(valueWidth, len(row[1]))
	}
	border := fmt.Sprintf("+-%s-+-%s-+", strings.Repeat("-", fieldWidth), strings.Repeat("-", valueWidth))

	fmt.Println("CURRENT")
	fmt.Println(border)
	fmt.Printf("| %-*s | %-*s |\n", fieldWidth, "Field", valueWidth, "Value")
	fmt.Println(border)
	for _, row := range rows {
		fmt.Printf("| %-*s | %-*s |\n", fieldWidth, row[0], valueWidth, row[1])
	}
	fmt.Println(border)
	fmt.Println()
}

func printDaily(daily Daily) {
	conditionWidth := len("Condition")
	for _, code := range daily.WeatherCode {
		conditionWidth = max(conditionWidth, len(weather(code)))
	}
	border := fmt.Sprintf("+-------------+-%s-+--------+--------+--------+-------+---------------+---------+---------+", strings.Repeat("-", conditionWidth))

	fmt.Println("NEXT 7 DAYS")
	fmt.Println(border)
	fmt.Printf("| %-11s | %-*s | %6s | %6s | %6s | %5s | %13s | %7s | %7s |\n", "Date", conditionWidth, "Condition", "High", "Mean", "Low", "Rain", "Precipitation", "Sunrise", "Sunset")
	fmt.Println(border)
	for i := range daily.Time {
		date := daily.Time[i]
		parsed, err := time.Parse("2006-01-02", date)
		if err == nil {
			date = parsed.Format("Mon Jan 02")
		}
		fmt.Printf("| %-11s | %-*s | %5.1fC | %5.1fC | %5.1fC | %4d%% | %10.1f mm | %7s | %7s |\n",
			date,
			conditionWidth, weather(daily.WeatherCode[i]),
			daily.TemperatureMax[i],
			daily.TemperatureMean[i],
			daily.TemperatureMin[i],
			daily.PrecipitationProbabilityMax[i],
			daily.PrecipitationSum[i],
			clock(daily.Sunrise[i]),
			clock(daily.Sunset[i]))
	}
	fmt.Println(border)
}

func printHourly(hourly Hourly, currentTime string) {
	start := hourlyStart(hourly, currentTime)
	end := min(start+24, len(hourly.Time))
	conditionWidth := len("Condition")
	for i := start; i < end; i++ {
		conditionWidth = max(conditionWidth, len(weather(hourly.WeatherCode[i])))
	}
	border := fmt.Sprintf("+-----------+-%s-+-------------+------------+----------+-------+---------------+", strings.Repeat("-", conditionWidth))

	fmt.Println("NEXT 24 HOURS")
	fmt.Println(border)
	fmt.Printf("| %-9s | %-*s | %11s | %10s | %8s | %5s | %13s |\n", "Time", conditionWidth, "Condition", "Temperature", "Feels like", "Humidity", "Rain", "Precipitation")
	fmt.Println(border)
	for i := start; i < end; i++ {
		label := hourly.Time[i]
		parsed, err := time.Parse(timeLayout, label)
		if err == nil {
			label = parsed.Format("Mon 15:04")
		}
		fmt.Printf("| %-9s | %-*s | %10.1fC | %9.1fC | %7d%% | %4d%% | %10.1f mm |\n",
			label,
			conditionWidth, weather(hourly.WeatherCode[i]),
			hourly.Temperature[i],
			hourly.ApparentTemperature[i],
			hourly.RelativeHumidity[i],
			hourly.PrecipitationProbability[i],
			hourly.Precipitation[i])
	}
	fmt.Println(border)
	fmt.Println()
}

func weather(code int) string {
	description, ok := weatherCodes[code]
	if !ok {
		return fmt.Sprintf("Code %d", code)
	}
	return description
}
