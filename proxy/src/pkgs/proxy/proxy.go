package proxy

import (
	"encoding/json"
	"log"
	"os"

	brokerconn "github.com/vaishnavsm/stun-proxy/proxy/src/pkgs/broker_conn"
)

const CONFIG_SCHEMA_MESSAGE = `
You've left the proxy config empty. Please set it with the "-proxyConfig" flag.
The value of this flag is an application, a JSON object following this Schema:
{
	"name": "name of the application",
	"addr": "address of the application to connect to"
}`

func printConfigSchema() {
	log.Println(CONFIG_SCHEMA_MESSAGE)
	os.Exit(0)
}

func validate(config *string, broker *string) {
	if config == nil {
		log.Fatal("proxy config is nil")
	}
	if broker == nil {
		log.Fatal("broker address is nil")
	}
	if *config == "{}" {
		printConfigSchema()
	}
}

type Config struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

func Run(configStr *string, broker *string, interrupt chan os.Signal) {
	validate(configStr, broker)

	c := &Config{}
	err := json.Unmarshal([]byte(*configStr), c)
	if err != nil {
		log.Fatal("error parsing json in config parameter", err)
	}

	b, err := brokerconn.New(*broker, interrupt)
	if err != nil {
		log.Fatal("could not connect to broker", err)
	}

	appId, err := b.RegisterApplication(c.Name)
	if err != nil {
		log.Fatalf("Could not register application %s: %v", c.Name, err)
	}

	app, err := New(c.Name, interrupt, b, c.Addr, appId)
	if err != nil {
		log.Fatalf("Could not create application %s: %v", c.Name, err)
	}

	app.Serve()
}
