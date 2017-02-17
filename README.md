# iot-simulator
IoT Simulator in Go

Make a file called config.json and use the sample-config.json as a sample

{
    "Interval":5,
    "IoTHubs":["gebahub.azure-devices.net","gebahub2.azure-devices.net"],
    "SasTokens":["SharedAccessSignature sr=...",
                 "SharedAccessSignature sr=..."],
    "DevGroups":[
        {"Prefix":"deva","DeviceNum":50,"Firmware":"1.0","IoTHub": 0},
        {"Prefix":"devb","DeviceNum":50,"Firmware":"1.1","IoTHub": 1}]
}

With the above config, there are two IoT Hubs each with its own SAS token (generate it with Device Explorer). The DevGroups control device creation? In the above case, 50 devices are created on IoT Hub 0 (gebahub) and 50 on IoT Hub 2 (gebahub2). The device names are prefix+device number.

Just running the simulator creates the devices and sending starts. Dummy data is sent every 5 seconds as specified in the above config.

Run the emulator with -r to delete the devices based on the definition of the DevGroups.

Be aware that device operations are throttled at 100 operations per minute per unit unless you use the expensive tier 3.
