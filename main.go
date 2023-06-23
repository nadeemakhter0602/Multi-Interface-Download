package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

// create a WaitGroup
var wg sync.WaitGroup

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

func createClient(laddr string) *http.Client {
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

func downloadRange(startBytes, endBytes int64, laddr, uri string, file *os.File) {
	// Schedule call to WaitGroup's Done to tell goroutine is completed
	defer wg.Done()
	contentRange := fmt.Sprintf("bytes %d-%d", startBytes, endBytes)
	client := createClient(laddr)
	request, err := http.NewRequest("GET", uri, nil)
	PanicErr(err)
	request.Header.Set("Content-Range", contentRange)
	response, err := client.Do(request)
	PanicErr(err)
	responseReader := response.Body
	defer responseReader.Close()
	buffer := make([]byte, 64)
	file.Seek(startBytes, io.SeekStart)
	for {
		bytesRead, err := responseReader.Read(buffer)
		file.Write(buffer[:bytesRead])
		if err == io.EOF {
			break
		}
	}
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
	file, err := os.Create(fName)
	PanicErr(err)
	defer file.Close()
	interval := math.Floor(float64(contentLength) / float64(len(laddrs)))
	offset := 0.0
	// add count of all goroutines
	wg.Add(len(laddrs))
	for _, laddr := range laddrs {
		startBytes := int64(offset)
		offset += interval
		endBytes := int64(offset)
		offset += 1
		// assign a goroutine for each interface
		go downloadRange(startBytes, endBytes, laddr, uri, file)
	}
	// wait for the goroutines to finish execution
	wg.Wait()
}
