package update

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/github/sha"
)

type GitHubContentResponse struct {
	Content struct {
		SHA string `json:"sha"`
	} `json:"content"`
	Message string `json:"message"`
}


func UpdateFile(ctx context.Context, nameFile, content, accessToken, owner, repo string) error {

	if err := ctx.Err(); err != nil {
        return err
    }

	oldSha, err := sha.GetFileSHA(ctx, nameFile, accessToken, owner, repo)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		// message is like a commit
		"message": fmt.Sprintf("%s note %s backup at %s",
			func() string {
				if oldSha != "" {
					return "Update"
				}
				return "Add"
			}(),
			nameFile,
			time.Now().UTC().Format(time.RFC3339),
		),
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
	}

	if oldSha != "" {
		body["sha"] = oldSha
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx,
		"PUT",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, nameFile),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	time.Sleep(200*time.Millisecond)

	return nil
}
