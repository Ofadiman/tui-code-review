package main

import (
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/ofadiman/tui-code-review/log"
	"net/http"
)

type AuthedTransport struct {
	token        string
	roundTripper http.RoundTripper
}

func (r *AuthedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+r.token)
	return r.roundTripper.RoundTrip(req)
}

func NewGithubApi(token string) *GitHubApi {
	httpClient := http.Client{
		Transport: &AuthedTransport{
			token:        token,
			roundTripper: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	return &GitHubApi{
		client: &graphqlClient,
	}
}

type GitHubApi struct {
	client *graphql.Client
	*log.Logger
}

func (r *GitHubApi) WithLogger(logger *log.Logger) *GitHubApi {
	r.Logger = logger

	return r
}

func (r *GitHubApi) UpdateClient(token string) {
	r.Logger.Info(fmt.Sprintf("creating a new graphql client with token %v", token))

	httpClient := http.Client{
		Transport: &AuthedTransport{
			token:        token,
			roundTripper: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)
	r.client = &graphqlClient
}
