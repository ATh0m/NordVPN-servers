package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

type Server struct {
	Name      string
	IPAddress string `json:"ip_address"`
	Load      int
	RTT       time.Duration
}

type Servers []Server

func (s Servers) Len() int {
	return len(s)
}

func (s Servers) Less(i, j int) bool {
	return s[i].RTT < s[j].RTT
}

func (s Servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getServers(url string) (servers Servers) {
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

func getPing(servers *Servers) {
	ping := fastping.NewPinger()
	ping.Network("udp")

	ping.MaxRTT = time.Second

	table := map[string]*Server{}

	for i, server := range *servers {
		table[server.IPAddress] = &(*servers)[i]
		ping.AddIP(server.IPAddress)
	}

	ping.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		(*table[addr.String()]).RTT = t
	}

	ping.Run()
}

func removeUnreachableServers(servers Servers) (result Servers) {
	for _, server := range servers {
		if server.RTT != 0 {
			result = append(result, server)
		}
	}
	return
}

func main() {
	url := "https://api.nordvpn.com/server"
	servers := getServers(url)

	getPing(&servers)
	sort.Sort(servers)

	servers = removeUnreachableServers(servers)

	for _, server := range servers {
		fmt.Printf("%-25v %-20v %v %%\n", server.Name, server.RTT, server.Load)
	}
}
