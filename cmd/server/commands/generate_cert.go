// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package commands

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	host       *string        = new(string)
	validFrom  *string        = new(string)
	validFor   *time.Duration = new(time.Duration)
	isCA       *bool          = new(bool)
	rsaBits    *int           = new(int)
	ecdsaCurve *string        = new(string)
	ed25519Key *bool          = new(bool)

	generateCertCmd = &cobra.Command{
		Use:   "cert",
		Short: "Generate a self-signed X.509 certificate for a TLS server. Outputs to 'cert.pem' and 'key.pem' and will overwrite existing files.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(*host) == 0 {
				log.Fatalf("Missing required --host parameter")
			}

			var priv any
			var err error
			switch *ecdsaCurve {
			case "":
				if *ed25519Key {
					_, priv, err = ed25519.GenerateKey(rand.Reader)
				} else {
					priv, err = rsa.GenerateKey(rand.Reader, *rsaBits)
				}
			case "P224":
				priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
			case "P256":
				priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			case "P384":
				priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			case "P521":
				priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			default:
				log.Fatalf("Unrecognized elliptic curve: %q", *ecdsaCurve)
			}
			if err != nil {
				log.Fatalf("Failed to generate private key: %v", err)
			}

			// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
			// KeyUsage bits set in the x509.Certificate template
			keyUsage := x509.KeyUsageDigitalSignature
			// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
			// the context of TLS this KeyUsage is particular to RSA key exchange and
			// authentication.
			if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
				keyUsage |= x509.KeyUsageKeyEncipherment
			}

			var notBefore time.Time
			if len(*validFrom) == 0 {
				notBefore = time.Now()
			} else {
				notBefore, err = time.Parse("Jan 2 15:04:05 2006", *validFrom)
				if err != nil {
					log.Fatalf("Failed to parse creation date: %v", err)
				}
			}

			notAfter := notBefore.Add(*validFor)

			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
			if err != nil {
				log.Fatalf("Failed to generate serial number: %v", err)
			}

			template := x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Organization: []string{"Acme Co"},
				},
				NotBefore: notBefore,
				NotAfter:  notAfter,

				KeyUsage:              keyUsage,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
			}

			hosts := strings.Split(*host, ",")
			for _, h := range hosts {
				if ip := net.ParseIP(h); ip != nil {
					template.IPAddresses = append(template.IPAddresses, ip)
				} else {
					template.DNSNames = append(template.DNSNames, h)
				}
			}

			if *isCA {
				template.IsCA = true
				template.KeyUsage |= x509.KeyUsageCertSign
			}

			derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
			if err != nil {
				log.Fatalf("Failed to create certificate: %v", err)
			}

			certOut, err := os.Create("cert.pem")
			if err != nil {
				log.Fatalf("Failed to open cert.pem for writing: %v", err)
			}
			if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
				log.Fatalf("Failed to write data to cert.pem: %v", err)
			}
			if err := certOut.Close(); err != nil {
				log.Fatalf("Error closing cert.pem: %v", err)
			}
			fmt.Print("wrote cert.pem\n")

			keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Fatalf("Failed to open key.pem for writing: %v", err)
			}
			privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
			if err != nil {
				log.Fatalf("Unable to marshal private key: %v", err)
			}
			if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
				log.Fatalf("Failed to write data to key.pem: %v", err)
			}
			if err := keyOut.Close(); err != nil {
				log.Fatalf("Error closing key.pem: %v", err)
			}
			fmt.Print("wrote key.pem\n")

		},
	}
)

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func init() {
	generateCertCmd.Flags().StringVar(host, "host", "localhost", "Comma-separated hostnames and IPs to generate a certificate for")
	generateCertCmd.Flags().StringVar(validFrom, "start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	generateCertCmd.Flags().DurationVar(validFor, "duration", 365*24*time.Hour, "Duration that certificate is valid for")
	generateCertCmd.Flags().BoolVar(isCA, "ca", false, "whether this cert should be its own Certificate Authority")
	generateCertCmd.Flags().IntVar(rsaBits, "rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	generateCertCmd.Flags().StringVar(ecdsaCurve, "ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	generateCertCmd.Flags().BoolVar(ed25519Key, "ed25519", false, "Generate an Ed25519 key")
}
