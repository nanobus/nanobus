/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	jwt, err := generateJWT()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jwt)
}

func generateJWT() (string, error) {
	rsaPublicKeyBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return "", fmt.Errorf("cannot read private key file: %w", err)
	}
	pubPem, _ := pem.Decode(rsaPublicKeyBytes)
	if pubPem == nil {
		return "", fmt.Errorf("cannot decode private key file")
	}
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(pubPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(pubPem.Bytes); err != nil {
			return "", err
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("parsed key was not a RSA private key")
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(privateKey)
}
