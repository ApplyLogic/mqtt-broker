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
	// Create MQTT server.
	log.Println("Create MQTT broker")
	opts := &mqttServer.Options{}
	server := mqttServer.New(opts)

	//tlsConfig, err := configureTLS()
	//if err != nil {
	//	log.Fatalf("Error configuring TLS: %v", err)
	//}

	// Add an authentication hook (example: allow all).  For production, replace with proper authentication.
	err := server.AddHook(new(auth.AllowHook), nil)
	if err != nil {
		log.Fatalf("Error adding authentication hook: %v", err)
	}

	// Create TCP listener with TLS.
	tcpListener := listeners.NewTCP(listeners.Config{
		Type:    "",
		ID:      cfg.MQTT.TCP.Id,
		Address: fmt.Sprintf(":%s", cfg.MQTT.TCP.Address),
		//TLSConfig: tlsConfig,
	})
	err = server.AddListener(tcpListener)
	if err != nil {
		log.Fatalf("Error adding listener: %v", err)
	}

	// Create Websocket listener with TLS.
	wsListener := listeners.NewWebsocket(listeners.Config{
		Type:    "",
		ID:      cfg.MQTT.WebSocket.Id,
		Address: fmt.Sprintf(":%s", cfg.MQTT.WebSocket.Address),
		//TLSConfig: tlsConfig,
	})
	err = server.AddListener(wsListener)
	if err != nil {
		log.Fatal(err)
	}

	// Create HTTP Stats listener with TLS.
	statsListener := listeners.NewTCP(listeners.Config{
		Type:    "",
		ID:      cfg.MQTT.HTTPStats.Id,
		Address: fmt.Sprintf(":%s", cfg.MQTT.HTTPStats.Address),
		//TLSConfig: tlsConfig,
	})
	err = server.AddListener(statsListener)
	if err != nil {
		log.Fatal(err)
	}

	return &Broker{
		Server: server,
	}
}

func configureTLS() (*tls.Config, error) {
	caCert, err := os.ReadFile("/Users/golanshay/Workspace/tls-certs/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair("/Users/golanshay/Workspace/tls-certs/server-cert.pem", "/Users/golanshay/Workspace/tls-certs/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
