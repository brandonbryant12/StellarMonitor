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
			Id            string `json:"id"`
			PagingToken   string `json:"paging_token"`
			SourceAccount string `json:"source_account"`
			Type          string `json:"type"`
			CreatedAt     string `json:"created_at"`
			Hash          string `json:"transaction_hash"`
			AssetType     string `json:"asset_type"`
			AssetCode     string `json:"asset_code"`
			AssetIssuer   string `json:"asset_issuer"`
			From          string `json:"from"`
			To            string `json:"to"`
			Amount        string `json:"amount"`
		} `json:"records"`
	} `json:"_embedded"`
}
