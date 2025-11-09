package google

import (
	"context"

	"github.com/Grs2080w/worker-knoteq/packages/google/client"
	"github.com/Grs2080w/worker-knoteq/packages/google/update"
	"github.com/Grs2080w/worker-knoteq/packages/supa/token"
	"google.golang.org/api/drive/v3"
)

type Google struct {
	Srv *drive.Service
}

// New: Create new Google client
func New(ctx context.Context, tok token.Token) (*Google, error) {

	srv, err := client.New(ctx, tok)
	if err != nil {
		return nil, err
	}

	return &Google{Srv: srv}, nil
}

// UpdateFile: Update file on Google Drive receive id and content
func (google *Google) UpdateFile(ctx context.Context, id string, content []byte) (error) {
	return update.UpdateGoogleFile(ctx, id, content, google.Srv)
}


