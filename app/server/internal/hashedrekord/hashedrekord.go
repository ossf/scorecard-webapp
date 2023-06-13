// Copyright 2023 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hashedrekord

import (
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
)

const Kind = "hashedrekord"

// https://github.com/sigstore/rekor/blob/f01f9cd2c55eaddba9be28624fea793a26ad28c4/pkg/types/hashedrekord/v0.0.1/hashedrekord_v0_0_1_schema.json
//
//nolint:lll
type Body struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       Spec   `json:"spec"`
}
type Spec struct {
	Data      Data      `json:"data"`
	Signature Signature `json:"signature"`
}
type Data struct {
	Hash map[string]string
}
type Signature struct {
	Content   string    `json:"content"`
	PublicKey PublicKey `json:"publicKey"`
}
type PublicKey struct {
	Content string `json:"content"`
}

// check if the rekord object matches a given blob (currently compares sha256 hash).
func (b Body) Matches(blob []byte) bool {
	if b.Spec.Data.Hash["algorithm"] != "sha256" {
		log.Println("hashed rekord entry has no sha256")
		return false
	}
	sha := sha256.Sum256(blob)
	have := hex.EncodeToString(sha[:])
	want := b.Spec.Data.Hash["value"]
	return have == want
}

// extracts x509 certs from the hashedrekord tlog entry.
// It uses the public key to pem decode the certificates.
func (b Body) Certs() ([]*x509.Certificate, error) {
	publicKey, err := base64.StdEncoding.DecodeString(b.Spec.Signature.PublicKey.Content)
	if err != nil {
		return nil, fmt.Errorf("decode rekord public key: %w", err)
	}

	remaining := publicKey
	var result []*x509.Certificate
	for len(remaining) > 0 {
		var certDer *pem.Block
		certDer, remaining = pem.Decode(remaining)
		if certDer == nil {
			return nil, fmt.Errorf("error during PEM decoding: %w", err)
		}

		cert, err := x509.ParseCertificate(certDer.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error during certificate parsing: %w", err)
		}
		result = append(result, cert)
	}
	return result, nil
}
