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
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"unicode/utf8"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/ossf/scorecard-webapp/app/server/internal/hashedrekord"
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

// Test_getCertInfoFromCert tests the getCertInfoFromCert function.
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
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := getCertPool(tt.cert)
			if (err != nil) != tt.wantError {
				t.Errorf("getCertPool() error = %v, wantErr %v", err, tt.wantError)
			}

			if tt.shouldBeNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_getTLogEntryFromURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		responsePath string
		uuid         string
		entry        *tlogEntry
		wantErr      bool
	}{
		{
			name:         "valid entry",
			responsePath: "/rekor/log-entries-response.json",
			uuid:         "24296fb24b8ad77abaa457505061c4a0ef34197534bdd8b474acfafdc4d76c437726e98153c7b253",
			entry: &tlogEntry{
				//nolint:lll
				Body:           "eyJhcGlWZXJzaW9uIjoiMC4wLjEiLCJraW5kIjoiaGFzaGVkcmVrb3JkIiwic3BlYyI6eyJkYXRhIjp7Imhhc2giOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiJjZDgzMjdkODY3ZmNlMDRiYzk3ZTE0OWRhNTBjMzc0NjM0MDg2OTU3NWY3YmY5NTlhNjcyODRlMzRiZmQ0NmJjIn19LCJzaWduYXR1cmUiOnsiY29udGVudCI6Ik1FWUNJUURDM0hOdUtYOGttLzdUT28rSFExV0dyL0ZJVFJvWUQ5Znc1UkNScWM3V0lnSWhBTkFuNit5ejB6aUdHSEgvNTJwRG01dkN1T1FnSG5RMEFrR3FpNlBPa0VVQiIsInB1YmxpY0tleSI6eyJjb250ZW50IjoiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTXdWRU5EUVd4cFowRjNTVUpCWjBsVlpGaDRXSFJEYTJOQ2RtdDRVR2hJYzIweVR6ZHhjekppZG5Jd2QwTm5XVWxMYjFwSmVtb3dSVUYzVFhjS1RucEZWazFDVFVkQk1WVkZRMmhOVFdNeWJHNWpNMUoyWTIxVmRWcEhWakpOVWpSM1NFRlpSRlpSVVVSRmVGWjZZVmRrZW1SSE9YbGFVekZ3WW01U2JBcGpiVEZzV2tkc2FHUkhWWGRJYUdOT1RXcE5kMDVxUlhwTlZHTjZUMVJSZVZkb1kwNU5hazEzVG1wRmVrMVVZekJQVkZGNVYycEJRVTFHYTNkRmQxbElDa3R2V2tsNmFqQkRRVkZaU1V0dldrbDZhakJFUVZGalJGRm5RVVZyVUV4c1ZraEdPRlJVUldNNVl6TXpibXhDV0hGRWVsTnFSM1JQWlVKb2EwVlhRMUFLTkV3NWMxaFlkVFF2TlhkbFl6TjNlVlk0TmtvemF6ZzVMeTlGVlRRdlN5c3JlRkk0Y2twa1pWRmlhbHBEVWtGUWFHRlBRMEZZWTNkblowWjZUVUUwUndwQk1WVmtSSGRGUWk5M1VVVkJkMGxJWjBSQlZFSm5UbFpJVTFWRlJFUkJTMEpuWjNKQ1owVkdRbEZqUkVGNlFXUkNaMDVXU0ZFMFJVWm5VVlZXVWsxTUNrbE9VbnBRWmpCVlRsRjBUMVJtVlZwTk1WTlVVVE5qZDBoM1dVUldVakJxUWtKbmQwWnZRVlV6T1ZCd2VqRlphMFZhWWpWeFRtcHdTMFpYYVhocE5Ga0tXa1E0ZDBsUldVUldVakJTUVZGSUwwSkNZM2RHV1VWVVl6Tk9hbUZJU25aWk1uUkJXakk1ZGxveWVHeE1iVTUyWWxSQmMwSm5iM0pDWjBWRlFWbFBMd3BOUVVWQ1FrSTFiMlJJVW5kamVtOTJUREprY0dSSGFERlphVFZxWWpJd2RtSkhPVzVoVnpSMllqSkdNV1JIWjNkTVoxbExTM2RaUWtKQlIwUjJla0ZDQ2tOQlVXZEVRalZ2WkVoU2QyTjZiM1pNTW1Sd1pFZG9NVmxwTldwaU1qQjJZa2M1Ym1GWE5IWmlNa1l4WkVkbmQyZFpiMGREYVhOSFFWRlJRakZ1YTBNS1FrRkpSV1pCVWpaQlNHZEJaR2RFWkZCVVFuRjRjMk5TVFcxTldraG9lVnBhZW1ORGIydHdaWFZPTkRoeVppdElhVzVMUVV4NWJuVnFaMEZCUVZscE1Rb3hORE5VUVVGQlJVRjNRa2hOUlZWRFNVaDNlWGgxYkhGSk1tMUpkVTQ0YWk5NFRsY3JVbUp4YkdkMFMyWXpjVmx4TnpGalNEWTRNV1pGYUhOQmFVVkJDbWhsY1ZZM1NWUXpTek5rZUVReFR6TlhjRWxGY3poc1kxUmhTMVU0SzJvMVIyZzNUMUU0T1ZwT2FUQjNRMmRaU1V0dldrbDZhakJGUVhkTlJGcDNRWGNLV2tGSmQwSlRUMFIzWVZKMVExaG5ZMFZDUmsxS1FtaDRPVllyU1c5MlpIZE9PWHA1U1ZKNlZXbHZUbFpyVEhSUU0yODFTSFZ2TlZaamNFeE1jRUk1YndwU1VYTlNRV3BCVlM5eVJtRnJVMmRKYW14NGMyY3daV0pKUzI4ME5WSndkR3REU2toNmRVWldTM1l6VERkTVFpdExjemxSY25oV2IyZzRXbWhoTVZOaENqUk9lbWxsTmtVOUNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSyJ9fX19",
				IntegratedTime: 1686677983,
				LogID:          "c0d23d6ad406973f9559f3ba2d1ca01f84147d8ffc5b8445c224f98b9591801d",
				LogIndex:       23652179,
			},
			wantErr: false,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			url := setupServer(t) + tt.responsePath
			uuid, entry, err := getTLogEntryFromURL(ctx, url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getTLogEntryFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if uuid != tt.uuid {
				t.Errorf("getTLogEntryFromURL() uuid: %s, wanted %s", uuid, tt.uuid)
			}
			ignoreVerification := cmpopts.IgnoreFields(tlogEntry{}, "Verification")
			if !cmp.Equal(entry, tt.entry, ignoreVerification) {
				t.Error(cmp.Diff(entry, tt.entry, ignoreVerification))
			}
		})
	}
}

