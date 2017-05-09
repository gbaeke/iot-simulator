package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type device struct {
	Name     string // name of the device (deviceprefix + number)
	Firmware string // firmware of device (results in other data to be sent)
	InHub    bool   // is the device in IoT Hub
	IoTHub   int    // index into IoTHubs slice in configuration object
}

type devicebody struct {
	DeviceID string
}

type devicemessage struct {
	Temperature float64
	Humidity    float64
}

// Function getDeviceList returns pointer to slice of device structs
// It takes the configuration read from a configuration JSON file to build the slice
func getDeviceList() *[]device {
	devices := make([]device, 0)

	for _, group := range Conf.DevGroups {

		// for current group, append devices to the list
		for i := 1; i <= group.DeviceNum; i++ {
			devices = append(devices, device{group.Prefix + strconv.Itoa(i), group.Firmware, false, group.IoTHub})
		}
	}

	return &devices
}

// function for use as a goroutine
func (d device) deviceSend(interval int, ch chan<- string) {
	// check if device exists
	if getResp, err := d.getDevice(); err == nil {
		if getResp.StatusCode == 404 {
			// create the device because it does not exist
			if createResp, err := d.createDevice(); err == nil {
				// device was created
				ch <- fmt.Sprint("Device created with response code ", createResp.StatusCode)
			} else {
				// there was an error creating the device
				ch <- fmt.Sprintf("Could not create device %s. Error %v", d.Name, err)
			}

		}
	} else {
		// there was an error calling getDevice
		ch <- fmt.Sprintf("Could not get device %s. Error %v", d.Name, err)
	}

	temperature := rand.Float64() * 40
	humidity := rand.Float64() * 40
	message := devicemessage{Temperature: temperature, Humidity: humidity}

	for {
		_, err := d.sendData(message)
		if err == nil {
			ch <- fmt.Sprintf("Sent message from %s", d.Name)
		} else {
			// there was a send error
			ch <- fmt.Sprintf("Error sending message from %s. Error %v", d.Name, err)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

}

func (d device) createDevice() (*http.Response, error) {

	reqBody, _ := json.Marshal(devicebody{DeviceID: d.Name})

	req, _ := http.NewRequest("PUT", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-02-03", bytes.NewBuffer(reqBody))

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])

	// do the request
	resp, err := client.Do(req)

	return resp, err
}

// get device in IoT Hub
// return 200 if device is found, 404 if not
func (d device) getDevice() (*http.Response, error) {

	req, _ := http.NewRequest("GET", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-02-03", nil)

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])

	// do the request
	resp, err := client.Do(req)

	return resp, err
}

func (d device) sendData(message devicemessage) (*http.Response, error) {

	reqBody, _ := json.Marshal(message)

	req, _ := http.NewRequest("POST", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"/messages/events?api-version=2016-02-03", bytes.NewBuffer(reqBody))

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])

	// do the request
	resp, err := client.Do(req)

	return resp, err

}

func (d device) deleteDevice() (*http.Response, error) {

	//fmt.Println("https://" + Conf.IoTHubs[d.IoTHub] + "/devices/" + d.Name + "?api-version=2016-11-14")

	req, _ := http.NewRequest("DELETE", "https://"+Conf.IoTHubs[d.IoTHub]+"/devices/"+d.Name+"?api-version=2016-11-14", nil)

	// add headers; DELETE requires If-Match with * for unconditional removal
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.SasTokens[d.IoTHub])
	req.Header.Add("If-Match", "*")

	// do the request
	resp, err := client.Do(req)

	return resp, err
}
