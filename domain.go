package main

import (
	"crypto/x509"
)

type Domain struct {
	DomainID int              `json:"domain_id"`
	Hostname string           `json:"hostname"`
	Port     string           `json:"port"`
	Cert     x509.Certificate `json:"-"`
}

func (d *Domain) String() string {
	return d.Hostname + ":" + d.Port
}

func (d *Domain) Issuer() string {
	issuer := d.Cert.Issuer
	org := issuer.Organization[0]

	if org != "" {
		return org
	}

	return d.Cert.Issuer.CommonName
}
