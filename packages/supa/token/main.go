package token

import (
	"context"

	"github.com/supabase-community/supabase-go"
)


type Token struct {
	Access_token string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Token_type string `json:"token_type"`
	Expiry int64 `json:"expires_at"`
}

func GetTokenUser(ctx context.Context, client *supabase.Client, provider string, user_id string) (Token, error) {

	if err := ctx.Err(); err != nil {
        return Token{}, err
    }

	var res Token

	_, err := client.From("accounts").Select("access_token,refresh_token,token_type,expires_at", "exact", false).Eq("userId", user_id).Eq("provider", provider).Single().ExecuteTo(&res)

	if err != nil {
		return Token{}, err
	}

	return res, nil

}