package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icedream/go-stagelinq"
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
	defer listener.Close()

	listener.AnnounceEvery(time.Second)

	deadline := time.After(timeout)
	foundDevices := []*stagelinq.Device{}

	log.Printf("Listening for devices for %s", timeout)

discoveryLoop:
	for {
		select {
		case <-deadline:
			break discoveryLoop
		default:
			device, deviceState, err := listener.Discover(timeout)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
				continue discoveryLoop
			}
			if device == nil {
				continue
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

			// discover provided services
			log.Println("\tattempting to connect to this device...")
			deviceConn, err := device.Connect(listener.Token(), []*stagelinq.Service{})
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			} else {
				defer deviceConn.Close()
				log.Println("\trequesting device data services…")
				services, err := deviceConn.RequestServices()
				if err != nil {
					log.Printf("WARNING: %s", err.Error())
					continue
				}

				for _, service := range services {
					log.Printf("\toffers %s at port %d", service.Name, service.Port)
					switch service.Name {
					case "BeatInfo":
						log.Println("\t\tconnecting to BeatInfo...")
						beatInfoTCPConn, err := device.Dial(service.Port)
						if err != nil {
							log.Printf("WARNING: %s", err.Error())
							continue
						}
						beatInfoConn, err := stagelinq.NewBeatInfoConnection(beatInfoTCPConn, listener.Token())
						if err != nil {
							log.Printf("WARNING: %s", err.Error())
							continue
						}

						log.Println("\t\trequesting start BeatInfo stream... PRESS CTRL-C TO ABORT!")
						abortC := make(chan os.Signal, 1)
						signal.Notify(abortC, os.Interrupt, syscall.SIGTERM)
						beatInfoConn.StartStream()

					beatInfoLoop:
						for {
							select {
							case bi := <-beatInfoConn.BeatInfoC():
								log.Printf("\t\t\t%+v", bi)
							case <-abortC:
								beatInfoConn.StopStream()
								break beatInfoLoop
							case err := <-beatInfoConn.ErrorC():
								log.Printf("WARNING: %s", err.Error())
								break beatInfoLoop
							}
						}
						beatInfoTCPConn.Close()
					}
				}

				log.Println("\tend of list of device data services")
			}
		}
	}

	log.Printf("Found devices: %d", len(foundDevices))
}
