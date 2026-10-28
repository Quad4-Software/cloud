package oidc_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Quad4-Software/cloud/pkg/log"
	"github.com/Quad4-Software/cloud/pkg/oidc"
	"github.com/Quad4-Software/cloud/services/proxy/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestAccessTokenAudiences(t *testing.T) {
	key := newRSAKey(t)
	tests := []struct {
		name      string
		audiences []string
		aud       any
		missing   bool
		wantErr   error
	}{
		{name: "disabled missing", missing: true},
		{name: "disabled foreign", aud: "immich"},
		{name: "disabled null", aud: nil},
		{name: "disabled empty array", aud: []string{}},
		{name: "explicitly empty configuration", audiences: []string{}, aud: "immich"},
		{name: "single audience", audiences: []string{"quad4"}, aud: "quad4"},
		{name: "array first match", audiences: []string{"quad4"}, aud: []string{"quad4", "immich"}},
		{name: "array last match", audiences: []string{"quad4"}, aud: []string{"immich", "quad4"}},
		{name: "any allowed audience", audiences: []string{"quad4", "quad4-api"}, aud: "quad4-api"},
		{name: "any allowed audience in array", audiences: []string{"quad4", "quad4-api"}, aud: []string{"immich", "quad4-api"}},
		{name: "duplicate audiences", audiences: []string{"quad4", "quad4"}, aud: []string{"quad4", "quad4"}},
		{name: "URI audience", audiences: []string{"https://cloud.example/api"}, aud: "https://cloud.example/api"},
		{name: "foreign", audiences: []string{"quad4"}, aud: "immich", wantErr: jwt.ErrTokenInvalidAudience},
		{name: "foreign array", audiences: []string{"quad4"}, aud: []string{"immich", "account"}, wantErr: jwt.ErrTokenInvalidAudience},
		{name: "case sensitive", audiences: []string{"quad4"}, aud: "Quad4 Cloud", wantErr: jwt.ErrTokenInvalidAudience},
		{name: "exact match", audiences: []string{"quad4"}, aud: "quad4-api", wantErr: jwt.ErrTokenInvalidAudience},
		{name: "no token normalization", audiences: []string{"quad4"}, aud: " quad4 ", wantErr: jwt.ErrTokenInvalidAudience},
		{name: "no wildcard matching", audiences: []string{"*"}, aud: "quad4", wantErr: jwt.ErrTokenInvalidAudience},
		{name: "missing", audiences: []string{"quad4"}, missing: true, wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "null", audiences: []string{"quad4"}, aud: nil, wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "empty string", audiences: []string{"quad4"}, aud: "", wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "empty array", audiences: []string{"quad4"}, aud: []string{}, wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "array empty string", audiences: []string{"quad4"}, aud: []string{""}, wantErr: jwt.ErrTokenRequiredClaimMissing},
		{name: "number", audiences: []string{"quad4"}, aud: 123, wantErr: jwt.ErrTokenMalformed},
		{name: "object", audiences: []string{"quad4"}, aud: map[string]string{"aud": "quad4"}, wantErr: jwt.ErrTokenMalformed},
		{name: "mixed array", audiences: []string{"quad4"}, aud: []any{"quad4", 123}, wantErr: jwt.ErrTokenMalformed},
		{name: "null array entry", audiences: []string{"quad4"}, aud: []any{"quad4", nil}, wantErr: jwt.ErrTokenMalformed},
		{name: "disabled still rejects invalid type", aud: 123, wantErr: jwt.ErrTokenMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"iss": "https://issuer.example",
				"sub": "alice",
				"sid": "session",
				"exp": time.Now().Add(time.Hour).Unix(),
			}
			if !tt.missing {
				claims["aud"] = tt.aud
			}
			client := newAccessTokenTestClient(t, key, tt.audiences, &oidc.ProviderMetadata{})
			registered, all, err := client.VerifyAccessToken(context.Background(), signAccessToken(t, key, claims))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Empty(t, all, "unverified claims must not be returned")
				return
			}
			require.NoError(t, err)
			require.Equal(t, "alice", registered.Subject)
			require.Equal(t, "session", registered.SessionID)
			require.Equal(t, "alice", all["sub"])
		})
	}
}

