package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"sort"
	"strings"
)

var roundedBorder = lipgloss.RoundedBorder()
var columnStyle = lipgloss.NewStyle().Border(roundedBorder).BorderForeground(lipgloss.Color("63"))

type PullRequestsScreen struct {
	*Window
	*Settings
	*Logger
	*GithubApi
	pullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest
}

type PullRequest struct {
	*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest
	state string
}

func NewPullRequestsScreen(globalState *Window, settings *Settings, logger *Logger, githubApi *GithubApi) *PullRequestsScreen {
	return &PullRequestsScreen{
		Window:    globalState,
		Settings:  settings,
		Logger:    logger,
		GithubApi: githubApi,
	}
}

func getGithubPullRequestsFromRepositories(repositoryInfoResponses []*getRepositoryInfoResponse) []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest {
	var pullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest

	for _, repositoryInfoResponse := range repositoryInfoResponses {
		pullRequests = append(pullRequests, repositoryInfoResponse.GetRepository().GetPullRequests().GetNodes()...)
	}

	return pullRequests
}

func sortPullRequestsForMe(pullRequestsForMe []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest, logger *Logger) {
	sort.Slice(pullRequestsForMe, func(i, j int) bool {
		if pullRequestsForMe[i].GetIsDraft() == true && pullRequestsForMe[j].GetIsDraft() == false {
			return false
		}

		if pullRequestsForMe[i].GetIsDraft() == false && pullRequestsForMe[j].GetIsDraft() == true {
			return true
		}

		if pullRequestsForMe[i].GetIsDraft() == pullRequestsForMe[j].GetIsDraft() {
			isFirstAwaiting := false
			isSecondAwaiting := false

			for _, node := range pullRequestsForMe[i].GetReviewRequests().GetNodes() {
				requestedReviewer, ok := node.GetRequestedReviewer().(*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequestReviewRequestsReviewRequestConnectionNodesReviewRequestRequestedReviewerUser)
				if ok {
					if requestedReviewer.GetLogin() == "GuilermeheGardosso" {
						isFirstAwaiting = true
					}
				} else {
					logger.Info("requested reviewer is not it the shape I expect")
				}
			}

			for _, node := range pullRequestsForMe[j].GetReviewRequests().GetNodes() {
				logger.Struct(node)

				requestedReviewer, ok := node.GetRequestedReviewer().(*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequestReviewRequestsReviewRequestConnectionNodesReviewRequestRequestedReviewerUser)
				if ok {
					if requestedReviewer.GetLogin() == "GuilermeheGardosso" {
						isSecondAwaiting = true
					}
				} else {
					logger.Info("requested reviewer is not it the shape I expect")
				}
			}

			isFirstRejected := false
			isSecondRejected := false
			isFirstApproved := false
			isSecondApproved := false

			for _, node := range pullRequestsForMe[i].GetLatestReviews().GetNodes() {
				if node.GetAuthor().GetLogin() == "GuilermeheGardosso" {
					if node.GetState() == PullRequestReviewStateApproved {
						isFirstApproved = true
					}

					if node.GetState() == PullRequestReviewStateChangesRequested {
						isFirstRejected = true
					}
				}
			}

			for _, node := range pullRequestsForMe[j].GetLatestReviews().GetNodes() {
				if node.GetAuthor().GetLogin() == "GuilermeheGardosso" {
					if node.GetState() == PullRequestReviewStateApproved {
						isSecondApproved = true
					}

					if node.GetState() == PullRequestReviewStateChangesRequested {
						isSecondRejected = true
					}
				}
			}

			if isFirstRejected && isSecondAwaiting {
				return false
			}

			if isFirstApproved && (isSecondAwaiting || isSecondRejected) {
				return false
			}

			if isFirstApproved && isSecondApproved || isFirstAwaiting && isSecondAwaiting || isFirstRejected && isSecondRejected {
				return pullRequestsForMe[i].GetCreatedAt().After(pullRequestsForMe[j].GetCreatedAt())
			}

			return true
		}

		return false
	})
}

