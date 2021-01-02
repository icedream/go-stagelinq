package main

import (
	"log"
	"time"

	"git.icedream.tech/icedream/stagelinq-receiver"
)

const (
	appName    = "Icedream StagelinQ Receiver"
	appVersion = "0.0.0"
	timeout    = 5 * time.Second
)

func main() {
	listener, err := stagelinq.ListenWithConfiguration(&stagelinq.ListenerConfiguration{
		DiscoveryTimeout: timeout,
		SoftwareName:     appName,
		SoftwareVersion:  appVersion,
		Name:             "testing",
	})
	if err != nil {
		panic(err)
	}

	deadline := time.After(timeout)
	foundDevices := []*stagelinq.Device{}

	log.Printf("Listening for devices for %s", timeout)

discoveryLoop:
	for {
		select {
		case <-deadline:
			break discoveryLoop
		default:
			device, deviceState, err := listener.Discover()
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
				continue discoveryLoop
			}
			// ignore device leaving messages since we do a one-off list
			if deviceState != stagelinq.DevicePresent {
				continue discoveryLoop
			}
			// check if we already found this device before
			for _, foundDevice := range foundDevices {
				if foundDevice.IsEqual(device) {
					continue discoveryLoop
				}
			}
			foundDevices = append(foundDevices, device)
			log.Printf("%s %q %q %q", device.IP.String(), device.Name, device.SoftwareName, device.SoftwareVersion)
		}
	}

	log.Printf("Found devices: %d", len(foundDevices))
}
