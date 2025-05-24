// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package autocert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/caddyserver/certmagic"
	"github.com/libdns/cloudflare"
	"github.com/mholt/acmez/v2"
	"github.com/mholt/acmez/v2/acme"
	"go.uber.org/atomic"
)

type Autocert struct {
	domains            []string
	ownerEmail         []string
	cloudflareAPIToken string
	privateKey         *ecdsa.PrivateKey
	certs              []acme.Certificate
	inProcess          *atomic.Bool
	prod               bool
}

func NewAutocert(domains []string, ownerEmail []string, cloudflareAPIToken string, prod bool) *Autocert {
	return &Autocert{
		domains:            domains,
		ownerEmail:         ownerEmail,
		cloudflareAPIToken: cloudflareAPIToken,
		inProcess:          atomic.NewBool(false),
		prod:               prod,
	}
}

func (a *Autocert) Start() {}

func (a *Autocert) Shutdown() {}

func (a *Autocert) RequestCertificate(ctx context.Context) (err error) {
	if !a.inProcess.CompareAndSwap(false, true) {
		return
	}
	defer a.inProcess.Store(false)

	log.Infof("request certificate for %v", a.domains)

	//fmt.Println(a.domains, len(a.domains), a.ownerEmail, len(a.ownerEmail), a.cloudflareAPIToken)

	// Initialize a DNS-01 solver, using Cloudflare APIs
	solver := &certmagic.DNS01Solver{
		DNSManager: certmagic.DNSManager{
			DNSProvider: &cloudflare.Provider{
				APIToken: a.cloudflareAPIToken,
			},
		},
	}
	// The CA endpoint to use (prod or staging)
	// switch to Production once fully tested
	// otherwise you might get rate-limited in Production
	// before you've had a chance to test that your client
	// works as expected
	caLocation := certmagic.LetsEncryptStagingCA
	if a.prod {
		caLocation = certmagic.LetsEncryptProductionCA
	}

	// Initialize an acmez client
	client := acmez.Client{
		Client: &acme.Client{
			Directory: caLocation,
			Logger:    log.Loggert(),
		},
		ChallengeSolvers: map[string]acmez.Solver{
			acme.ChallengeTypeDNS01: solver,
		},
	}

	// Generate a private key for your Let's Encrypt account
	var accountPrivateKey *ecdsa.PrivateKey
	accountPrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		err = fmt.Errorf("%s: %w", "could not generate an account key", err)
		return
	}
	// Create a Let's Encrypt account
	account := acme.Account{
		Contact:              a.ownerEmail,
		TermsOfServiceAgreed: true,
		PrivateKey:           accountPrivateKey,
	}

	var acc acme.Account
	acc, err = client.NewAccount(ctx, account)
	if err != nil {
		err = fmt.Errorf("%s: %w", "could not create new account:", err)
		return
	}

	// Generate a private key for the certificate
	a.privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		err = fmt.Errorf("%s: %w", "generating certificate key", err)
		return
	}

	a.certs, err = client.ObtainCertificateForSANs(ctx, acc, a.privateKey, a.domains)
	if err != nil {
		err = fmt.Errorf("%s: %w", "could not obtain certificate", err)
		return
	}

	// since the client returns more than one cert, it is up to you
	// to choose the most appropriate one (such as one which contains
	// the full chain, including any intermediate certificates)
	//for _, cert := range a.certs {
	//	fmt.Println(string(cert.ChainPEM))
	//}

	log.Info("successful")

	return nil
}

func (a *Autocert) PrivateKey() (pemEncoded []byte) {
	x509Encoded, _ := x509.MarshalECPrivateKey(a.privateKey)
	pemEncoded = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return
}

func (a *Autocert) PublicKey() (pemEncodedPub []byte) {
	for _, cert := range a.certs {
		pemEncodedPub = append(pemEncodedPub, cert.ChainPEM...)
	}
	return
}
