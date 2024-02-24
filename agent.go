package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"time"
	"trackssl.com/agent/client"
)

var SLEEP_DURATION = 5 * time.Minute

type Agent struct {
}

func (a *Agent) fetchCert(domain client.Domain) *x509.Certificate {
	conn, _ := net.Dial("tcp", domain.String())

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         domain.Hostname,
	}

	tlsConn := tls.Client(conn, tlsCfg)
	tlsConn.Handshake()
	defer conn.Close()

	return tlsConn.ConnectionState().PeerCertificates[0]
}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) run() {
	client, err := client.NewClient()

	if err != nil {
		fmt.Println("Error creating client")
		fmt.Println(err)
	}

	fmt.Println("TrackSSL Agent started...")

	for {
		fmt.Println("Retriving domain list...")
		domains := client.FetchDomains()
		fmt.Printf("Retrieved %d domains\n", len(domains))
		fmt.Println("Fetching certificates...")
		for _, domain := range domains {
			c := a.fetchCert(*domain)
			domain.Cert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}))
			client.SendCert(domain)
		}

		fmt.Printf("Sleeping for %s\n", SLEEP_DURATION)
		fmt.Printf("Next run at %s\n", time.Now().Add(SLEEP_DURATION))
		time.Sleep(SLEEP_DURATION)
	}
}

func main() {
	agent := NewAgent()
	agent.run()
}
