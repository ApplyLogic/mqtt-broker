package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MQTT *MQTT
	TLS  *TLS
}

type MQTT struct {
	TCP       TCP
	WebSocket WebSocket
	HTTPStats HTTPStats
}

type TCP struct {
	Id      string
	Address string
}

type WebSocket struct {
	Id      string
	Address string
}

type HTTPStats struct {
	Id      string
	Address string
}

type TLS struct {
	CACertFile     string
	ServerCertFile string
	ServerKeyFile  string
}

func LoanConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		MQTT: &MQTT{
			TCP: TCP{
				Id:      viper.GetString("TCP_ID"),
				Address: viper.GetString("TCP_ADDRESS"),
			},
			WebSocket: WebSocket{
				Id:      viper.GetString("WS_ID"),
				Address: viper.GetString("WS_ADDRESS"),
			},
			HTTPStats: HTTPStats{
				Id:      viper.GetString("HTTP_STATS_ID"),
				Address: viper.GetString("HTTP_STATS_ADDRESS"),
			},
		},
		TLS: &TLS{
			ServerCertFile: viper.GetString("TLS_SERVER_CERT_FILE"),
			ServerKeyFile:  viper.GetString("TLS_SERVER_KEY_FILE"),
			//if you want clients to authenticate only with certs issued by your CA
			CACertFile: viper.GetString("TLS_CA_CERT_FILE"),
		},
	}
	return cfg, nil
}
