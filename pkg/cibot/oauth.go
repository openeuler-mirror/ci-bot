package cibot

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func GetToken(code, client, lang string) (*oauth2.Token, error) {

	url := os.Getenv("WEBSITE_URL")
	if url == "" {
		url = "https://openeuler.org"
	}

	secret := os.Getenv("GITEE_OAUTH2_SECRET_EN")

	if lang == "zh" {
		secret = os.Getenv("GITEE_OAUTH2_SECRET_ZH")
	}

	if !strings.HasSuffix(url, "/") {
		url = fmt.Sprintf("%s/", url)
	}

	redirect := fmt.Sprintf("%s%s/cla.html", strings.Trim(url, " "), strings.Trim(lang, " "))

	ctx := context.Background()
	config := Setup(client, redirect, secret)

	return config.Exchange(ctx, code)

}

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://gitee.com/oauth/authorize",
	TokenURL: "https://gitee.com/oauth/token",
}

func Setup(client, redirect, secret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     client,
		ClientSecret: secret,
		Scopes:       []string{"emails", "user_info"},
		Endpoint:     Endpoint,
		RedirectURL:  redirect,
	}
}
