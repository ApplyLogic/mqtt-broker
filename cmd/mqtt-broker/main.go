package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ApplyLogic/mqtt-broker/config"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func main() {

	cfg, err := config.LoanConfig()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Load tls cert from your cert file
	cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/test_cert.pem", cfg.TLS.CertPath), fmt.Sprintf("%s/test_cert.key", cfg.TLS.CertPath))
	if err != nil {
		log.Fatal(err)
	}

	// Basic TLS Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Optionally, if you want clients to authenticate only with certs issued by your CA,
	// you might want to use something like this:
	if cfg.TLS.CACertFile != "" {
		pemCACert, err := os.ReadFile(cfg.TLS.CACertFile)
		if err != nil {
			log.Fatal(err)
		}
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(pemCACert)
		if ok {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = certPool
		}
	}

	server := mqtt.New(nil)
	_ = server.AddHook(new(auth.AllowHook), nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      cfg.MQTT.TCP.Id,
		Address: fmt.Sprintf(":%s", cfg.MQTT.TCP.Address),
		//TLSConfig: tlsConfig,
	})
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	ws := listeners.NewWebsocket(listeners.Config{
		ID:      cfg.MQTT.WebSocket.Id,
		Address: fmt.Sprintf(":%s", cfg.MQTT.WebSocket.Address),
		//TLSConfig: tlsConfig,
	})
	err = server.AddListener(ws)
	if err != nil {
		log.Fatal(err)
	}

	stats := listeners.NewHTTPStats(
		listeners.Config{
			ID:      cfg.MQTT.HTTPStats.Id,
			Address: fmt.Sprintf(":%s", cfg.MQTT.HTTPStats.Address),
			//TLSConfig: tlsConfig,
		}, server.Info,
	)
	err = server.AddListener(stats)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
		server.Log.Info("Broker started")
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	server.Log.Info("main.go finished")
}
