package main

import (
	"context"
	"flag"
	"math"
	"net"
	"net/http"
	"strings"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getFileDetails(uri string) (int64, string) {
	res, err := http.Head(uri)
	PanicErr(err)
	contentLength := res.ContentLength
	headers := res.Header
	_, ok := headers["Accept-Ranges"]
	if !ok {
		panic("Server does not support HTTP Range Requests")
	}
	slashSplit := strings.Split(uri, "/")
	lastStringSlash := slashSplit[len(slashSplit)-1]
	fName := strings.Split(lastStringSlash, "?")[0]
	return contentLength, fName
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

func downloadRange(startBytes, endBytes int64, laddr, uri, fName string) {
	panic("Not Implemented")
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
	contentLength, fName := getFileDetails(uri)
	interval := math.Floor(float64(contentLength) / float64(len(laddrs)))
	offset := 0.0
	for _, laddr := range laddrs {
		startBytes := int64(offset)
		offset += interval
		endBytes := int64(offset)
		offset += 1
		go downloadRange(startBytes, endBytes, laddr, uri, fName)
	}
}
