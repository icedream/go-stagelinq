package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/icedream/go-stagelinq/eaas/proto/networktrust"
	"golang.org/x/text/encoding/unicode"
	"google.golang.org/grpc"
)

const (
	appName    = "Icedream StagelinQ Storage"
	appVersion = "0.0.0"
	timeout    = 5 * time.Second
)

var (
	eaasMagic         = []byte{'E', 'A', 'A', 'S', 0x01, 0x00}
	eaasResponseMagic = []byte{'E', 'A', 'A', 'S', 0x01, 0x01}
)

var networkStringEncoding = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

func writeNetworkString(w io.Writer, v string) (err error) {
	converted, err := networkStringEncoding.NewEncoder().Bytes([]byte(v))
	if err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(converted))); err != nil {
		return
	}
	_, err = w.Write(converted)
	return
}

func main() {
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		panic(err)
	}

	// listener, err := stagelinq.ListenWithConfiguration(&stagelinq.ListenerConfiguration{
	// 	DiscoveryTimeout: timeout,
	// 	SoftwareName:     appName,
	// 	SoftwareVersion:  appVersion,
	// 	Name:             "testing",
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer listener.Close()

	// listener.AnnounceEvery(time.Second)

	ctx := context.TODO()

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

	// set up grpc listener
	grpcListener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 50010,
	})
	if err != nil {
		panic(err)
	}
	enginelibrary.RegisterEngineLibraryServiceServer(grpcServer, &EngineLibraryServiceServer{})
	networktrust.RegisterNetworkTrustServiceServer(grpcServer, &NetworkTrustServer{})
	go func() {
		log.Println("Listening on GRPC")
		_ = grpcServer.Serve(grpcListener)
	}()

	// set up http listener
	s.Addr = ":50020"
	s.Handler = newHTTPServiceHandler()
	go func() {
		log.Println("Listening on HTTP")
		_ = s.ListenAndServe()
	}()

	// listen for broadcasts
	udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 11224,
	})
	if err != nil {
		panic(err)
	}
	udpC := make(chan *net.UDPAddr, 2)
	go func() {
		b := make([]byte, 6)
		for {
			n, addr, err := udpListener.ReadFromUDP(b)
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
			if err != nil {
				log.Println("UDP error, ignoring:", err)
				continue
			}
			if n != 6 {
				log.Println("UDP message too short, ignoring")
				continue
			}
			if !bytes.Equal(b, eaasMagic) {
				log.Println("UDP broadcast invalid, ignoring")
				continue
			}
			udpC <- addr
		}
	}()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "demo"
	}
	go func() {
		log.Println("Listening on UDP")
		for {
			select {
			case addr := <-udpC:
				msg := new(bytes.Buffer)
				msg.Write(eaasResponseMagic)
				msg.Write(token[:])
				writeNetworkString(msg, hostname)
				uri := fmt.Sprintf("grpc://%s:%d", "192.168.188.120", 50010)
				binary.Write(msg, binary.BigEndian, uint32(len(uri)))
				writeNetworkString(msg, appVersion)
				msg.Write([]byte{0, 0, 0, 2, 0, 0x5f}) // TODO
				b := msg.Bytes()
				log.Println("Sending UDP beacon\n", hex.Dump(b))
				udpListener.WriteToUDP(b, addr)
			case <-ctx.Done():
				_ = udpListener.Close()
				return
			}
		}
	}()

	// wait for interrupt/term
	<-ctx.Done()
}
