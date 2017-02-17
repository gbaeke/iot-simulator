package main

import "os"
import "encoding/json"

// define a struct for a device group (used in configuration struct)
type devgroup struct {
	Prefix    string // prefix for device name
	DeviceNum int    // number of devices with this prefix
	Firmware  string // firmware version like 1.20
	IoTHub    int    // index into IoTHubs slice in configuration struct

}

// define a struct for configuration
type configuration struct {
	Interval  int        // interval to send in seconds
	IoTHubs   []string   // slice of strings to hold IoTHubs to send to
	SasTokens []string   // slice of strings to hold SaS tokens
	DevGroups []devgroup // slice of devgroup
}

// Conf is global configuration variable
var Conf configuration

// read the configuration from configuration file
func getConf(configFile string) (configuration, error) {
	file, _ := os.Open(configFile)
	defer file.Close()

	decoder := json.NewDecoder(file)
	Conf := configuration{}
	err := decoder.Decode(&Conf)

	// in case of err Conf will be empty struct
	return Conf, err

}
