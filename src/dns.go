package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/miekg/dns"
)

func resolve(name string, rootNameServer string) net.IP {
	// Start from the nameserver
	nameserver := net.ParseIP(rootNameServer)
	for {
		reply, err := dnsQuery(name, nameserver)
		if err != nil {
			log.Fatal("DNS query timeout: ", err)
		}
		if ip := getAnswer(reply); ip != nil {
			return ip
		} else if nsIP := getGlue(reply); nsIP != nil {
			nameserver = nsIP
		} else if domain := getNS(reply); domain != "" {
			nameserver = resolve(domain, rootNameServer)
		} else {
			break
		}
	}
	return nil
}

// This function parses the "answer section" of DNS response, where the IP address of your website is present
func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Ans: ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

// This is additional section, where “glue records” live.
// Glue records holds IP address of nameserver where your query is routed
func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("Glue: ", record)
			return record.(*dns.A).A
		}
	}

	return nil
}

// This is the "authority section", it has domain names of the other nameservers where your query is routed
func getNS(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("NS: ", record)
			return record.(*dns.NS).Ns
		}
	}
	return ""
}

// This function takes care of preparing the DNS query and sending them over UDP
func dnsQuery(name string, server net.IP) (*dns.Msg, error) {
	fmt.Printf("dig @%s %s\n", server.String(), name)
	// prepare the dns query
	msg := new(dns.Msg)
	// set the domain name we are querying in question
	msg.SetQuestion(name, dns.TypeA)
	// initalizing DNS client
	client := new(dns.Client)
	// send the request over UDP
	reply, _, err := client.Exchange(msg, server.String()+":53")
	return reply, err
}

func main() {
	// standard domestic root nameserver(.) ip
	var rootNameServer = "198.41.0.4"
	// extracting domain name for which we will find IP address
	name := os.Args[1]
	fmt.Println("Finding the IP for: ", name)
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	// run the DNS resolver and print the result
	ip := resolve(name, rootNameServer)
	fmt.Println()
	if ip == nil {
		red := color.New(color.FgHiRed)
		boldRed := red.Add(color.Bold)
		boldRed.Printf("Unable to find the ip for %s\n", name)
	} else {
		cyan := color.New(color.FgCyan)
		cyan.Printf("IP: %s\n", ip)
	}
}
