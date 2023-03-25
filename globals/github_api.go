package globals

import (
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"net/http"
)

type authedTransport struct {
	token        string
	roundTripper http.RoundTripper
}

func (r *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+r.token)
	return r.roundTripper.RoundTrip(req)
}

func newGithubApi(token string) *githubApi {
	httpClient := http.Client{
		Transport: &authedTransport{
			token:        token,
			roundTripper: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)

	return &githubApi{
		Client: &graphqlClient,
	}
}

type githubApi struct {
	Client *graphql.Client
	*logger
}

func (r *githubApi) UpdateClient(token string) {
	r.logger.Info(fmt.Sprintf("creating a new graphql Client with token %v", token))

	httpClient := http.Client{
		Transport: &authedTransport{
			token:        token,
			roundTripper: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)
	r.Client = &graphqlClient
}
