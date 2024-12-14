package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	githubURL, ok := os.LookupEnv("GITHUB_SERVER_URL")
	if !ok {
		fmt.Println("GitHub server URL is not provided")
		os.Exit(1)
	}

	repository, ok := os.LookupEnv("GITHUB_REPOSITORY")
	if !ok {
		fmt.Println("Repository is not provided")
		os.Exit(1)
	}

	prNumber, ok := os.LookupEnv("INPUT_PR_NUMBER")
	if !ok {
		fmt.Println("PR number is not provided")
		os.Exit(1)
	}

	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("GitHub token is not provided")
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("%s/repos/%s/pulls/%s", githubURL, repository, prNumber)

	// Create a new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received non-200 response status: %s\n", resp.Status)
		os.Exit(1)
	}

	type PullRequest struct {
		Title string `json:"title"`
	}

	// Decode the response
	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	// Construct the PR link
	prLink := fmt.Sprintf("[%s](%s/%s/pull/%s)", pr.Title, githubURL, repository, prNumber)

	// Print the PR link
	fmt.Printf("::notice file=main.go,line=14::%s\n", prLink)

	// Write outputs to the $GITHUB_OUTPUT file
	outputFile, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(fmt.Sprintf("pr_link=%s\n", prLink)); err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		os.Exit(1)
	}
}