func findPullRequestsForMe(pullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest, user string) []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest {
	var final []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest
	for _, pullRequest := range pullRequests {
		isSubmittedByMe := false
		isRequestingMyReview := false
		isAlreadyReviewedByMe := false

		if pullRequest.GetAuthor().GetLogin() == user {
			isSubmittedByMe = true
		}

		for _, reviewRequest := range pullRequest.GetReviewRequests().GetNodes() {
			requestedReviewer, ok := reviewRequest.GetRequestedReviewer().(*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequestReviewRequestsReviewRequestConnectionNodesReviewRequestRequestedReviewerUser)
			if ok {
				if requestedReviewer.GetLogin() == user {
					isRequestingMyReview = true
				}
			}
		}

		for _, x := range pullRequest.GetLatestReviews().GetNodes() {
			if x.GetAuthor().GetLogin() == user {
				isAlreadyReviewedByMe = true
			}
		}

		if isSubmittedByMe == false && (isRequestingMyReview || isAlreadyReviewedByMe) {
			final = append(final, pullRequest)
		}
	}

	return final
}

func (r *PullRequestsScreen) Init() tea.Cmd {
	if r.Settings.GithubToken == "" {
		return nil
	}

	channel := make(chan *getRepositoryInfoResponse)
	responses := make([]*getRepositoryInfoResponse, len(r.Settings.Repositories))

	for _, repositoryUrl := range r.Settings.Repositories {
		go func(repositoryUrl string) {
			urlParts := strings.Split(repositoryUrl, "/")
			username := urlParts[len(urlParts)-2]
			repositoryName := urlParts[len(urlParts)-1]

			r.Logger.Info(fmt.Sprintf("sending request to %v/%v", username, repositoryName))

			var response *getRepositoryInfoResponse
			var err error
			response, err = getRepositoryInfo(context.Background(), *r.GithubApi.client, username, repositoryName)
			if err != nil {
				r.Logger.Info("error is not nil")
				channel <- nil

				r.Logger.Error(err)

				if strings.Contains(err.Error(), "401") {
					r.Settings.UpdateGitHubToken("")
				}
			} else {
				r.Logger.Info("passing response to channel")
				channel <- response
			}
		}(repositoryUrl)
	}

	for i := 0; i < len(r.Settings.Repositories); i++ {
		responses[i] = <-channel
	}

	allPullRequestsFromWatchedRepositories := getGithubPullRequestsFromRepositories(responses)

	// TODO: Take user from settings.
	pullRequestsForMe := findPullRequestsForMe(allPullRequestsFromWatchedRepositories, "GuilermeheGardosso")

	sortPullRequestsForMe(pullRequestsForMe, r.Logger)

	r.Logger.Struct(pullRequestsForMe)

	r.pullRequests = pullRequestsForMe

	return nil
}

func (r *PullRequestsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "s":
				{
					r.Logger.KeyPress("s")
					return r, nil
				}
			}
		}
	}

	return r, nil
}

func (r *PullRequestsScreen) View() string {
	columnStyle.Width(r.Window.Width - roundedBorder.GetLeftSize() - roundedBorder.GetRightSize())
	header := columnStyle.Render("Pull requests")

	var todo string
	for _, value := range r.pullRequests {
		todo += fmt.Sprintf("isDraft: %v\n", value.GetIsDraft())
		todo += fmt.Sprintf("author: %v\n", value.GetAuthor().GetLogin())
		todo += "\n"
	}
	f := wordwrap.NewWriter(r.Window.Width - roundedBorder.GetLeftSize() - roundedBorder.GetRightSize())
	f.Breakpoints = []rune{' '}
	_, err := f.Write([]byte(todo))
	if err != nil {
		r.Logger.Error(err)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, columnStyle.Render(f.String()))
}
