// Copyright 2022 Security Scorecard Authors
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

package app

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"net/url"
	"testing"

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
							Path:   "https://test.com/foo/bar/workflow@c8416b0b2bf627c349ca92fc8e3de51a64b005cf",
						},
					},
				},
			},
			want: certInfo{
				repoFullName:  "https://test.com/",
				workflowPath:  "oo/bar/workflow",
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
