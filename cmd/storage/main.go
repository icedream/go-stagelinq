package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icedream/go-stagelinq/eaas"
	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
	"google.golang.org/grpc"
)

const (
	appName    = "Icedream StagelinQ Storage"
	appVersion = "0.0.0"
	timeout    = 5 * time.Second
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "eaas-demoserver"
	}
}

func main() {
	// Generate random token to identify with.
	//
	// Engine uses the token to know whether you just logged onto the network or
	// whether you're a library that just restarted. For our demo purposes this
	// doesn't matter too much though, so we just regenerate on bootup.
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		panic(err)
	}

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, stopNotify := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)
	defer stopNotify()

	var s http.Server
	grpcServer := grpc.NewServer()
	go func() {
		<-ctx.Done()

		grpcServer.Stop()

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		s.Shutdown(timeout)
	}()

	// Set up gRPC API
	grpcPort := eaas.DefaultEAASGRPCPort
	grpcListener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: int(grpcPort),
	})
	if err != nil {
		panic(err)
	}
	enginelibrary.RegisterEngineLibraryServiceServer(grpcServer, &EngineLibraryServiceServer{})
	networktrust.RegisterNetworkTrustServiceServer(grpcServer, &NetworkTrustServiceServer{})
	go func() {
		log.Println("Listening on GRPC")
		_ = grpcServer.Serve(grpcListener)
	}()

	// Set up HTTP server
	s.Addr = fmt.Sprintf(":%d", eaas.DefaultEAASHTTPPort)
	s.Handler = eaasHTTPHandler()
	go func() {
		log.Println("Listening on HTTP")
		_ = s.ListenAndServe()
	}()

	// Listen for beacon UDP broadcasts
	log.Println("Beacon starting")
	beacon, err := eaas.StartBeaconWithConfiguration(&eaas.BeaconConfiguration{
		Name:            hostname,
		SoftwareVersion: appVersion,
		GRPCPort:        grpcPort,
		Token:           demoToken,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println("Beacon shutting down")
		beacon.Shutdown()
	}()

	// Wait for interrupt/term
	log.Println("Running")
	<-ctx.Done()
}
