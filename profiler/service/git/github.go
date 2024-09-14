package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xhermitx/gitpulse-01/profiler/configs"
	"github.com/xhermitx/gitpulse-01/profiler/types"
)

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

const (
	GH_URL = "https://api.github.com/graphql"
	query  = `query getUserDetails($github_id: String!) {
				user(login: $github_id) {
					# Fetch basic user info
					name
					login
					__typename
					avatarUrl
					bio
					email
					websiteUrl
					followers {
						totalCount
					}

					# Fetch contributions data
					contributionsCollection {
						contributionCalendar {
								totalContributions
							}
						}
					}

					# Fetch repositories created by the user and sort by stars
					repositories(first: 1, orderBy: { field: STARGAZERS, direction: DESC }, privacy: PUBLIC) {
						nodes {
							name
							description
							stargazers {
								totalCount
							}
							languages(first: 5) {
								nodes {
									name
								}
							}
							repositoryTopics(first: 5) {
								nodes {
									topic {
										name
									}
								}
							}
							url
						}
					}

					# Fetch repositories the user has contributed to (not their own) and sort by stars
					repositoriesContributedTo(first: 1, orderBy: { field: STARGAZERS, direction: DESC }, privacy: PUBLIC) {
						nodes {
							name
							owner {
								login
							}
							description
							stargazers {
								totalCount
							}
							languages(first: 5) {
								nodes {
									name
								}
							}
							repositoryTopics(first: 5) {
								nodes {
									topic {
										name
									}
								}
							}
							url
						}
					}
				}
			}`
)

func (g *GitService) FetchUserDetails(github_id string) (*types.GitUser, error) {
	ghQuery := types.GitQuery{
		Query: query,
		Variables: map[string]string{
			"github_id": github_id,
		},
	}

	body, err := json.Marshal(ghQuery)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, GH_URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.Envs.GithubToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ghResp types.GitResponse
	err = json.Unmarshal(responseBody, &ghResp)
	if err != nil {
		return nil, err
	}

	if len(ghResp.Errors) > 0 {
		return nil, fmt.Errorf("error occured while fetching username: %v", ghResp.Errors)
	}
	if ghResp.Data.User.AccountType != "User" {
		return nil, fmt.Errorf("username %s is not of type user", ghResp.Data.User.Username)
	}

	return &ghResp.Data.User, nil
}
