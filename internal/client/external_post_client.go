package client

import (
	"encoding/json"
	"fmt"
	"io"
	"lovender_backend/internal/models"
	"net/http"
	"sort"
	"time"
)

// 外部投稿API用のクライアント
type ExternalPostClient struct {
	baseURL    string
	httpClient *http.Client
}

// コンストラクタ
func NewExternalPostClient() *ExternalPostClient {
	return &ExternalPostClient{
		baseURL: "http://176.34.25.68:8000",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// アカウント名で投稿を取得
func (c *ExternalPostClient) GetPostsByUsername(accountName string) ([]models.ExternalPost, error) {
	url := fmt.Sprintf("%s/v1/posts/username/%s", c.baseURL, accountName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts for %s: %w", accountName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for account %s", resp.StatusCode, accountName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.ExternalPostsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Posts, nil
}

// 最新の投稿を指定件数取得
func (c *ExternalPostClient) GetLatestPostsByUsername(accountName string, limit int) ([]models.ExternalPost, error) {
	posts, err := c.GetPostsByUsername(accountName)
	if err != nil {
		return nil, err
	}

	// idで降順ソート
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	// 指定件数まで切り取り
	if len(posts) > limit {
		posts = posts[:limit]
	}

	return posts, nil
}
