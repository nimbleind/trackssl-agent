package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var DEFAULT_TRACKSSL_URL = "https://home.andylibby.org"

type CreateResponse struct {
	Message string `json:"message"`
}
type Client struct {
	TracksslUrl string
	AuthToken   string
	AgentToken  string
}

var (
	ERROR_NO_AUTH_TOKEN     = fmt.Errorf("No auth token")

	ERROR_NO_AGENT_TOKEN    = fmt.Errorf("No agent token")
	ERROR_CONNECTION_FAILED = fmt.Errorf("Connection failed")
)

func (c *Client) checkValues() error {
	if c.AuthToken == "" {
		return ERROR_NO_AUTH_TOKEN
	}

	if c.AgentToken == "" {
		return ERROR_NO_AGENT_TOKEN
	}

	return nil
}

func NewClient() (*Client, error) {
	client := &Client{
		TracksslUrl: os.Getenv("TRACKSSL_URL"),
		AuthToken:   os.Getenv("TRACKSSL_AUTH_TOKEN"),
		AgentToken:  os.Getenv("TRACKSSL_AGENT_TOKEN"),
	}

	if client.TracksslUrl == "" {
		client.TracksslUrl = DEFAULT_TRACKSSL_URL
	}

	if err := client.checkValues(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) DomainsUrl() string {
	return fmt.Sprintf("%s/api/v1/agents/%s/domains.json", c.TracksslUrl, c.AgentToken)
}

func (c *Client) CertificateUrl() string {
	return fmt.Sprintf("%s/api/v1/agents/%s/certificate.json", c.TracksslUrl, c.AgentToken)
}

func (c *Client) SendCert(domain *Domain) error {
	http_client := &http.Client{}
	json_data, _ := json.MarshalIndent(DomainRequest{Data: domain}, "", "  ")
	req, _ := http.NewRequest("POST", c.CertificateUrl(), bytes.NewBuffer(json_data))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http_client.Do(req)

	if err != nil {
		return ERROR_CONNECTION_FAILED
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	response := CreateResponse{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return err
	}

	if response.Message != "created" {
		return fmt.Errorf("Failed to create certificate: %s", response.Message)
	}

	return nil
}

func (c *Client) FetchDomains() []*Domain {
	http_client := &http.Client{}

	req, err := http.NewRequest("GET", c.DomainsUrl(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := http_client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	response := DomainResponse{}
	json.Unmarshal(body, &response)
	return response.Data
}
