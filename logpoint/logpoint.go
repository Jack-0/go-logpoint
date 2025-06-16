package logpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jack-0/go-logpoint/models"
)

type Logpoint struct {
	url      string
	username string
	secret   string
	client   *http.Client
}

func New(url string, username string, secret string) *Logpoint {

	// verify url
	// strip ending / to avoid issues

	return &Logpoint{
		url:      url,
		client:   &http.Client{},
		username: username,
		secret:   secret,
	}
}

func (logpoint *Logpoint) Query(query string, timeRange string, limit int, repos []string) {
	payload := map[string]interface{}{
		"query":      query,
		"time_range": timeRange,
		"limit":      limit,
		"repos":      repos,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	encodedPayloadStr := url.QueryEscape(string(jsonPayload))
	payloadStr := fmt.Sprintf("username=%s&secret_key=%s&requestData=%s", logpoint.username, logpoint.secret, encodedPayloadStr)
	data := strings.NewReader(payloadStr)

	fmt.Println(data)

	client := &http.Client{}
	route := logpoint.url + "/getsearchlogs"
	method := "POST"
	req, err := http.NewRequest(method, route, data)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(body)
	fmt.Println(string(body))
}

func (logpoint *Logpoint) GetRepos() models.RepoRequestResponse {

	// payload := map[string]interface{}{
	// 		"username":   logpoint.username,
	// 		"secret_key": logpoint.secret,
	// 		"type"
	// 	}

	method := "POST"
	payloadStr := fmt.Sprintf("username=%s&secret_key=%s&type=logpoint_repos", logpoint.username, logpoint.secret)
	payload := strings.NewReader(payloadStr)

	client := &http.Client{}
	route := logpoint.url + "/getalloweddata"
	req, err := http.NewRequest(method, route, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))
	var response models.RepoRequestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}
