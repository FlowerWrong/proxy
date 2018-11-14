// https://github.com/songgao/water#tun-on-macos
package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/FlowerWrong/kone/tcpip"
	"github.com/songgao/water"
)

// sudo go run utun.go
// sudo ifconfig utun2 10.1.0.10 10.1.0.20 up
// ping 10.1.0.20
func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	b := make([]byte, 1500)
	for {
		n, err := ifce.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		packet := b[:n]
		if tcpip.IsIPv4(packet) {
			ipPacket := tcpip.IPv4Packet(packet)
			icmpPacket := tcpip.ICMPPacket(ipPacket.Payload())
			if icmpPacket.Type() == tcpip.ICMPRequest && icmpPacket.Code() == 0 {
				log.Printf("icmp echo request: %s -> %s\n", ipPacket.SourceIP(), ipPacket.DestinationIP())
				// forge a reply
				icmpPacket.SetType(tcpip.ICMPEcho)
				srcIP := ipPacket.SourceIP()
				dstIP := ipPacket.DestinationIP()
				ipPacket.SetSourceIP(dstIP)
				ipPacket.SetDestinationIP(srcIP)

				icmpPacket.ResetChecksum()
				ipPacket.ResetChecksum()
				ifce.Write(ipPacket)
			} else {
				log.Printf("icmp: %s -> %s\n", ipPacket.SourceIP(), ipPacket.DestinationIP())
			}
		}
	}
}
