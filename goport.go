package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	snapshot_len int32 = 1024
	err          error
	timeout      = 30 * time.Second
	handle       *pcap.Handle
)

func main() {
	ifaceParser := argparse.NewParser("iface", "Capture the packets")
	iface := ifaceParser.String("i", "interface", &argparse.Options{Help: "Network interface"})
	list := ifaceParser.Flag("l", "list", &argparse.Options{Help: "List all devices"})
	ble := ifaceParser.Flag("b", "ble", &argparse.Options{Help: "Discover Bluetooth devices"})
	err = ifaceParser.Parse(os.Args)
	if err != nil {
		fmt.Print(ifaceParser.Usage(err))
	}

	if *iface != "" && checkDevice(iface) {
		capturePackets(iface)
	} else if *iface == "" && *list {
		for num, i := range listDevices() {
			fmt.Printf("Device %d: %s \n", num, i.Name)
		}
	} else if *ble {
		runBLE()
	}
}

func listDevices() []net.Interface {
	devices, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	return devices
}

func checkDevice(input *string) bool {
	allDevices := listDevices()
	for _, i := range allDevices {
		if strings.Compare(*input, i.Name) == 1 {
			return true
		}
	}
	return false
}

func capturePackets(deviceChoice *string) {
	handle, err = pcap.OpenLive(*deviceChoice, snapshot_len, false, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filter := "port 8443"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] Start capturing on %s...", *deviceChoice)
	packetsSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetsSource.Packets() {
		fmt.Println(packet)
	}
}
