// Package opendatameteo is a package that allows the application
// to retrieve meteo data from the open-meteo.com/ API.
// For the scope of this project, we only need the current temperature,
// the current wind speed and the current air pressure.
//
// The API is updated every hour or so, so we can't get the current
// precise temperature, wind speed and air pressure. This is the
// reason why the results are so stepped.
package opendatameteo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// This methods returns the temperature at a given airport for a given date.
func GetTemperature(latitude float64, longitude float64, date string) float64 {
	meteo := getMeteoFromApi(latitude, longitude, date)
	current_weather := meteo["current_weather"].(map[string]interface{})

	return current_weather["temperature"].(float64)
}

// This methods returns the wind speed at a given airport for a given date.
func GetWindSpeed(latitude float64, longitude float64, date string) float64 {
	meteo := getMeteoFromApi(latitude, longitude, date)
	current_weather := meteo["current_weather"].(map[string]interface{})

	return current_weather["windspeed"].(float64)
}

// This methods returns the air pressure at a given airport for a given date.
func GetAirPressure(latitude float64, longitude float64, date string) float64 {
	meteo := getMeteoFromApi(latitude, longitude, date)
	current_weather := meteo["current_weather"].(map[string]interface{})

	// Getting the current air pressure, as it not given in the current_weather object
	current_time := current_weather["time"].(string)
	hour, errConv := strconv.Atoi(current_time[11:13])
	if errConv != nil {
		log.Fatal(errConv)
	}

	hourly_data := meteo["hourly"].(map[string]interface{}) //.([]interface{})
	hourly_pressure_data := hourly_data["surface_pressure"].([]interface{})

	return hourly_pressure_data[hour].(float64)
}

// This method returns the meteo data from the open-meteo.com API
// for a given airport and a given date.
func getMeteoFromApi(latitude float64, longitude float64, date string) map[string]interface{} {
	url := buildApiUrl(latitude, longitude, date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	jsonOut := map[string]interface{}{}
	jsonErr := json.Unmarshal(body, &jsonOut)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}

	return jsonOut
}

// This method builds the API url to call the open-meteo.com API
// for a given airport and a given date.
func buildApiUrl(latitude float64, longitude float64, date string) string {
	return "https://api.open-meteo.com/v1/forecast?" +
		"latitude=" + fmt.Sprintf("%.3f", latitude) +
		"&longitude=" + fmt.Sprintf("%.3f", longitude) +
		"&hourly=temperature_2m,surface_pressure,windspeed_80m" +
		"&current_weather=true" +
		"&timezone=Europe%2FBerlin" +
		"&start_date=" + date +
		"&end_date=" + date
}
