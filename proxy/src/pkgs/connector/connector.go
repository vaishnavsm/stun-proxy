package connector

import (
	"log"
	"os"
)

func validate(addr *string, broker *string) {
	if addr == nil {
		log.Fatal("connector address is nil")
	}
	if broker == nil {
		log.Fatal("broker address is nil")
	}
}

func Run(addr *string, broker *string, interrupt chan os.Signal) {
	validate(addr, broker)
}
