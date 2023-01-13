package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// standard domestic root nameserver ip
var rootNameServer = "198.41.0.4"

func resolve(name string) net.IP {
	// Start from the nameserver
	nameserver := net.ParseIP(rootNameServer)
	for {
		reply := dnsQuery(name, nameserver)
		if ip := getAnswer(reply); ip != nil {
			fmt.Println("IP: ", ip)
			return ip
		} else if nsIP := getGlue(reply); nsIP != nil {
			nameserver = nsIP
		} else {
			panic("Cannot resolve the IP")
		}
	}
}

func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println(" ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println(" ", record)
			return record.(*dns.A).A
		}
	}

	return nil
}

func dnsQuery(name string, server net.IP) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server.String(), name)
	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	client := new(dns.Client)
	reply, _, _ := client.Exchange(msg, server.String()+":53")
	return reply
}

func main() {
	name := os.Args[1]
	fmt.Println("Finding the IP for: ", name)
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	fmt.Println("Result: ", resolve(name))
}

// Steps
// Send DNS query to root name server
// Root name server will return the NS address of specific TLD
// Ask the same DNS query to that NS TLD
// It will provide the address of Authorative NS
// Ask the IP(dns query) to that address
// Yay we got our IP
