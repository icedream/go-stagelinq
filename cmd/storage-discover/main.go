package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
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
	key      ed25519.PrivateKey
	id       uuid.UUID
)

func init() {
	flag.StringVar(&grpcURL, "server", "", "GRPC URL of the remote Engine Library to connect to. If empty, will discover devices instead.")
	flag.Parse()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "eaas-demo"
	}

	if f, err := os.Open("eaas-key.bin"); err == nil {
		defer f.Close()
		keyBytes, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		readKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
		if err != nil {
			panic(err)
		}
		if edkey, ok := readKey.(ed25519.PrivateKey); !ok {
			panic("eaas-key.bin is not an ed25519 private key")
		} else {
			key = edkey
		}
	}
	if key == nil {
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			panic(err)
		}
		os.WriteFile("eaas-key.bin", keyBytes, 0o600)
		key = priv
	}

	if f, err := os.Open("eaas-id.txt"); err == nil {
		defer f.Close()
		keyBytes, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		id, err = uuid.ParseBytes(keyBytes)
		if err != nil {
			panic(err)
		}
	}
	if key == nil {
		id, err = uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		keyBytes, err := id.MarshalBinary()
		if err != nil {
			panic(err)
		}
		os.WriteFile("eaas-id.txt", keyBytes, 0o600)
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

func runEngineLibraryUI(grpcURL string) {
	ctx := context.Background()
	connection, err := eaas.DialContext(ctx, grpcURL)
	if err != nil {
		panic(err)
	}

	// pk := string(key.Public().(ed25519.PublicKey))
	pk := id.String()
	log.Println("Waiting for approval on the other end...")
	resp, err := connection.CreateTrust(ctx, &networktrust.CreateTrustRequest{
		DeviceName: &hostname,
		Ed25519Pk:  &pk,
	})
	if err != nil {
		panic(err)
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
	var pageSize uint32 = 100
	for _, playlist := range getLibraryResp.GetPlaylists() {
		log.Printf("Playlist %q (%q)", playlist.GetTitle(), playlist.GetListType())

		getTracksResp, err := connection.GetTracks(ctx, &enginelibrary.GetTracksRequest{
			PlaylistId: playlist.Id,
			PageSize:   &pageSize,
		})
		if err != nil {
			panic(err)
		}
		for _, track := range getTracksResp.GetTracks() {
			metadata := track.GetMetadata()
			if metadata == nil {
				continue
			}
			log.Printf("\tTrack %s", metadata.String())
		}
	}
}

func runDiscovery() {
	listener, err := eaas.ListenWithConfiguration(&eaas.ListenerConfiguration{
		DiscoveryTimeout: timeout,
	})
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	listener.SendBeaconEvery(5 * time.Second)

	deadline := time.After(timeout)
	foundDevices := []*eaas.Device{}

	log.Printf("Listening for devices for %s", timeout)

discoveryLoop:
	for {
		select {
		case <-deadline:
			break discoveryLoop
		default:
			device, err := listener.Discover(timeout)
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
