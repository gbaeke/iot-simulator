package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

// define a flag -r to remove all devices
var removeDevices = flag.Bool("r", false, "Remove devices from IoT Hub")

// http client
var client *http.Client

func main() {
	// logical processors for goroutines
	runtime.GOMAXPROCS(runtime.NumCPU())

	// create a channel to allow goroutines to communicate
	// in this case goroutines will communicate with main() which is also a goroutine
	ch := make(chan string)

	// http client
	client = &http.Client{
		Timeout: time.Second * 10,
	}

	// parse flags
	flag.Parse()

	// read configuration into conf struct
	config, err := getConf("config.json")

	// store config in global Conf variable
	Conf = config

	if err != nil {
		// could be due to file not found or some JSON parsing error
		fmt.Println("Error occurred: ", err)
		os.Exit(1)
	}

	// get device list (returns pointer) created from the configuration
	devList := getDeviceList()

	if *removeDevices {
		for _, device := range *devList {
			deleteResp, err := device.deleteDevice()
			if err != nil {
				fmt.Printf("Could not delete device %s. Error: %v\n", device.Name, err)
			} else {
				fmt.Printf("Deleted device %s with reponse code %d\n", device.Name, deleteResp.StatusCode)
			}

		}
		os.Exit(0)
	}

	for _, device := range *devList {
		// for each device, run the deviceSend method as a goroutine
		go device.deviceSend(Conf.Interval, ch)
	}

	// not sure this is good practice but we keep reading from
	// the channel forever to pick up messages from the goroutines
	for {
		fmt.Println(<-ch)
	}

}
