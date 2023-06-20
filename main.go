package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	laddr := flag.String("i", "", "local ip address(es) for the interface(s)")
	url := flag.String("u", "", "URL of the file")
	flag.Parse()
	if *laddr == "" || *url == "" {
		flag.Usage()
		return
	}
	laddrs := strings.Split(*laddr, ",")
	uri := *url
	fmt.Println("Local IP Addresses :", laddrs)
	fmt.Println("URL of the file :", uri)
}
