package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	URL = "https://api.github.com/users/%s/events"
)

type GHActivityRepo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GHActicityActor struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type GHActivity struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	// GHActicityActor `json:"actor"`
	GHActivityRepo `json:"repo"`
	// Payload         interface{} //optional????
}

type Activity struct {
	action string
	times  int
}

func GitHubRequest(user string) []byte {
	u := fmt.Sprintf(URL, user)
	by := make([]byte, 0)
	c := http.DefaultClient
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Add("accept", "application/vnd.github+json")
	// resp, err := http.Get(u)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("Error calling the GitHub's Api: ", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("User %s not found. Please check the name.\n", user)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Fail on collection data from github. Code %s\n", resp.Status)
		os.Exit(1)
	}

	fmt.Println("StatusCode: ", resp.Status)
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		by = append(by, sc.Bytes()...)
	}
	defer resp.Body.Close()

	return by
}

/*
Parse the GHActivity to optimal format map
*/
func parseToMap(activities []GHActivity) map[string]map[string]int {
	data := map[string]map[string]int{}
	for _, actv := range activities {
		_, bol := data[actv.GHActivityRepo.Name]
		if !bol {
			data[actv.GHActivityRepo.Name] = map[string]int{}
		}
		data[actv.GHActivityRepo.Name][actv.Type] += 1
	}

	return data
}

/*
Pretty print the data based on the repository and action's type
*/
func prettyPrint(omap map[string]map[string]int) {
	ptrn := "   > %s, %d time(s) at repo %s.\n"
	fmt.Println("Activities:")
	for repo, data := range omap {
		for action, v := range data {
			switch action {
			case "CreateEvent":
				fmt.Printf("   > Create a branch from %s.\n", repo)
			case "PushEvent":
				fmt.Printf("   > Pushed %d commits to %s.\n", v, repo)
			case "WatchEvent":
				fmt.Printf("   > Watched %s.\n", repo)
			case "IssuesEvent":
				fmt.Printf("   > Opened a new issue in %s.\n", repo)
			case "ForkEvent":
				fmt.Printf("   > Forked the repo %s.\n", repo)
			default:
				fmt.Printf(ptrn, action, v, repo)
			}
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Expecting the github username. Example:\n   ./github-activity username")
		os.Exit(1)
	}
	by := GitHubRequest(args[0])
	if string(by) == "" {
		fmt.Printf("No recent activity found for the %s user's activities.\n", args[0])
		os.Exit(1)
	}
	var activities []GHActivity
	if err := json.Unmarshal(by, &activities); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	optMap := parseToMap(activities)
	prettyPrint(optMap)
}