func setupServer(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("./testdata" + r.URL.Path)
		if err != nil {
			t.Logf("os.ReadFile: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func Test_tlogEntry_rekord(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		entry   tlogEntry
		want    hashedrekord.Body
		wantErr bool
	}{
		{
			name: "basic",
			entry: tlogEntry{
				//nolint:lll
				Body:           "eyJhcGlWZXJzaW9uIjoiMC4wLjEiLCJraW5kIjoiaGFzaGVkcmVrb3JkIiwic3BlYyI6eyJkYXRhIjp7Imhhc2giOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiJjZDgzMjdkODY3ZmNlMDRiYzk3ZTE0OWRhNTBjMzc0NjM0MDg2OTU3NWY3YmY5NTlhNjcyODRlMzRiZmQ0NmJjIn19LCJzaWduYXR1cmUiOnsiY29udGVudCI6Ik1FWUNJUURDM0hOdUtYOGttLzdUT28rSFExV0dyL0ZJVFJvWUQ5Znc1UkNScWM3V0lnSWhBTkFuNit5ejB6aUdHSEgvNTJwRG01dkN1T1FnSG5RMEFrR3FpNlBPa0VVQiIsInB1YmxpY0tleSI6eyJjb250ZW50IjoiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTXdWRU5EUVd4cFowRjNTVUpCWjBsVlpGaDRXSFJEYTJOQ2RtdDRVR2hJYzIweVR6ZHhjekppZG5Jd2QwTm5XVWxMYjFwSmVtb3dSVUYzVFhjS1RucEZWazFDVFVkQk1WVkZRMmhOVFdNeWJHNWpNMUoyWTIxVmRWcEhWakpOVWpSM1NFRlpSRlpSVVVSRmVGWjZZVmRrZW1SSE9YbGFVekZ3WW01U2JBcGpiVEZzV2tkc2FHUkhWWGRJYUdOT1RXcE5kMDVxUlhwTlZHTjZUMVJSZVZkb1kwNU5hazEzVG1wRmVrMVVZekJQVkZGNVYycEJRVTFHYTNkRmQxbElDa3R2V2tsNmFqQkRRVkZaU1V0dldrbDZhakJFUVZGalJGRm5RVVZyVUV4c1ZraEdPRlJVUldNNVl6TXpibXhDV0hGRWVsTnFSM1JQWlVKb2EwVlhRMUFLTkV3NWMxaFlkVFF2TlhkbFl6TjNlVlk0TmtvemF6ZzVMeTlGVlRRdlN5c3JlRkk0Y2twa1pWRmlhbHBEVWtGUWFHRlBRMEZZWTNkblowWjZUVUUwUndwQk1WVmtSSGRGUWk5M1VVVkJkMGxJWjBSQlZFSm5UbFpJVTFWRlJFUkJTMEpuWjNKQ1owVkdRbEZqUkVGNlFXUkNaMDVXU0ZFMFJVWm5VVlZXVWsxTUNrbE9VbnBRWmpCVlRsRjBUMVJtVlZwTk1WTlVVVE5qZDBoM1dVUldVakJxUWtKbmQwWnZRVlV6T1ZCd2VqRlphMFZhWWpWeFRtcHdTMFpYYVhocE5Ga0tXa1E0ZDBsUldVUldVakJTUVZGSUwwSkNZM2RHV1VWVVl6Tk9hbUZJU25aWk1uUkJXakk1ZGxveWVHeE1iVTUyWWxSQmMwSm5iM0pDWjBWRlFWbFBMd3BOUVVWQ1FrSTFiMlJJVW5kamVtOTJUREprY0dSSGFERlphVFZxWWpJd2RtSkhPVzVoVnpSMllqSkdNV1JIWjNkTVoxbExTM2RaUWtKQlIwUjJla0ZDQ2tOQlVXZEVRalZ2WkVoU2QyTjZiM1pNTW1Sd1pFZG9NVmxwTldwaU1qQjJZa2M1Ym1GWE5IWmlNa1l4WkVkbmQyZFpiMGREYVhOSFFWRlJRakZ1YTBNS1FrRkpSV1pCVWpaQlNHZEJaR2RFWkZCVVFuRjRjMk5TVFcxTldraG9lVnBhZW1ORGIydHdaWFZPTkRoeVppdElhVzVMUVV4NWJuVnFaMEZCUVZscE1Rb3hORE5VUVVGQlJVRjNRa2hOUlZWRFNVaDNlWGgxYkhGSk1tMUpkVTQ0YWk5NFRsY3JVbUp4YkdkMFMyWXpjVmx4TnpGalNEWTRNV1pGYUhOQmFVVkJDbWhsY1ZZM1NWUXpTek5rZUVReFR6TlhjRWxGY3poc1kxUmhTMVU0SzJvMVIyZzNUMUU0T1ZwT2FUQjNRMmRaU1V0dldrbDZhakJGUVhkTlJGcDNRWGNLV2tGSmQwSlRUMFIzWVZKMVExaG5ZMFZDUmsxS1FtaDRPVllyU1c5MlpIZE9PWHA1U1ZKNlZXbHZUbFpyVEhSUU0yODFTSFZ2TlZaamNFeE1jRUk1YndwU1VYTlNRV3BCVlM5eVJtRnJVMmRKYW14NGMyY3daV0pKUzI4ME5WSndkR3REU2toNmRVWldTM1l6VERkTVFpdExjemxSY25oV2IyZzRXbWhoTVZOaENqUk9lbWxsTmtVOUNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSyJ9fX19",
				IntegratedTime: 1686677983,
				LogID:          "c0d23d6ad406973f9559f3ba2d1ca01f84147d8ffc5b8445c224f98b9591801d",
				LogIndex:       23652179,
			},
			want: hashedrekord.Body{
				APIVersion: "0.0.1",
				Kind:       hashedrekord.Kind,
				Spec: hashedrekord.Spec{
					Data: hashedrekord.Data{
						Hash: hashedrekord.Hash{
							Algorithm: "sha256",
							Value:     "cd8327d867fce04bc97e149da50c3746340869575f7bf959a67284e34bfd46bc",
						},
					},
					Signature: hashedrekord.Signature{
						Content: "MEYCIQDC3HNuKX8km/7TOo+HQ1WGr/FITRoYD9fw5RCRqc7WIgIhANAn6+yz0ziGGHH/52pDm5vCuOQgHnQ0AkGqi6POkEUB",
						PublicKey: hashedrekord.PublicKey{
							//nolint:lll
							Content: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMwVENDQWxpZ0F3SUJBZ0lVZFh4WHRDa2NCdmt4UGhIc20yTzdxczJidnIwd0NnWUlLb1pJemowRUF3TXcKTnpFVk1CTUdBMVVFQ2hNTWMybG5jM1J2Y21VdVpHVjJNUjR3SEFZRFZRUURFeFZ6YVdkemRHOXlaUzFwYm5SbApjbTFsWkdsaGRHVXdIaGNOTWpNd05qRXpNVGN6T1RReVdoY05Nak13TmpFek1UYzBPVFF5V2pBQU1Ga3dFd1lICktvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVrUExsVkhGOFRURWM5YzMzbmxCWHFEelNqR3RPZUJoa0VXQ1AKNEw5c1hYdTQvNXdlYzN3eVY4Nkozazg5Ly9FVTQvSysreFI4ckpkZVFialpDUkFQaGFPQ0FYY3dnZ0Z6TUE0RwpBMVVkRHdFQi93UUVBd0lIZ0RBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREF6QWRCZ05WSFE0RUZnUVVWUk1MCklOUnpQZjBVTlF0T1RmVVpNMVNUUTNjd0h3WURWUjBqQkJnd0ZvQVUzOVBwejFZa0VaYjVxTmpwS0ZXaXhpNFkKWkQ4d0lRWURWUjBSQVFIL0JCY3dGWUVUYzNOamFISnZZMnRBWjI5dloyeGxMbU52YlRBc0Jnb3JCZ0VFQVlPLwpNQUVCQkI1b2RIUndjem92TDJkcGRHaDFZaTVqYjIwdmJHOW5hVzR2YjJGMWRHZ3dMZ1lLS3dZQkJBR0R2ekFCCkNBUWdEQjVvZEhSd2N6b3ZMMmRwZEdoMVlpNWpiMjB2Ykc5bmFXNHZiMkYxZEdnd2dZb0dDaXNHQVFRQjFua0MKQkFJRWZBUjZBSGdBZGdEZFBUQnF4c2NSTW1NWkhoeVpaemNDb2twZXVONDhyZitIaW5LQUx5bnVqZ0FBQVlpMQoxNDNUQUFBRUF3QkhNRVVDSUh3eXh1bHFJMm1JdU44ai94TlcrUmJxbGd0S2YzcVlxNzFjSDY4MWZFaHNBaUVBCmhlcVY3SVQzSzNkeEQxTzNXcElFczhsY1RhS1U4K2o1R2g3T1E4OVpOaTB3Q2dZSUtvWkl6ajBFQXdNRFp3QXcKWkFJd0JTT0R3YVJ1Q1hnY0VCRk1KQmh4OVYrSW92ZHdOOXp5SVJ6VWlvTlZrTHRQM281SHVvNVZjcExMcEI5bwpSUXNSQWpBVS9yRmFrU2dJamx4c2cwZWJJS280NVJwdGtDSkh6dUZWS3YzTDdMQitLczlRcnhWb2g4WmhhMVNhCjROemllNkU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, err := tt.entry.rekord()
			if (err != nil) != tt.wantErr {
				t.Fatalf("(tlogEntry).rekord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(r, tt.want) {
				t.Error(cmp.Diff(r, tt.want))
			}
		})
	}
}
