package client

type Domain struct {
	DomainID int    `json:"domain_id"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Cert    string `json:"certificate"`
	Error   string `json:"error"`
}

type DomainRequest struct {
	Data *Domain `json:"domain"`
}

type DomainResponse struct {
	Data []*Domain `json:"data"`
}

func (d *Domain) String() string {
	return d.Hostname + ":" + d.Port
}
