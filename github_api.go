package main

import (
	"fmt"
	"github.com/Khan/genqlient/graphql"
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

func NewGithubApi(token string) *GithubApi {
	httpClient := http.Client{
		Transport: &AuthedTransport{
			token:        token,
			roundTripper: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	return &GithubApi{
		client: &graphqlClient,
	}
}

type GithubApi struct {
	client *graphql.Client
	*Logger
}

func (r *GithubApi) UpdateClient(token string) {
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
