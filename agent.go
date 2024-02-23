package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type AgentCfg struct {
	AgentUrl   string
	AuthToken  string
	AgentToken string
}

func NewAgentCfg() *AgentCfg {
	url := os.Getenv("TRACKSSL_AGENT_URL")

	if url == "" {
		url = "https://home.andylibby.org:443"
	}

	return &AgentCfg{
		AgentUrl:   url,
		AuthToken:  os.Getenv("TRACKSSL_AUTH_TOKEN"),
		AgentToken: os.Getenv("TRACKSSL_AGENT_TOKEN"),
	}
}

type Agent struct {
	AgentCfg *AgentCfg
}

func (a *Agent) DomainsUrl() string {
	return fmt.Sprintf("%s/api/v1/agents/%s/domains.json", a.AgentCfg.AgentUrl, a.AgentCfg.AgentToken)
}

func (a *Agent) fetchDomains() []Domain {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", a.DomainsUrl(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AgentCfg.AuthToken))

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	response := Response{}
	json.Unmarshal(body, &response)
	return response.Data
}

func (a *Agent) fetchCert(domain Domain) *x509.Certificate {
	conn, _ := net.Dial("tcp", domain.String())

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: domain.Hostname,
	}

	tlsConn := tls.Client(conn, tlsCfg)
	tlsConn.Handshake()
	defer conn.Close()

	return tlsConn.ConnectionState().PeerCertificates[0]
}

func NewAgent() *Agent {
	return &Agent{
		AgentCfg: NewAgentCfg(),
	}
}

