package logpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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

// Helper function to call the getSearchLog endpoint with a json payload (API uses x-www-form-urlencoded)
func getSearchLogs[T any](logpoint *Logpoint, requestData map[string]interface{}) (*T, error) {
	fmt.Printf("\nRequest data = %+v\n", requestData)

	jsonPayload, err := json.Marshal(requestData)
	if err != nil {
		return nil, error(err)
	}
	encodedPayloadStr := url.QueryEscape(string(jsonPayload))
	payloadStr := fmt.Sprintf("username=%s&secret_key=%s&requestData=%s", logpoint.username, logpoint.secret, encodedPayloadStr)
	data := strings.NewReader(payloadStr)

	client := &http.Client{}
	route := logpoint.url + "/getsearchlogs"
	method := "POST"
	req, err := http.NewRequest(method, route, data)

	if err != nil {
		return nil, error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, error(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, error(err)
	}

	var response T
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, error(err)
	}
	return &response, nil
}

func (logpoint *Logpoint) Query(query string, timeRange string, limit int, repos []string, timeoutSeconds int) (*models.QueryRequestResponse, error) {
	// handle timeout
	if timeoutSeconds < 5 || timeoutSeconds > 90 {
		fmt.Println("Timeout could be problematic try a range within 5-90 seconds")
	}

	requestData := map[string]interface{}{
		"timeout":    timeoutSeconds,
		"query":      query,
		"time_range": timeRange,
		"limit":      limit, // TODO: review limit usage
		"repos":      repos,
	}
	return getSearchLogs[models.QueryRequestResponse](logpoint, requestData)
}

func (logpoint *Logpoint) QueryResult(searchId string) ([]interface{}, error) {
	// From the api docs:
	// Retrieve search result logs based on the search_id. The server sends the search result logs in chunks. You need to continue sending the request with the same parameters until you receive a response where final is equal to TRUE. It indicates that you have received all the search result logs.

	payload := map[string]interface{}{
		"searchId": searchId,
	}

	time.Sleep(200 * time.Millisecond) // add slight delay for Logpoint
	// time.Sleep( * time.Second) // add slight delay for Logpoint

	// TODO: compair logpoint response with this query with and without the delay.
	// Then use whatever var says search complete to wait until that is done before pulling the data
	// profit
	res, err := getSearchLogs[models.SearchRequestResponse](logpoint, payload)
	if err != nil {
		return nil, error(err)
	}
	fmt.Printf("RES: complete=%t, totalPages=%d, num_aggregated=%d", res.Complete, res.TotalPages, res.NumAggregated)

	if !res.Success {
		return nil, fmt.Errorf("%s", res.Message)
	}

	rows := []interface{}{}
	finished := res.Final

	// Logpoints docs are unclear and pagination process is poor. Here we assume if we get 0 total pages then
	// we have the data we request event though res.Final could be false...
	if res.TotalPages == 0 {
		finished = true
	}
	// payload := map[string]interface{}{
	// 	"searchId": res,
	// }

	// Recursively fetch data
	attempts := 0
	for !finished {
		if attempts > 50 {
			return nil, fmt.Errorf("Pagination failed; logpoints docs are a bit naff. Possible ticket required.")
		}
		rows = append(rows, res.Rows...)
		res, err = getSearchLogs[models.SearchRequestResponse](logpoint, payload)
		if err != nil {
			return nil, error(err)
		}
		finished = res.Final
		// time.Sleep(60 * time.Second) // add slight delay for Logpoint
		time.Sleep(100 * time.Millisecond) // add slight delay for Logpoint
		attempts++
	}

	rows = append(rows, res.Rows...)

	return rows, nil
}

func (logpoint *Logpoint) GetRepos() (*models.RepoRequestResponse, error) {
	method := "POST"
	payloadStr := fmt.Sprintf("username=%s&secret_key=%s&type=logpoint_repos", logpoint.username, logpoint.secret)
	payload := strings.NewReader(payloadStr)

	client := &http.Client{}
	route := logpoint.url + "/getalloweddata"
	req, err := http.NewRequest(method, route, payload)

	if err != nil {
		return nil, error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, error(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, error(err)
	}

	// fmt.Println(string(body))
	var response models.RepoRequestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, error(err)
	}

	if !response.Success {
		fmt.Errorf("Failed to query repos")
	}

	return &response, nil
}
