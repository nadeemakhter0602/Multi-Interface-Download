package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createClient(laddr, uri string) *http.Client {
	addr, err := net.ResolveTCPAddr("tcp", laddr+":0")
	PanicErr(err)
	// create dialer
	dialer := &net.Dialer{LocalAddr: addr}
	// create dial context
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.Dial(network, addr)
		return conn, err
	}
	// create HTTP transport
	transport := &http.Transport{DialContext: dialContext}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func main() {
	// define flags
	laddr := flag.String("i", "", "local ip address(es) for the interface(s)")
	url := flag.String("u", "", "URL of the file")
	// parse flags
	flag.Parse()
	// check if flag values are empty
	if *laddr == "" || *url == "" {
		flag.Usage()
		return
	}
	// split multiple comma seperated local ip addresses
	laddrs := strings.Split(*laddr, ",")
	uri := *url
	fmt.Println("Local IP Addresses :", laddrs)
	fmt.Println("URL of the file :", uri)
}
