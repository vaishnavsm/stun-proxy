package config

import "os"

type Config struct {
	Name       string
	Addr       string
	BrokerAddr string
	Interrupt  chan os.Signal
}
