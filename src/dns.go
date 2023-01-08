package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var rootNameServer = "198.41.0.4"
	fmt.Println("Going to build our toy dns resolver ✨", rootNameServer)

	for j := 1; j <= 10; j++ {
		fmt.Println("A random number = ", rand.Intn(101))
	}
}

// Steps
// Send DNS query to root name server
// Root name server will return the NS address of specific TLD
// Ask the same DNS query to that NS TLD
// It will provide the address of Authorative NS
// Ask the IP(dns query) to that address
// Yay we got our IP