func TestAccessTokenValidationWithAudiences(t *testing.T) {
	key, otherKey := newRSAKey(t), newRSAKey(t)
	tests := []struct {
		name       string
		issuer     string
		provider   *oidc.ProviderMetadata
		signingKey *signingKey
		exp        time.Time
		nbf        time.Time
		wantErr    error
	}{
		{name: "invalid signature", signingKey: otherKey, wantErr: jwt.ErrTokenSignatureInvalid},
		{name: "invalid issuer", issuer: "https://other.example", wantErr: jwt.ErrTokenInvalidIssuer},
		{name: "expired", exp: time.Now().Add(-time.Hour), wantErr: jwt.ErrTokenExpired},
		{name: "not yet valid", nbf: time.Now().Add(time.Hour), wantErr: jwt.ErrTokenNotValidYet},
		{name: "AD FS access token issuer", issuer: "https://adfs.example", provider: &oidc.ProviderMetadata{AccessTokenIssuer: "https://adfs.example"}},
		{name: "AD FS rejects discovery issuer", provider: &oidc.ProviderMetadata{AccessTokenIssuer: "https://adfs.example"}, wantErr: jwt.ErrTokenInvalidIssuer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.issuer == "" {
				tt.issuer = "https://issuer.example"
			}
			if tt.provider == nil {
				tt.provider = &oidc.ProviderMetadata{}
			}
			if tt.signingKey == nil {
				tt.signingKey = key
			}
			if tt.exp.IsZero() {
				tt.exp = time.Now().Add(time.Hour)
			}
			claims := jwt.MapClaims{"iss": tt.issuer, "sub": "alice", "aud": "quad4", "exp": tt.exp.Unix()}
			if !tt.nbf.IsZero() {
				claims["nbf"] = tt.nbf.Unix()
			}
			client := newAccessTokenTestClient(t, key, []string{"quad4"}, tt.provider)
			_, _, err := client.VerifyAccessToken(context.Background(), signAccessToken(t, tt.signingKey, claims))
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAccessTokenAudienceConfiguration(t *testing.T) {
	for _, method := range []string{config.AccessTokenVerificationNone, ""} {
		t.Run("incompatible method "+method, func(t *testing.T) {
			// No HTTP client is supplied: invalid configuration must fail before discovery.
			client, err := oidc.NewOIDCClient(
				oidc.WithAccessTokenVerifyMethod(method),
				oidc.WithAccessTokenAudiences([]string{"quad4"}),
			)
			require.ErrorContains(t, err, "requires the jwt verification method")
			require.Nil(t, client)
		})
	}
	for _, audiences := range [][]string{{""}, {" \t"}, {"quad4", ""}} {
		client, err := oidc.NewOIDCClient(
			oidc.WithAccessTokenVerifyMethod(config.AccessTokenVerificationJWT),
			oidc.WithAccessTokenAudiences(audiences),
		)
		require.ErrorContains(t, err, "empty or whitespace-only")
		require.Nil(t, client)
	}
	t.Run("none remains compatible when disabled", func(t *testing.T) {
		client, err := oidc.NewOIDCClient(
			oidc.WithLogger(log.NopLogger()),
			oidc.WithAccessTokenVerifyMethod(config.AccessTokenVerificationNone),
			oidc.WithProviderMetadata(&oidc.ProviderMetadata{}),
		)
		require.NoError(t, err)
		_, _, err = client.VerifyAccessToken(context.Background(), "opaque-token")
		require.NoError(t, err)
	})
	t.Run("caller cannot mutate the policy", func(t *testing.T) {
		key := newRSAKey(t)
		audiences := []string{"quad4"}
		client := newAccessTokenTestClient(t, key, audiences, &oidc.ProviderMetadata{})
		audiences[0] = "immich"
		_, _, err := client.VerifyAccessToken(context.Background(), signAccessToken(t, key,
			jwt.MapClaims{"iss": "https://issuer.example", "aud": "immich"}))
		require.ErrorIs(t, err, jwt.ErrTokenInvalidAudience)
	})
}

func TestAccessTokenAudiencesDoNotApplyToLogoutTokens(t *testing.T) {
	key := newRSAKey(t)
	client := newAccessTokenTestClient(t, key, []string{"quad4"}, &oidc.ProviderMetadata{})
	token := signAccessToken(t, key, jwt.MapClaims{
		"iss": "https://issuer.example",
		"sub": "alice",
		"aud": "web-client",
		"events": map[string]any{
			"http://schemas.openid.net/event/backchannel-logout": map[string]any{},
		},
	})
	_, err := client.VerifyLogoutToken(context.Background(), token)
	require.NoError(t, err)
}

func newAccessTokenTestClient(t *testing.T, key *signingKey, audiences []string, provider *oidc.ProviderMetadata) oidc.OIDCClient {
	t.Helper()
	client, err := oidc.NewOIDCClient(
		oidc.WithLogger(log.NopLogger()),
		oidc.WithOidcIssuer("https://issuer.example"),
		oidc.WithAccessTokenVerifyMethod(config.AccessTokenVerificationJWT),
		oidc.WithAccessTokenAudiences(audiences),
		oidc.WithJWKS(key.jwks),
		oidc.WithProviderMetadata(provider),
	)
	require.NoError(t, err)
	return client
}

func signAccessToken(t *testing.T, key *signingKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "1"
	signed, err := token.SignedString(key.priv)
	require.NoError(t, err)
	return signed
}
