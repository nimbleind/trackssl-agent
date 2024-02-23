package main

import (
	"fmt"
)

func main() {
	agent := NewAgent()
	domains := agent.fetchDomains()

	for _, domain := range domains {
		c := agent.fetchCert(domain)

		fmt.Println("================================================")
		fmt.Println(domain.String())
		fmt.Println(c.NotBefore)
		fmt.Println(c.NotAfter)
		fmt.Println(c.Issuer)
		fmt.Println(c.Subject)
		fmt.Println(c.DNSNames)
	}
}
