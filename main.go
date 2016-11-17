package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	Name     string
	IPAdress string `json:"ip_address"`
	Country  string
	Load     int
}

func getServers(url string) (servers []Server) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &servers)

	return
}

func main() {
	url := "https://api.nordvpn.com/server"
	servers := getServers(url)

	for _, server := range servers {
		fmt.Println(server.IPAdress)
	}
}
