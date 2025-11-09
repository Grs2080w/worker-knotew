package update

import (
	"bytes"
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
)

func UpdateGoogleFile(ctx context.Context, id string, content []byte, srv *drive.Service) error {

	if err := ctx.Err(); err != nil {
        return err
    }

	reader := bytes.NewReader(content)

	_, err := srv.Files.Update(id, &drive.File{}).Media(reader).Context(ctx).Do()

	if err != nil {
		return fmt.Errorf("error updating google drive file: %w", err)
	}

	return nil;
}