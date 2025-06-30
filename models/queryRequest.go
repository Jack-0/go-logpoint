package models

type QueryRequestResponse struct {
	Success     bool     `json:"success"`
	SearchId    string   `json:"searchId"`
	QueryFilter string   `json:"query_filter"`
	Latest      bool     `json:"latest"`
	Lookup      bool     `json:"lookup"`
	QueryType   string   `json:"query_type"`
	ClientType  string   `json:"client_type"`
	TimeRange   []uint64 `json:"time_range"`
}
