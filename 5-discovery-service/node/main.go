package main

import (
	"discovery-service/node/consul"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

func resolveHostIp() string {
	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ip := networkIp.IP.String()
			return ip
		}
	}
	return ""
}

func findNode(client consul.Client, serviceName string) {
	serviceEntry, _, err := client.Service(serviceName, "")

	if err != nil {
		fmt.Errorf("Error get consul client: ", err.Error())
	}

	go func() {
		for _, entry := range serviceEntry {
			checkHealth("http://" + entry.Service.Address + ":" + strconv.Itoa(entry.Service.Port) + "/health")
		}
	}()
}

func checkHealth(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get health:\n", string(body))
}

func main() {
	addressConsul := flag.String("consul", "localhost:8500", "address discovery service. 'Format host:port'")
	nodeID := flag.String("nodeID", "default-ID", "node ID")
	serviceName := flag.String("serviceName", "default-service", "service name")
	listenPort := flag.Int("port", 11111, "listen service port")
	flag.Parse()

	fmt.Println("Hello, I am ", *nodeID)

	client, err := consul.NewConsulClient(*addressConsul)
	if err != nil {
		fmt.Errorf("Error create consul client: ", err.Error())
	}

	ipAddress := resolveHostIp()

	client.Register(*serviceName, ipAddress, "http", *listenPort)

	findNode(client, *serviceName)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/health from " + *nodeID)
		fmt.Fprintf(w, "Hello, %q. I am %q", html.EscapeString(r.URL.Path), *nodeID)
	})

	port := ":" + strconv.Itoa(*listenPort)
	http.ListenAndServe(port, nil)

	client.DeRegister(*serviceName)
}
