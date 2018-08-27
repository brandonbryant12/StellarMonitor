package main

type Ledger struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Prev struct {
			Href string `json:"href"`
		} `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Record []struct {
			Id             string `json:"id"`
			Paging_Token   string `json:"paging_token"`
			Source_Account string `json:"source_account"`
			Type           string `json:"type"`
			CreatedAt      string `json:"created_at"`
			Hash           string `json:"transaction_hash"`
			Currency       string `json:"asset_type"`
			From           string `json:"from"`
			To             string `json:"to"`
			Amount         string `json:"amount"`
		} `json:"records"`
	} `json:"_embedded"`
}
