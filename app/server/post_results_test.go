// Copyright 2022 OpenSSF Scorecard Authors
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

package server

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"net/url"
	"testing"
	"unicode/utf8"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"github.com/stretchr/testify/assert"
)

func Test_extractCertInfo(t *testing.T) {
	t.Parallel()
	type args struct {
		cert *x509.Certificate
	}
	tests := []struct {
		name    string
		args    args
		want    certInfo
		errType error
	}{
		{
			name: "certificate has no URIs",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
				},
			},
			want:    certInfo{},
			errType: errCertMissingURI,
		},
		{
			name: "cert has empty repository ref fulcioRepoRefKey",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
					Extensions: []pkix.Extension{
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 6},
						},
					},
				},
			},
			errType: errEmptyCertRef,
		},
		{
			name: "cert has empty repository path for fulcioRepoPathKey",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
					Extensions: []pkix.Extension{
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 5},
						},
					},
				},
			},
			errType: errEmptyCertPath,
		},
		{
			name: "Cert Missing URI",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
					Extensions: []pkix.Extension{
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 5},
							Value:    fulcioIntermediate,
						},
					},
				},
			},
			errType: errCertMissingURI,
		},
		{
			name: "Cert Workflow path is empty",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
					Extensions: []pkix.Extension{
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 5},
							Value:    fulcioIntermediate,
						},
					},
					URIs: []*url.URL{
						{
							Scheme: "https",
							Host:   "test.com",
						},
					},
				},
			},
			errType: errCertWorkflowPathEmpty,
		},
		{
			name: "Valid Cert",
			args: args{
				cert: &x509.Certificate{
					Subject: pkix.Name{
						CommonName: "test",
					},
					Extensions: []pkix.Extension{
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 5},
							Value:    []byte("https://test.com/"),
						},
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 6},
							Value:    []byte("https://test.com/"),
						},
						{
							Critical: true,
							Id:       asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 3},
							Value:    []byte("https://test.com/"),
						},
					},
					URIs: []*url.URL{
						{
							Scheme: "https",
							Host:   "test.com",
							Path:   "/foo/bar/workflow@c8416b0b2bf627c349ca92fc8e3de51a64b005cf",
						},
					},
				},
			},
			want: certInfo{
				repoFullName:  "https://test.com/",
				workflowPath:  "foo/bar/workflow",
				workflowRef:   "c8416b0b2bf627c349ca92fc8e3de51a64b005cf",
				repoBranchRef: "https://test.com/",
				repoSHA:       "https://test.com/",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := extractCertInfo(tt.args.cert)
			if err != nil {
				assert.Equal(t, tt.errType, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func FuzzExtractCertInfo(f *testing.F) {
	f.Fuzz(func(t *testing.T, commonName, value string, critical bool, data []byte) {
		f := fuzz.NewConsumer(data)
		asn := []int{}
		f.CreateSlice(&asn)
		if !utf8.ValidString(commonName) || !utf8.ValidString(value) {
			t.Skip()
		}
		if len(data) == 0 {
			t.Skip()
		}
		if len(asn) < 8 {
			t.Skip()
		}
		cert := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: commonName,
			},
			Extensions: []pkix.Extension{
				{
					Critical: critical,
					Id:       asn1.ObjectIdentifier{asn[0], asn[1], asn[2], asn[3], asn[4], asn[5], asn[6], asn[7]},
					Value:    []byte(value),
				},
			},
		}
		extractCertInfo(cert)
	})
}

func Test_splitFullPath(t *testing.T) {
	t.Parallel()
	type results struct {
		org, repo, subPath string
		ok                 bool
	}
	tests := []struct {
		name string
		path string
		want results
	}{
		{
			name: "valid path",
			path: "org/repo/rest/of/path@ref",
			want: results{
				org:     "org",
				repo:    "repo",
				subPath: "rest/of/path@ref",
				ok:      true,
			},
		},
		{
			name: "malformed path",
			path: "malformed/path",
			want: results{
				org:     "",
				repo:    "",
				subPath: "",
				ok:      false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o, r, p, ok := splitFullPath(tt.path)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.org, o)
			assert.Equal(t, tt.want.repo, r)
			assert.Equal(t, tt.want.subPath, p)
		})
	}
}

// Test_getCertInfoFromCert tests the getCertInfoFromCert function
func Test_getCertPool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		cert        []byte
		shouldBeNil bool // not comparing the cert pool, just if it's nil or not. It is hard to compare the cert pool
		wantError   bool
	}{
		{nil, true, true},
		{[]byte(""), true, true},
		{[]byte("certificate"), true, true},
		{[]byte(
			`-----BEGIN CERTIFICATE-----
MIIDYzCCAkugAwIBAgIRAPgUfht89Rg0uiZ/sUqqM+swDQYJKoZIhvcNAQELBQAw
UDEhMB8GA1UEChMYTWFubmluZyBQdWJsaWNhdGlvbnMgQ28uMQ4wDAYDVQQLEwVC
b29rczEbMBkGA1UEAxMSR28gV2ViIFByb2dyYW1taW5nMB4XDTIzMDEwODE2Mjkx
MloXDTI1MTAwNDE2MjkxMlowUDEhMB8GA1UEChMYTWFubmluZyBQdWJsaWNhdGlv
bnMgQ28uMQ4wDAYDVQQLEwVCb29rczEbMBkGA1UEAxMSR28gV2ViIFByb2dyYW1t
aW5nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz06OXOWcXUOihEuJ
yMDyGGAzGVUmLbLen4q2LiW8WiHBgM4gbADUZzK6d42XsNFQuptbZsjKl6/P0z4A
LmALIXHwln4Acju5IbkqdQbKxyqbwHILNyXUoOXmbVwQwXwNxdKsRi2nySzsNUR7
JiUhYZKlL26zu1wRo3JjnViXLsecf/1G+nvWyypcvf1iRmLv6oaMqroVAfnmksnG
Oi2Peyd5fDbuMF9nI9qqPZ//clgPa9Z6P4kF+r+VAdgFo2Kt2mUzf8sn5jSvmdTu
gYkuRoQdn6bQYY68FMEJivBXOPIwRqdbrnF7wl+PNEuZVYQgeSh5uNX1+WO1DD/y
sRRm4QIDAQABozgwNjAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUH
AwEwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAeaDNcWWiOY0r
S6gr3AFhbPVD0sFDna+BbDfVOL6u4a9XVw2bg0ul0wEZFEq2g4qHASJDP6aB0og+
rZCXPdmRvWKa1J1UsBBXwrZ3AEYv1kqOU2GJhD5AlG+zASAym7InapCQ5yU4eVB4
5n0dayWStHI5It+ub2ubDcpZvjB+kCTRRLAy7PSSxa2rY7csIYgEOALOJk1VqO2M
VT47afTMFrgOFyZ33DArNO034Dnu+Uz/V1SAeebPGv0vdl65plLh8ekoLgBZW87k
e0i4463IxAwWdpCk29FOn3o0GZAtCWhDznIM70bTunZxl6QRCjdN0Z2sDEl+jPit
MYSKu39B6Q==
-----END CERTIFICATE-----`), false, false},
	}

	for _, tt := range tests {
		t.Parallel()
		got, err := getCertPool(tt.cert)
		if (err != nil) != tt.wantError {
			t.Errorf("getCertPool() error = %v, wantErr %v", err, tt.wantError)
			continue
		}

		if tt.shouldBeNil {
			assert.Nil(t, got)
		} else {
			assert.NotNil(t, got)
		}
	}
}
