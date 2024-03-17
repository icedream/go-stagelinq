package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/icedream/go-stagelinq/eaas"
	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
	"github.com/rivo/tview"
)

const (
	appName    = "Icedream StagelinQ Receiver"
	appVersion = "0.0.0"
	timeout    = 15 * time.Second
)

var (
	grpcURL  string
	hostname string
	identity string
)

func init() {
	flag.StringVar(&grpcURL, "server", "", "GRPC URL of the remote Engine Library to connect to. If empty, will discover devices instead.")
	flag.Parse()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "eaas-demo"
	}
}

type App struct {
	*tview.Application
}

func main() {
	if len(grpcURL) == 0 {
		runDiscovery()
		return
	}

	runEngineLibraryUI(grpcURL)
}

func marshalJSON(v any) []byte {
	s, _ := json.Marshal(v)
	return s
}

func runEngineLibraryUI(grpcURL string) {
	// load our identity so we don't have to repeatedly re-verify
	identity, err := loadUUIDKey()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	connection, err := eaas.DialContext(ctx, grpcURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for approval on the other end...")
	resp, err := connection.CreateTrust(ctx, &networktrust.CreateTrustRequest{
		DeviceName: &hostname,
		// I honestly don't know why in the proto it was defined as "Ed25519Pk"...
		Ed25519Pk: &identity,
		// ...or why there even is a WireguardPort field, too?!
	})
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case resp.GetGranted() != nil:
		log.Println("Access granted")
	case resp.GetBusy() != nil:
		log.Fatal("Busy")
	case resp.GetDenied() != nil:
		log.Fatal("Access denied")
	default:
		panic("unexpected response")
	}

	getLibraryResp, err := connection.GetLibrary(ctx, &enginelibrary.GetLibraryRequest{})
	if err != nil {
		panic(err)
	}
	var pageSize uint32 = 25
	getTracksResp, err := connection.GetTracks(ctx, &enginelibrary.GetTracksRequest{
		PageSize: &pageSize,
	})
	if err != nil {
		panic(err)
	}
	for _, track := range getTracksResp.GetTracks() {
		log.Printf("Track: %s", string(marshalJSON(track)))
		getTrackResp, err := connection.GetTrack(ctx, &enginelibrary.GetTrackRequest{
			TrackId: track.GetMetadata().Id,
		})
		if err != nil {
			log.Println("\tfailed to GetTrack on this track")
			continue
		}
		log.Printf("\t%+v", getTrackResp)
	}
	for _, playlist := range getLibraryResp.GetPlaylists() {
		log.Printf("Playlist: %s", string(marshalJSON(playlist)))
		getTracksResp, err := connection.GetTracks(ctx, &enginelibrary.GetTracksRequest{
			PlaylistId: playlist.Id,
			PageSize:   &pageSize,
		})
		if errors.Is(err, io.EOF) {
			// BUG - empty playlist causes EOF, reconnect
			connection, err = eaas.DialContext(ctx, grpcURL)
			if err != nil {
				panic(err)
			}
		}
		if err != nil {
			panic(err)
		}
		for _, track := range getTracksResp.GetTracks() {
			log.Printf("\tTrack: ID %s", track.GetMetadata().GetId())
		}
	}
}

func runDiscovery() {
	discoverer, err := eaas.NewDiscovererWithConfiguration(&eaas.DiscovererConfiguration{
		DiscoveryTimeout: timeout,
	})
	if err != nil {
		panic(err)
	}
	defer discoverer.Shutdown()

	discoverer.ScanEvery(5 * time.Second)

	deadline := time.After(timeout)
	foundDevices := []*eaas.Device{}

	log.Printf("Listening for devices for %s", timeout)

discoveryLoop:
	for {
		select {
		case <-deadline:
			break discoveryLoop
		default:
			device, err := discoverer.Discover(timeout)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
				continue discoveryLoop
			}
			if device == nil {
				continue
			}
			// check if we already found this device before
			for _, foundDevice := range foundDevices {
				if foundDevice.IsEqual(device) {
					continue discoveryLoop
				}
			}
			foundDevices = append(foundDevices, device)
			log.Printf("%s %q %q", device.Hostname, device.URL, device.SoftwareVersion)
		}
	}

	log.Printf("Found devices: %d", len(foundDevices))
}
