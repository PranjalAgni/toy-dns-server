package main

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

func main() {
	// standard domestic root nameserver ip
	var rootNameServer = "198.41.0.4"
	fmt.Println("Going to build our toy dns resolver ✨", rootNameServer)

}

func dnsQuery(name string, server net.IP) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server.String(), name)
	msg := new(dns.Msg)
	msg.SetQuestion("name", dns.TypeA)
	client := new(dns.Client)
	reply, _, _ := client.Exchange(msg, server.String()+":53")
	return reply

}

// Steps
// Send DNS query to root name server
// Root name server will return the NS address of specific TLD
// Ask the same DNS query to that NS TLD
// It will provide the address of Authorative NS
// Ask the IP(dns query) to that address
// Yay we got our IP
