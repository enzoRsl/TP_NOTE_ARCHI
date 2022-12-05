package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	dotenv "github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// type Airport struct {
// 	AirportIataCode string  `yaml:"airportIata"`
// 	Latitude        float64 `yaml:"latitude"`
// 	Longitude       float64 `yaml:"longitude"`
// }

type Airport struct {
	AirportIataCode string  `yaml:"airportIataCode"`
	Latitude        float64 `yaml:"latitude"`
	Longitude       float64 `yaml:"longitude"`
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// Used for full precision at atomic level
func GetCurrentDateAndHour() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

func GetCurrentHour() string {
	return time.Now().Format("15")
}

// InitEnvVariables loads the environment variables from the .env file
func InitEnvVariables() {
	err := dotenv.Load()
	if err != nil {
		// we try to load the .env file with the relative path
		err = dotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// GetCliParams returns the command line arguments given to the program
func GetCliParams() (string, Airport) {
	// We ensure that the sensor id is given as an argument
	// The reason is that we could have multiple sensors
	if len(os.Args) != 3 {
		panic("You must give the sensor id as a command line argument, as well as the airport" +
			" 3 letters IATA code (example: NTE)" +
			"\nUsage: ./temperature <sensor_id> <IATA code>" +
			"\n./Example (for the NTE airport): temperature 1 NTE")
	}

	sensorId := os.Args[1]
	airportIataCode := os.Args[2]

	airport, err := GetAirportFromConfig("config.yml", airportIataCode)
	if err != nil {
		log.Fatal(err)
	}

	return sensorId, airport
}

// GetAirportFromConfig retrieves the airport from the config file
func GetAirportFromConfig(configPath string, airportIataCode string) (Airport, error) {
	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var airports []Airport
	err = yaml.Unmarshal(configFile, &airports)
	if err != nil {
		panic(err)
	}

	// Make the airports variable accessible via the airport Iata code
	for _, airport := range airports {
		if airport.AirportIataCode == airportIataCode {
			return airport, nil
		}
	}

	return Airport{}, fmt.Errorf("Airport with IATA code %s not found in the config.yml", airportIataCode)
}

// GetPublicationInterval retrieves the publication interval from the .env file (given in seconds)
func GetPublicationIntervalInSeconds() time.Duration {
	publicationInterval, errConv := strconv.Atoi(os.Getenv("PUB_INTERVAL_SECONDS"))
	if errConv != nil {
		log.Fatal(errConv)
	}
	return time.Duration(publicationInterval)
}
