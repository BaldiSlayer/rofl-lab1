package githubclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v66/github"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	retryMax = 3
)

type Client struct {
	gh *github.Client
}

func New(ghToken string) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retryMax
	standardClient := retryClient.StandardClient() // *http.Client

	return &Client{
		gh: github.NewClient(standardClient).WithAuthToken(ghToken),
	}
}

func (c *Client) GistCreate(ctx context.Context, gist Gist) (string, error) {
	isPublic := false

	input := &github.Gist{
		Description: &gist.Description,
		Public:      &isPublic,
		Files: func() map[github.GistFilename]github.GistFile {
			res := make(map[github.GistFilename]github.GistFile)
			for i := range gist.Files {
				file := gist.Files[i]
				res[github.GistFilename(file.Name)] = github.GistFile{
					Size:     nil,
					Filename: nil,
					Language: nil,
					Type:     nil,
					RawURL:   nil,
					Content:  &file.Content,
				}
			}
			return res
		}(),
	}

	createdGist, resp, err := c.gh.Gists.Create(ctx, *&input)
	if err != nil {
		return "", fmt.Errorf("error requesting POST /gists: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("github POST /gists request returned non-201 code: %d", resp.StatusCode)
	}

	link := createdGist.GetHTMLURL()
	if link == "" {
		return "", fmt.Errorf("github POST /gists request returned empty link")
	}

	return link, nil
}
