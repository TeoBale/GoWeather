package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"io"
	"net/http"
	"time"
)

const WEATHER_API = "https://api.open-meteo.com/v1/forecast?latitude=45.4918&longitude=9.2981&hourly=temperature_2m,relative_humidity_2m,precipitation_probability,precipitation,weather_code,cloud_cover,wind_speed_10m&timeformat=unixtime&forecast_days=3"

type Weather struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time                     string `json:"time"`
		Temperature2M            string `json:"temperature_2m"`
		RelativeHumidity2M       string `json:"relative_humidity_2m"`
		PrecipitationProbability string `json:"precipitation_probability"`
		Precipitation            string `json:"precipitation"`
		WeatherCode              string `json:"weather_code"`
		CloudCover               string `json:"cloud_cover"`
		WindSpeed10M             string `json:"wind_speed_10m"`
	} `json:"hourly_units"`
	Hourly struct {
		Time                     []int     `json:"time"`
		Temperature2M            []float64 `json:"temperature_2m"`
		RelativeHumidity2M       []int     `json:"relative_humidity_2m"`
		PrecipitationProbability []int     `json:"precipitation_probability"`
		Precipitation            []float32 `json:"precipitation"`
		WeatherCode              []int     `json:"weather_code"`
		CloudCover               []int     `json:"cloud_cover"`
		WindSpeed10M             []float32 `json:"wind_speed_10m"`
	} `json:"hourly"`
}

func main() {
	res, err := http.Get(WEATHER_API)
	if err != nil {
		fmt.Print(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather api not available !")
	}

	out, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(out, &weather)
	if err != nil {
		panic(err)
	}

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "TIME"},
			{Align: simpletable.AlignCenter, Text: "TEMP"},
			{Align: simpletable.AlignCenter, Text: "HUMIDITY"},
			{Align: simpletable.AlignCenter, Text: "PREC %"},
			{Align: simpletable.AlignCenter, Text: "PREC mm"},
			{Align: simpletable.AlignCenter, Text: "CODE"},
			{Align: simpletable.AlignCenter, Text: "CLOUDE"},
			{Align: simpletable.AlignCenter, Text: "WIND"},
		},
	}

	for i, timestamp := range weather.Hourly.Time {
		date := time.Unix(int64(timestamp), 0).Format("15:04")

		r := []*simpletable.Cell{
			{Text: fmt.Sprintf("%v", date)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", weather.Hourly.Temperature2M[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v°", weather.Hourly.RelativeHumidity2M[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v %%", weather.Hourly.PrecipitationProbability[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%vmm", weather.Hourly.Precipitation[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", weather.Hourly.WeatherCode[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v %%", weather.Hourly.CloudCover[i])},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", weather.Hourly.WindSpeed10M[i])},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleUnicode)
	fmt.Println(table.String())

}
