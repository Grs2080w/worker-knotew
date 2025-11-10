package github

import (
	"context"

	"github.com/Grs2080w/worker-knotew/packages/github/update"
)

type gitHub struct {
	owner string
	repo string
	access_token string
}


func New(accessToken, owner, repo string) *gitHub {
	return &gitHub{
		owner: owner,
		repo: repo,
		access_token: accessToken,
	}
}


func (g *gitHub) UpdateFile(ctx context.Context, nameFile, content string) error {
	return update.UpdateFile(ctx, nameFile, content, g.access_token, g.owner, g.repo)
}