package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"time"
	//"net"
	"log"
)

var (
	snapshot_len int32 = 1024
	err          error
	timeout      = 30 * time.Second
	handle       *pcap.Handle
)

func main() {
	devices := listDevices()
	for num, i := range devices {
		fmt.Printf("Device %d: %s \n", num, i.Name)
	}
	ifaceParser := argparse.NewParser("iface", "Pass the interface using to scan")
	iface := ifaceParser.String("i", "interface", &argparse.Options{Required: true,
		Help: "Put network interface"})
	err = ifaceParser.Parse(os.Args)
	if err == nil {
		capturePackets(iface)
	} else {
		log.Fatal(err)
	}
}

func listDevices() []pcap.Interface {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	return devices
}

func capturePackets(deviceChoice *string) {
	handle, err = pcap.OpenLive(*deviceChoice, snapshot_len, false, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filter := "port 80"
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
