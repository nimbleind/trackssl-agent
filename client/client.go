package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreateResponse struct {
	Message string `json:"message"`
}

type Client struct {
	TracksslUrl string
	AuthToken   string
	AgentToken  string
	HttpClient  *http.Client
}

var (
	ERROR_CONNECTION_FAILED = fmt.Errorf("Connection failed")
)

func (c *Client) DomainsUrl() string {
	return fmt.Sprintf("%s/api/v1/agents/%s/domains.json", c.TracksslUrl, c.AgentToken)
}

func (c *Client) CertificateUrl() string {
	return fmt.Sprintf("%s/api/v1/agents/%s/certificate.json", c.TracksslUrl, c.AgentToken)
}

func (c *Client) SendCert(domain *Domain) error {
	http_client := &http.Client{}

	if c.HttpClient != nil {
		http_client = c.HttpClient
	}

	json_data, err := json.MarshalIndent(DomainRequest{Data: domain}, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal domain: %v", err)
	}
	req, err := http.NewRequest("POST", c.CertificateUrl(), bytes.NewBuffer(json_data))

	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
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
		return fmt.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "created" {
		return fmt.Errorf("Failed to create certificate: %s", response.Message)
	}

	return nil
}

func (c *Client) FetchDomains() ([]*Domain, error) {
	http_client := &http.Client{}

	if c.HttpClient != nil {
		http_client = c.HttpClient
	}

	req, err := http.NewRequest("GET", c.DomainsUrl(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	resp, err := http_client.Do(req)

	if err != nil {
		fmt.Println(fmt.Errorf("Failed to fetch domains from TrackSSL: %v", err))
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to fetch domains from TrackSSL: %d", resp.StatusCode)
	} else {
		body, _ := io.ReadAll(resp.Body)

		response := DomainResponse{}
		json.Unmarshal(body, &response)
		return response.Data, nil
	}
}
