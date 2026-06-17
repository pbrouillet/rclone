package main

import "testing"

func TestExtractCode(t *testing.T) {
	for _, tc := range []struct {
		name    string
		uri     string
		want    string
		wantErr bool
	}{
		{
			name: "OOBWithCode",
			uri:  "urn:ietf:wg:oauth:2.0:oob?code=ABC123&session_state=xyz",
			want: "ABC123",
		},
		{
			name: "NativeClientWithCode",
			uri:  "https://login.microsoftonline.com/common/oauth2/nativeclient?code=DEF456",
			want: "DEF456",
		},
		{
			name: "URLEncodedCode",
			uri:  "urn:ietf:wg:oauth:2.0:oob?code=a%2Fb%2Bc",
			want: "a/b+c",
		},
		{
			name:    "ErrorResponse",
			uri:     "urn:ietf:wg:oauth:2.0:oob?error=access_denied&error_description=user+cancelled",
			wantErr: true,
		},
		{
			name:    "NoCode",
			uri:     "urn:ietf:wg:oauth:2.0:oob",
			wantErr: true,
		},
		{
			name:    "Empty",
			uri:     "",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractCode(tc.uri)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got code %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
