package configurations

import (
	"os"
	"strings"
	"strconv"
	"log"
	"fmt"
	"errors"
)

type Config struct {
	Port       int        `json:"port"`
	MathApiUrl string     `json:"mathApiUrl"`
	Cors       CorsConfig `json:"cors"`
}

type CorsConfig struct {
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedHeaders   []string `json:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
}

const PortEnvVar = "PORT"
const MathApiUrlEnvVar = "MATH_API_URL"
const CorsEnvVar = "CORS"

func LoadConfig() (config *Config, err error) {
	log.Printf("Loading config from environment variables...\n")
	config = &Config{}

	envVarMissing := func(envVar string) (err error) {
		return errors.New(fmt.Sprintf("Environment variable \"%v\" not provided!\n", envVar))
	}

	if value, present := os.LookupEnv(PortEnvVar); present {
		port, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Environment variable \"%v\" is invalid! Must be an integer, but got \"%v\".\n", PortEnvVar, value))
		}
		config.Port = int(port)
	} else {
		return nil, envVarMissing(PortEnvVar)
	}

	if value, present := os.LookupEnv(MathApiUrlEnvVar); present {
		config.MathApiUrl = value
	} else {
		return nil, envVarMissing(PortEnvVar)
	}

	if value, present := os.LookupEnv(CorsEnvVar); present {
		corsConfigurationValues := strings.Split(value, ";")
		if len(corsConfigurationValues) != 4 {
			return nil, errors.New(fmt.Sprintf("Environment variable \"%v\" is invalid! Expected 4 values seperated by \";\", but got \"%v\"", CorsEnvVar, len(corsConfigurationValues)))
		}

		config.Cors.AllowedOrigins = strings.Split(corsConfigurationValues[0], ",")

		config.Cors.AllowedMethods = strings.Split(corsConfigurationValues[1], ",")

		config.Cors.AllowedHeaders = strings.Split(corsConfigurationValues[2], ",")

		allowCredentials, err := strconv.ParseBool(corsConfigurationValues[3])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Environment variable \"%v\" is invalid! Last sgment must be a boolean, but got \"%v\".\n", CorsEnvVar, corsConfigurationValues[3]))
		}
		config.Cors.AllowCredentials = allowCredentials
	} else {
		return nil, envVarMissing(CorsEnvVar)
	}

	log.Printf("Successfully loaded config.\n")
	return config, nil
}
