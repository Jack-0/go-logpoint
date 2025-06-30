package models

type SearchRequestResponse struct {
	// The total number of logs aggregated in the result set
	NumAggregated uint64 `json:"num_aggregated"`
	// Lists the columns returned by the query, such as count().
	Columns []string `json:"columns"`
	// Defines the type of query executed (Example: chart).
	QueryType string `json:"query_type"`
	// Contains the actual search result logs. Each object represents a log entry with relevant fields (e.g., device_ip).
	Rows []any `json:"rows"`
	// Defines the fields used to group the search results (e.g., device_ip).
	Grouping []string `json:"grouping"`
	// Version of the search result format or API being used.
	Version uint64 `json:"version"`
	// Lists any fields deemed interesting in the result set.
	InterestingFields []string `json:"interesting_fields"`
	// Contains two timestamps that define the start and end of the search time range.
	TimeRange []uint64 `json:"time_range"`
	// The original search_id used to initiate the search.
	OriginalSearchId string `json:"orig_search_id"`
	// Returns True if all search result logs are retrieved; otherwise it returns False.
	Final bool `json:"finial"`
	// Returns True if the API call is successful; otherwise it returns False.
	Success bool `json:"success"`
	// The total number of pages of results.
	TotalPages int64 `json:"totalPages"`
	// Returns True the search is successful; otherwise it returns False.
	Complete bool `json:"complete"`
	// Returns True additional visualizations should be displayed with the search result; otherwise it returns False.
	ShowAdditionalPanels bool `json:"showAdditionalPanels"`
	// Additional status information about the search, such as progress or execution details.
	Status any `json:"status"`
	// Error message?
	Message string `json:"message"`
}
