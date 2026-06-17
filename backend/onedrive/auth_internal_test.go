package onedrive

import (
	"context"
	"testing"

	"github.com/rclone/rclone/lib/oauthutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakeOauthConfigWebAuth checks that web_auth produces a first-party
// public-client config with the SharePoint resource audience and OOB redirect.
func TestMakeOauthConfigWebAuth(t *testing.T) {
	ctx := context.Background()

	t.Run("Defaults", func(t *testing.T) {
		opt := &Options{
			Region:          regionGlobal,
			WebAuth:         true,
			WebAuthClientID: webAuthClientID,
			TenantURL:       "https://contoso-my.sharepoint.com/_api",
			AccessScopes:    scopeAccess,
		}
		conf, err := makeOauthConfig(ctx, opt)
		require.NoError(t, err)

		assert.Equal(t, webAuthClientID, conf.ClientID)
		assert.Empty(t, conf.ClientSecret)
		assert.Equal(t, oauthutil.TitleBarRedirectURL, conf.RedirectURL)
		assert.Equal(t, []string{"https://contoso-my.sharepoint.com/.default", "offline_access"}, []string(conf.Scopes))
		assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/authorize", conf.AuthURL)
		assert.Equal(t, "https://login.microsoftonline.com/common/oauth2/v2.0/token", conf.TokenURL)
	})

	t.Run("CustomClientIDAndSharePointSite", func(t *testing.T) {
		opt := &Options{
			Region:          regionGlobal,
			WebAuth:         true,
			WebAuthClientID: "11111111-2222-3333-4444-555555555555",
			TenantURL:       "https://contoso.sharepoint.com/sites/team/_api",
			AccessScopes:    scopeAccess,
		}
		conf, err := makeOauthConfig(ctx, opt)
		require.NoError(t, err)

		assert.Equal(t, "11111111-2222-3333-4444-555555555555", conf.ClientID)
		assert.Equal(t, []string{"https://contoso.sharepoint.com/.default", "offline_access"}, []string(conf.Scopes))
	})

	t.Run("EmptyClientIDFallsBackToDefault", func(t *testing.T) {
		opt := &Options{
			Region:       regionGlobal,
			WebAuth:      true,
			TenantURL:    "https://contoso-my.sharepoint.com/_api",
			AccessScopes: scopeAccess,
		}
		conf, err := makeOauthConfig(ctx, opt)
		require.NoError(t, err)
		assert.Equal(t, webAuthClientID, conf.ClientID)
	})

	t.Run("MissingTenantURLFails", func(t *testing.T) {
		opt := &Options{
			Region:       regionGlobal,
			WebAuth:      true,
			AccessScopes: scopeAccess,
		}
		_, err := makeOauthConfig(ctx, opt)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tenant_url")
	})

	t.Run("InvalidTenantURLFails", func(t *testing.T) {
		opt := &Options{
			Region:       regionGlobal,
			WebAuth:      true,
			TenantURL:    "not-a-url",
			AccessScopes: scopeAccess,
		}
		_, err := makeOauthConfig(ctx, opt)
		require.Error(t, err)
	})
}

// TestMakeOauthConfigDefaultUnaffected makes sure the normal Graph flow is not
// changed by the web_auth additions.
func TestMakeOauthConfigDefaultUnaffected(t *testing.T) {
	ctx := context.Background()
	opt := &Options{
		Region:       regionGlobal,
		AccessScopes: scopeAccess,
	}
	conf, err := makeOauthConfig(ctx, opt)
	require.NoError(t, err)
	assert.Equal(t, rcloneClientID, conf.ClientID)
	assert.NotEmpty(t, conf.ClientSecret)
	assert.Equal(t, []string(scopeAccess), []string(conf.Scopes))
}
