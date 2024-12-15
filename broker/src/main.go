package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/vaishnavsm/stun-proxy/broker/src/pkgs/broker_server"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"details":      "connect websocket to /ws to initiate STUN connection",
		"what is this": "this is a STUN based",
		"github":       "https://github.com/vaishnavsm/stun-proxy",
		"author":       "vaishnavsm",
		"more":         "vaishnavsm.com",
	})
}

func main() {
	flag.Parse()

	if addr == nil {
		log.Fatal("Listen Address is nil")
	}

	b, err := broker_server.New()
	if err != nil {
		log.Fatal("Error starting broker", err)
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", b.Serve)

	log.Printf("starting broker on %s", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
