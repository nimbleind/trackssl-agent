package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"trackssl.com/agent/client"
)

var (
	DEFAULT_TRACKSSL_URL = "https://app.trackssl.com"
	SLEEP_DURATION       = 5 * time.Minute
	ERROR_NO_AUTH_TOKEN  = fmt.Errorf("No auth token")
	ERROR_NO_AGENT_TOKEN = fmt.Errorf("No agent token")
)

type Agent struct {
	TracksslUrl *string
	AuthToken   *string
	AgentToken  *string
}

func (a *Agent) fetchCert(domain client.Domain) (*x509.Certificate, error) {
	conn, err := net.Dial("tcp", domain.String())
	if err != nil {
		return nil, fmt.Errorf("Error connecting %s: %s", domain.String(), err)
	}

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         domain.Hostname,
	}

	tlsConn := tls.Client(conn, tlsCfg)
	tlsConn.Handshake()
	defer conn.Close()

	return tlsConn.ConnectionState().PeerCertificates[0], nil
}

func (a *Agent) checkConfig() error {
	if a.AuthToken == nil {
		return ERROR_NO_AUTH_TOKEN
	}

	if a.AgentToken == nil {
		return ERROR_NO_AGENT_TOKEN
	}

	return nil
}

func NewAgent() (*Agent, error) {
	agent := &Agent{}

	agent.TracksslUrl = flag.String("trackssl-url", os.Getenv("TRACKSSL_URL"), "TrackSSL URL")
	agent.AuthToken = flag.String("auth-token", os.Getenv("TRACKSSL_AUTH_TOKEN"), "TrackSSL Auth Token")
	agent.AgentToken = flag.String("agent-token", os.Getenv("TRACKSSL_AGENT_TOKEN"), "TrackSSL Agent Token")
	flag.Parse()

	if *agent.TracksslUrl == "" {
		agent.TracksslUrl = &DEFAULT_TRACKSSL_URL
	}

	return agent, agent.checkConfig()
}

func (a *Agent) NewClient() *client.Client {
	return &client.Client{
		TracksslUrl: *a.TracksslUrl,
		AuthToken:   *a.AuthToken,
		AgentToken:  *a.AgentToken,
	}
}

func (a *Agent) run() {
	client := a.NewClient()

	fmt.Printf("Initializing TrackSSL agent with interval %v\n", SLEEP_DURATION)

	for {
		fmt.Println("Retriving domain list...")
		domains, err := client.FetchDomains()

		if err != nil {
			fmt.Println(err)
		}

		if len(domains) > 0 {
			fmt.Printf("Retrieved %d domains\n", len(domains))
			fmt.Println("Fetching certificates...")
			for _, domain := range domains {
				c, err := a.fetchCert(*domain)

				if err != nil {
					fmt.Println(err)
					domain.Error = fmt.Sprintf("%s", err)
					domain.Cert = ""
				} else {
					domain.Cert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}))
				}

				err = client.SendCert(domain)
				if err != nil {
					fmt.Printf("Failed to send certificate for %s: %s\n", domain.Hostname, err)
				}

			}
		}

		fmt.Printf("Sleeping for %s\n", SLEEP_DURATION)
		fmt.Printf("Next run at %s\n", time.Now().Add(SLEEP_DURATION))
		time.Sleep(SLEEP_DURATION)
	}
}

func main() {
	agent, err := NewAgent()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	agent.run()
}
