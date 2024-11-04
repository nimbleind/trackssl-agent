package client

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

var (
	tracksslUrl = "https://home.andylibby.org"
	authToken   = "10644d067a80b7ad97b151d6e9fb954ab4b264ebd0860a8912596ce907aa97d4a13221125b634dfbdc7401e1d31f1a13b7aef2e415f0c61b9377190dd87d813e"
	agentToken  = "tucKuJPTYhTSEbEYi2uCk5iT"
	c           = Client{
		TracksslUrl: tracksslUrl,
		AuthToken:   authToken,
		AgentToken:  agentToken,
	}
)

func TestClient_DomainsUrl(t *testing.T) {
	expected := fmt.Sprintf("%s/api/v1/agents/%s/domains.json", tracksslUrl, agentToken)

	if c.DomainsUrl() != expected {
		t.Errorf("Expected %s, got %s", expected, c.DomainsUrl())
		return
	}

}

func TestClient_CertificateUrl(t *testing.T) {

	expected := fmt.Sprintf("%s/api/v1/agents/%s/certificate.json", tracksslUrl, agentToken)

	if c.CertificateUrl() != expected {
		t.Errorf("Expected %s, got %s", expected, c.CertificateUrl())
	}

}

func TestClient_SendCertSuccess(t *testing.T) {
	cert_bytes, err := os.ReadFile("testdata/morekidsonbikespa.org.crt")

	if err != nil {
		fmt.Errorf("Failed to read certificate: %v", err)
		return
	}

	r, err := recorder.New("fixtures/client/send_client_success")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	c.HttpClient = r.GetDefaultClient()

	domain := Domain{
		DomainID: 34027,
		Hostname: "morekidsonbikespa.org",
		Cert:     string(cert_bytes),
	}

	err = c.SendCert(&domain)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}

func TestClient_FetchDomainsSuccess(t *testing.T) {
	r, err := recorder.New("fixtures/client/fetch_domains_success")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	c.HttpClient = r.GetDefaultClient()

	domains, err := c.FetchDomains()

	if err != nil {
		t.Errorf("Expected domains, got %s", err)
	}

	if len(domains) != 3 {
		t.Errorf("Expected 3 domains, got %d", len(domains))
	}

}
