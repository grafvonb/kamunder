package testx

import (
	"net/http"

	config2 "github.com/grafvonb/kamunder/config"
	"github.com/grafvonb/kamunder/internal/clients/auth/oauth2"
)

type tokenJSON200 = struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int     `json:"expires_in"`
	IdToken      *string `json:"id_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	Scope        *string `json:"scope,omitempty"`
	TokenType    string  `json:"token_type"`
}

func TestAuthJSON200Response(status int, token string, raw string) *oauth2.RequestTokenResponse {
	return &oauth2.RequestTokenResponse{
		Body: []byte(raw),
		JSON200: &tokenJSON200{
			AccessToken: token,
			TokenType:   "Bearer",
		},
		HTTPResponse: &http.Response{StatusCode: status},
	}
}

func TestConfig() *config2.Config {
	return &config2.Config{
		App: config2.App{
			Tenant: "tenant",
		},
		Auth: config2.Auth{
			OAuth2: config2.AuthOAuth2ClientCredentials{
				TokenURL:     "http://localhost/token",
				ClientID:     "test",
				ClientSecret: "test",
			},
			Cookie: config2.AuthCookieSession{
				BaseURL:  "http://localhost/cookie",
				Username: "test",
				Password: "test",
			},
		},
		APIs: config2.APIs{
			Camunda: config2.API{
				BaseURL: "http://localhost/camunda/v2",
			},
			Operate: config2.API{
				BaseURL: "http://localhost/operate",
			},
			Tasklist: config2.API{
				BaseURL: "http://localhost/tasklist",
			},
		},
	}
}
