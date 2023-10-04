// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package webpush

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateVAPIDKeys will create a private and public VAPID key pair
func GenerateVAPIDKeys() (privateKey, publicKey string, err error) {
	// Get the private key from the P256 curve
	curve := elliptic.P256()

	private, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return
	}

	public := elliptic.Marshal(curve, x, y)

	// Convert to base64
	publicKey = base64.RawURLEncoding.EncodeToString(public)
	privateKey = base64.RawURLEncoding.EncodeToString(private)

	return
}

// Generates the ECDSA public and private keys for the JWT encryption
func generateVAPIDHeaderKeys(privateKey []byte) *ecdsa.PrivateKey {
	// Public key
	curve := elliptic.P256()
	px, py := curve.ScalarMult(
		curve.Params().Gx,
		curve.Params().Gy,
		privateKey,
	)

	pubKey := ecdsa.PublicKey{
		Curve: curve,
		X:     px,
		Y:     py,
	}

	// Private key
	d := &big.Int{}
	d.SetBytes(privateKey)

	return &ecdsa.PrivateKey{
		PublicKey: pubKey,
		D:         d,
	}
}

// getVAPIDAuthorizationHeader
func getVAPIDAuthorizationHeader(
	endpoint,
	subscriber,
	vapidPublicKey,
	vapidPrivateKey string,
) (string, error) {
	// Create the JWT token
	subURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	sub := ""
	if strings.HasPrefix(subscriber, "https:") {
		sub = subscriber
	} else {
		sub = fmt.Sprintf("mailto:%s", subscriber)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"aud": fmt.Sprintf("%s://%s", subURL.Scheme, subURL.Host),
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"sub": sub,
	})

	// Decode the VAPID private key
	decodedVapidPrivateKey, err := decodeVapidKey(vapidPrivateKey)
	if err != nil {
		return "", err
	}

	privKey := generateVAPIDHeaderKeys(decodedVapidPrivateKey)

	// Sign token with private key
	jwtString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	// Decode the VAPID public key
	pubKey, err := decodeVapidKey(vapidPublicKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"vapid t=%s, k=%s",
		jwtString,
		base64.RawURLEncoding.EncodeToString(pubKey),
	), nil
}

// Need to decode the vapid private key in multiple base64 formats
// Solution from: https://github.com/SherClockHolmes/webpush-go/issues/29
func decodeVapidKey(key string) ([]byte, error) {
	bytes, err := base64.URLEncoding.DecodeString(key)
	if err == nil {
		return bytes, nil
	}

	return base64.RawURLEncoding.DecodeString(key)
}
