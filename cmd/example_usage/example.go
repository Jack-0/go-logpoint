package main

import (
	"fmt"

	"github.com/jack-0/go-logpoint/logpoint"
)

func main() {
	logpointURL := "https://<required>"
	logpointUsername := "<required>"
	logpointUserSecret := "<required>"

	if logpointUserSecret == "<required>" || logpointURL == "https://<required>" || logpointUsername == "<required>" {
		panic("Ensure you alter all required variables")
	}

	lp := logpoint.New(logpointURL, logpointUsername, logpointUserSecret, true)
	a, err := lp.GetRepos()
	if err != nil {
		panic(err)
	}

	fmt.Println("📃 Allowed Repos for given user:")
	fmt.Println(a.AllowedRepos)

	allRepos := []string{}
	for _, repo := range a.AllowedRepos {
		allRepos = append(allRepos, repo.Repo)
	}

	fmt.Println("🔍 Query 100 Logs from those repos in the last hour:")
	b, err := lp.Query("", "Last 1 hour", 100, allRepos, 10)
	if err != nil {
		panic(err)
	}

	qr, err := lp.QueryResult(b.SearchId)

	if err != nil {
		panic(err)
	}

	//// Uncomment the following to see results
	fmt.Println("...Uncomment example code to see full results...")
	// fmt.Println("🎁 Query Results:")
	// for _, item := range qr {
	// 	fmt.Println(item)
	// }

	fmt.Println("🏁 Total items:")
	fmt.Println(len(qr.Rows))
	fmt.Println("🏁 Total Number Aggregated:")
	fmt.Println(qr.Meta.NumAggregated)
}
