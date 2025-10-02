package models

//推し1人取得レスポンス
type GetOshiResponse struct {
	Oshi GetOshiResponseItem `json:"oshi"`
}

type GetOshiResponseItem struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Color      string   `json:"color"`
	URLs       []string `json:"urls"`
	Categories []string `json:"categories"`
}
