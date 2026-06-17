package oauthutil

import (
	"net/url"
	"testing"

	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestConfig() *Config {
	return &Config{
		ClientID:     "test-client-id",
		ClientSecret: "",
		AuthURL:      "https://login.example.com/authorize",
		TokenURL:     "https://login.example.com/token",
		Scopes:       []string{"scope1", "offline_access"},
		RedirectURL:  TitleBarRedirectURL,
	}
}

// TestGetAuthURLPKCE checks that enabling PKCE adds an S256 code challenge to
// the authorization URL and stores the verifier in the config map.
func TestGetAuthURLPKCE(t *testing.T) {
	m := configmap.Simple{}
	opt := &Options{
		OAuth2Config: newTestConfig(),
		UsePKCE:      true,
	}

	authURL, state, err := getAuthURL("test", m, opt.OAuth2Config, opt)
	require.NoError(t, err)
	assert.NotEmpty(t, state)

	u, err := url.Parse(authURL)
	require.NoError(t, err)
	q := u.Query()
	assert.Equal(t, "S256", q.Get("code_challenge_method"))
	assert.NotEmpty(t, q.Get("code_challenge"))

	// The verifier must be persisted so configExchange can recover it.
	verifier, ok := m.Get(configPKCEVerifier)
	assert.True(t, ok)
	assert.NotEmpty(t, verifier)
}

// TestGetAuthURLNoPKCE checks that without PKCE no code challenge is added and
// no verifier is stored.
func TestGetAuthURLNoPKCE(t *testing.T) {
	m := configmap.Simple{}
	opt := &Options{
		OAuth2Config: newTestConfig(),
	}

	authURL, _, err := getAuthURL("test", m, opt.OAuth2Config, opt)
	require.NoError(t, err)

	u, err := url.Parse(authURL)
	require.NoError(t, err)
	q := u.Query()
	assert.Empty(t, q.Get("code_challenge"))
	assert.Empty(t, q.Get("code_challenge_method"))

	verifier, ok := m.Get(configPKCEVerifier)
	assert.False(t, ok)
	assert.Empty(t, verifier)
}
