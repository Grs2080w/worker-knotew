package client

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/supa/token"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func New(ctx context.Context, tok token.Token) (*drive.Service, error) {

	if err := ctx.Err(); err != nil {
        return nil, err
    }
	
	client, err := getClient(ctx, tok)
	if err != nil {
		return nil, err
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}

	return srv, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, tok token.Token) (*http.Client, error) {

	if err := ctx.Err(); err != nil {
        return nil, err
    }

	godotenv.Load()

	// This credentials must be the same as the used by the main app, including the redirect
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirect := os.Getenv("GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		log.Fatalf("Missing GOOGLE_CLIENT_ID and/or GOOGLE_CLIENT_SECRET")
	}

	config := &oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		RedirectURL: redirect,
		Scopes: []string{drive.DriveFileScope},
		Endpoint: google.Endpoint,
	}

	return config.Client(context.Background(), &oauth2.Token{
		AccessToken:  tok.Access_token,
		RefreshToken: tok.Refresh_token,
		TokenType:    tok.Token_type,
		Expiry:       time.Unix(tok.Expiry, 0),
	}), nil
}
