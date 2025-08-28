package config

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Keycloak struct {
	ClientID     string `env:"KEYCLOCK_CLIENT_ID" envDefault:"wongnok"`
	ClientSecret string `env:"KEYCLOCK_CLIENT_SECRET" envDefault:"R6UCkSremn7nOYkNzqxJVUcVNPnG5fu7"`
	RedirectURL  string `env:"KEYCLOCK_REDIRECT_URL" envDefault:"http://localhost:8000/api/v1/callback"`
	Realm        string `env:"KEYCLOAK_REALM" envDefault:"pea-devpool-2025"`
	URL          string `env:"KEYCLOAK_URL" envDefault:"https://sso-dev.odd.works"`
}

func (kc Keycloak) RealmURL() string {
	return fmt.Sprintf("%s/realms/%s", kc.URL, kc.Realm)
}

func (kc Keycloak) LogoutURL() string {
	return fmt.Sprintf("%s/protocol/openid-connect/logout", kc.RealmURL())
}

type IOAuth2Config interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type IOIDCTokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

type IOIDCIDToken interface {
	Claims(v any) error
}
