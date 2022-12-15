package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dcasado/wol-api/magicpacket"
)

const (
	listenAddressEnvVariable = "LISTEN_ADDRESS"
	listenPortEnvVariable    = "LISTEN_PORT"

	wolEndpoint    = "/wol"
	healthEndpoint = "/health"

	// Only 48-bit MACs are allowed
	macLenth = 6
)

type WOLBody struct {
	MAC string `json:"mac"`
}

func main() {
	listenAddress := getListenAddressEnvVariable()
	listenPort := getListenPortEnvVariable()

	log.Printf("Starting server listening on %s:%s", listenAddress, listenPort)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(wolEndpoint, handleWOL)
	serveMux.HandleFunc(healthEndpoint, handleHealth)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", listenAddress, listenPort),
		Handler: serveMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting the server: %s", err)
	}
}

func getListenAddressEnvVariable() string {
	value := os.Getenv(listenAddressEnvVariable)
	if len(value) != 0 {
		return value
	}
	return "127.0.0.1"
}

func getListenPortEnvVariable() string {
	value := os.Getenv(listenPortEnvVariable)
	if len(value) != 0 {
		return value
	}
	return "9099"
}

func handleWOL(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var b WOLBody

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&b)
	if err != nil {
		http.Error(responseWriter, "Malformed request", http.StatusBadRequest)
		return
	}

	mac, err := net.ParseMAC(b.MAC)
	if err != nil {
		http.Error(responseWriter, "Malformed MAC address", http.StatusBadRequest)
		return
	}
	if len(mac) > macLenth {
		http.Error(responseWriter, "Only 48-bit MAC address are allowed", http.StatusBadRequest)
		return
	}

	mp := magicpacket.New(mac)

	ba, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		http.Error(responseWriter, "Error resolving broadcast address", http.StatusInternalServerError)
		return
	}

	c, err := net.DialUDP("udp", nil, ba)
	if err != nil {
		http.Error(responseWriter, "Error opening UDP connection", http.StatusInternalServerError)
		return
	}
	defer c.Close()

	_, err = c.Write(mp)
	if err != nil {
		http.Error(responseWriter, "Error writing to the UDP connection", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/text")
	responseWriter.Write([]byte("Ok"))

}

func handleHealth(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/text")
	responseWriter.Write([]byte("Ok"))
}
