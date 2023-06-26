# Multi-Interface-Download

A simple program for distributed downloading of files over http across multiple interfaces.

It takes multiple comma-seperated local IP addresses of interfaces, and the HTTP URL of the file to download. 

It then uses each interface to get a seperate range of the file using [HTTP Range Requests](https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests), and appropriately writes the data to the file in 64 byte chunks.

The Local IP address(es) of the interface(s) used to download the file can be retrieved by running `ipconfig`. We can use multiple different interfaces, e.g. an Ethernet connection and a seperate WiFi connection, etc.

## Example

Downloading Arch Linux using 2 different interfaces (IP addresses of the interfaces are removed):

```
~$ go run main.go -i XXX.XXX.XXX.XXX,XXX.XXX.XXX.XXX -u https://geo.mirror.pkgbuild.com/iso/2023.06.01/archlinux-x86_64.iso
Starting download of file archlinux-x86_64.iso of size 828715008 bytes
[100.00/100.00]
Download complete in 4m21.6296483s
```