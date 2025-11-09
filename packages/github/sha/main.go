package sha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetFileSHA(ctx context.Context, nameFile, accessToken, owner, repo string) (string, error) {

	if err := ctx.Err(); err != nil {
        return "", err
    }

	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, nameFile),
		nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var data struct {
			SHA string `json:"sha"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}
		return data.SHA, nil
	}

	return "", nil
}
