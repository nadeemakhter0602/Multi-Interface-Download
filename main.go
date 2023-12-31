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
	"time"
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

func downloadRange(startBytes, endBytes int64, laddr, uri string, file *os.File, ch chan int) {
	contentRange := fmt.Sprintf("bytes=%d-%d", startBytes, endBytes)
	client := createClient(laddr)
	request, err := http.NewRequest("GET", uri, nil)
	PanicErr(err)
	request.Header.Set("Range", contentRange)
	response, err := client.Do(request)
	PanicErr(err)
	responseReader := response.Body
	defer responseReader.Close()
	// create a 64 byte buffer
	buffer := make([]byte, 64)
	offset := startBytes
	// write to file in 64 byte chunks
	for {
		bytesRead, readErr := responseReader.Read(buffer)
		bytesWritten, writeErr := file.WriteAt(buffer[:bytesRead], offset)
		PanicErr(writeErr)
		ch <- bytesWritten
		offset += int64(bytesWritten)
		if readErr == io.EOF {
			break
		}
		PanicErr(readErr)
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
	interval := math.Ceil(float64(contentLength) / float64(len(laddrs)))
	offset := 0.0
	start := time.Now()
	ch := make(chan int)
	fmt.Printf("Starting download of file %s of size %d bytes\n", fName, contentLength)
	for _, laddr := range laddrs {
		startBytes := int64(offset)
		offset += interval
		endBytes := int64(math.Min(offset, float64(contentLength-1)))
		offset += 1
		// assign a goroutine for each interface
		go downloadRange(startBytes, endBytes, laddr, uri, file, ch)
	}
	totalBytesWritten := 0.0
	fSize := float64(contentLength)
	for totalBytesWritten < fSize {
		bytesWritten := <-ch
		totalBytesWritten += float64(bytesWritten)
		scale := (totalBytesWritten / fSize) * 100
		fmt.Printf("[%.2f/100.00]\r", scale)
	}
	fmt.Printf("\nDownload complete in %s\n", time.Since(start))
}
