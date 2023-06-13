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
	"os"
	"testing"
)

func Test_Body_certs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		body      Body
		wantCount int
		wantErr   bool
	}{
		{
			name: "basic",
			body: Body{
				APIVersion: "0.0.1",
				Kind:       Kind,
				Spec: Spec{
					Data: Data{
						Hash: map[string]string{
							"algorithm": "sha256",
							"value":     "cd8327d867fce04bc97e149da50c3746340869575f7bf959a67284e34bfd46bc",
						},
					},
					Signature: Signature{
						Content: "MEYCIQDC3HNuKX8km/7TOo+HQ1WGr/FITRoYD9fw5RCRqc7WIgIhANAn6+yz0ziGGHH/52pDm5vCuOQgHnQ0AkGqi6POkEUB",
						PublicKey: PublicKey{
							//nolint:lll
							Content: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMwVENDQWxpZ0F3SUJBZ0lVZFh4WHRDa2NCdmt4UGhIc20yTzdxczJidnIwd0NnWUlLb1pJemowRUF3TXcKTnpFVk1CTUdBMVVFQ2hNTWMybG5jM1J2Y21VdVpHVjJNUjR3SEFZRFZRUURFeFZ6YVdkemRHOXlaUzFwYm5SbApjbTFsWkdsaGRHVXdIaGNOTWpNd05qRXpNVGN6T1RReVdoY05Nak13TmpFek1UYzBPVFF5V2pBQU1Ga3dFd1lICktvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVrUExsVkhGOFRURWM5YzMzbmxCWHFEelNqR3RPZUJoa0VXQ1AKNEw5c1hYdTQvNXdlYzN3eVY4Nkozazg5Ly9FVTQvSysreFI4ckpkZVFialpDUkFQaGFPQ0FYY3dnZ0Z6TUE0RwpBMVVkRHdFQi93UUVBd0lIZ0RBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREF6QWRCZ05WSFE0RUZnUVVWUk1MCklOUnpQZjBVTlF0T1RmVVpNMVNUUTNjd0h3WURWUjBqQkJnd0ZvQVUzOVBwejFZa0VaYjVxTmpwS0ZXaXhpNFkKWkQ4d0lRWURWUjBSQVFIL0JCY3dGWUVUYzNOamFISnZZMnRBWjI5dloyeGxMbU52YlRBc0Jnb3JCZ0VFQVlPLwpNQUVCQkI1b2RIUndjem92TDJkcGRHaDFZaTVqYjIwdmJHOW5hVzR2YjJGMWRHZ3dMZ1lLS3dZQkJBR0R2ekFCCkNBUWdEQjVvZEhSd2N6b3ZMMmRwZEdoMVlpNWpiMjB2Ykc5bmFXNHZiMkYxZEdnd2dZb0dDaXNHQVFRQjFua0MKQkFJRWZBUjZBSGdBZGdEZFBUQnF4c2NSTW1NWkhoeVpaemNDb2twZXVONDhyZitIaW5LQUx5bnVqZ0FBQVlpMQoxNDNUQUFBRUF3QkhNRVVDSUh3eXh1bHFJMm1JdU44ai94TlcrUmJxbGd0S2YzcVlxNzFjSDY4MWZFaHNBaUVBCmhlcVY3SVQzSzNkeEQxTzNXcElFczhsY1RhS1U4K2o1R2g3T1E4OVpOaTB3Q2dZSUtvWkl6ajBFQXdNRFp3QXcKWkFJd0JTT0R3YVJ1Q1hnY0VCRk1KQmh4OVYrSW92ZHdOOXp5SVJ6VWlvTlZrTHRQM281SHVvNVZjcExMcEI5bwpSUXNSQWpBVS9yRmFrU2dJamx4c2cwZWJJS280NVJwdGtDSkh6dUZWS3YzTDdMQitLczlRcnhWb2g4WmhhMVNhCjROemllNkU9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
						},
					},
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			certs, err := tt.body.Certs()
			if (err != nil) != tt.wantErr {
				t.Fatalf("(hashedRekordBody).certs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(certs) != tt.wantCount {
				t.Errorf("(hashedRekordBody).certs() parsed %d certs, wanted %d", len(certs), tt.wantCount)
			}
		})
	}
}

func Test_Body_Matches(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		blobPath string
		body     Body
		want     bool
	}{
		{
			name:     "valid match",
			blobPath: "./testdata/uploaded-blob.md",
			body: Body{
				Spec: Spec{
					Data: Data{
						Hash: map[string]string{
							"algorithm": "sha256",
							"value":     "cd8327d867fce04bc97e149da50c3746340869575f7bf959a67284e34bfd46bc",
						},
					},
				},
			},
			want: true,
		},
		{
			name:     "not a match, mismatched hash",
			blobPath: "./testdata/uploaded-blob.md",
			body: Body{
				Spec: Spec{
					Data: Data{
						Hash: map[string]string{
							"algorithm": "sha256",
							"value":     "foo",
						},
					},
				},
			},
			want: false,
		},
		{
			name:     "not sha256",
			blobPath: "./testdata/uploaded-blob.md",
			body: Body{
				Spec: Spec{
					Data: Data{
						Hash: map[string]string{
							"algorithm": "foo",
							"value":     "bar",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b, err := os.ReadFile(tt.blobPath)
			if err != nil {
				t.Fatalf("unable to read test file (%q): %v", tt.blobPath, err)
			}
			got := tt.body.Matches(b)
			if got != tt.want {
				t.Errorf("(hashedRekordBody).matches() got %t, wanted %t", got, tt.want)
			}
		})
	}
}
