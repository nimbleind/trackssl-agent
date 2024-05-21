package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var DEFAULT_TRACKSSL_URL = "https://app.trackssl.com"

type CreateResponse struct {
	Message string `json:"message"`
}

type Client struct {
	TracksslUrl string
	AuthToken   string
	AgentToken  string
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

func (c *Client) FetchDomains() ([]*Domain, error) {
	http_client := &http.Client{}

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
