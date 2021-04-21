// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package alexa

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPError is a convenience method for logging a message and writing the provided error message
// and error code to the HTTP response.
func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Error(logMsg)
	}

	http.Error(w, err, errCode)
}

// IsValidAlexaRequest handles all the necessary steps to validate that an incoming http.Request has actually come from
// the Alexa service. If an error occurs during the validation process, an http.Error will be written to the provided http.ResponseWriter.
// The required steps for request validation can be found on this page:
// --insecure-skip-verify flag will disable all validations
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#hosting-a-custom-skill-as-a-web-service
func IsValidAlexaRequest(w http.ResponseWriter, r *http.Request) bool {

	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", 401)
		return false
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		HTTPError(w, err.Error(), "Not Authorized", 401)
		return false
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
		return false
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		HTTPError(w, "Amazon certificate invalid.", "Not Authorized", 401)
		return false
	}

	return true
}

func readCert(certURL string) ([]byte, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil || certPool == nil {
		log.Error("Can't open system cert pools")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: certPool, InsecureSkipVerify: insecureSkipVerify},
	}
	hc := &http.Client{Timeout: 2 * time.Second, Transport: tr}

	cert, err := hc.Get(certURL)
	if err != nil {
		return nil, errors.New("could not download Amazon cert file: " + err.Error())
	}
	defer cert.Body.Close()
	certContents, err := io.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("could not read Amazon cert file: " + err.Error())
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}
