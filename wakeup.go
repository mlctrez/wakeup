package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/linde12/gowol"
)

func main() {
	var hostsUrl string
	flag.StringVar(&hostsUrl, "hostsUrl", "https://optiplex/hosts", "the url to retrieve a list of hosts from")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("provide one or more short host names as arguments")
		os.Exit(1)
	}

	resp, err := http.Get(hostsUrl)
	if err != nil {
		log.Fatal(err)
	}
	api := struct {
		Hosts []struct {
			Name       string `json:"Id"`
			IpAddress  string `json:"Ip"`
			MacAddress string `json:"Mac"`
		} `json:"hosts"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&api)
	if err != nil {
		log.Fatal(err)
	}

	for _, h := range api.Hosts {
		for _, a := range args {
			if strings.ToLower(h.Name) == strings.ToLower(a) {
				packet, err := gowol.NewMagicPacket(h.MacAddress)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("wakeup", h.Name, "at", h.MacAddress)
				err = packet.Send("255.255.255.255")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
