package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ApplyLogic/mqtt-broker/config"
	mqttServer "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log"
	"os"
)

type Broker struct {
	Server *mqttServer.Server
}

func New(cfg *config.Config) *Broker {

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

	server := mqttServer.New(nil)
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

	return &Broker{
		Server: server,
	}
}
