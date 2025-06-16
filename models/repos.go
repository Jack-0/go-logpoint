package models

type RepoRequestResponse struct {
	Success      bool           `json:"success"`
	AllowedRepos []RepoItem     `json:"allowed_repos"`
	Logpoint     []LogpointItem `json:"logpoint"`
}

type RepoItem struct {
	Repo    string `json:"repo"`
	Address string `json:"address"`
}
type LogpointItem struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}
