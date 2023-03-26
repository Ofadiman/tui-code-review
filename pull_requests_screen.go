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

type PullRequestsScreen struct {
	*Window
	*Settings
	*Logger
	*GithubApi
	pullRequests []*PullRequest
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

// TODO: handle case of commented pull request.
const (
	PULL_REQUEST_AWAITING = "PULL_REQUEST_AWAITING"
	PULL_REQUEST_REJECTED = "PULL_REQUEST_REJECTED"
	PULL_REQUEST_APPROVED = "PULL_REQUEST_APPROVED"
	PULL_REQUEST_DRAFT    = "PULL_REQUEST_DRAFT"
)

func mapGithubPullRequestsToApplicationPullRequests(githubPullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest, user string) []*PullRequest {
	var applicationPullRequests []*PullRequest
	for _, githubPullRequest := range githubPullRequests {
		pullRequest := PullRequest{
			getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest: githubPullRequest,
		}

		if githubPullRequest.GetIsDraft() == true {
			pullRequest.state = PULL_REQUEST_DRAFT
		} else {
			for _, latestReviews := range githubPullRequest.GetLatestReviews().GetNodes() {
				if latestReviews.GetAuthor().GetLogin() == user {
					if latestReviews.GetState() == PullRequestReviewStateApproved {
						pullRequest.state = PULL_REQUEST_APPROVED
					} else if latestReviews.GetState() == PullRequestReviewStateChangesRequested {
						pullRequest.state = PULL_REQUEST_REJECTED
					}
				}
			}

			for _, reviewRequest := range githubPullRequest.GetReviewRequests().GetNodes() {
				requestedReviewer, ok := reviewRequest.GetRequestedReviewer().(*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequestReviewRequestsReviewRequestConnectionNodesReviewRequestRequestedReviewerUser)
				if ok {
					if requestedReviewer.GetLogin() == user {
						pullRequest.state = PULL_REQUEST_AWAITING
					}
				}
			}
		}

		applicationPullRequests = append(applicationPullRequests, &pullRequest)
	}

	return applicationPullRequests
}

func getGithubPullRequestsFromRepositories(repositoryInfoResponses []*getRepositoryInfoResponse) []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest {
	var pullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest

	for _, repositoryInfoResponse := range repositoryInfoResponses {
		pullRequests = append(pullRequests, repositoryInfoResponse.GetRepository().GetPullRequests().GetNodes()...)
	}

	return pullRequests
}

func sortPullRequestsForMe(pullRequestsForMe []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest, logger *Logger, username string) {
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
					if requestedReviewer.GetLogin() == username {
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
					if requestedReviewer.GetLogin() == username {
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
				if node.GetAuthor().GetLogin() == username {
					if node.GetState() == PullRequestReviewStateApproved {
						isFirstApproved = true
					}

					if node.GetState() == PullRequestReviewStateChangesRequested {
						isFirstRejected = true
					}
				}
			}

			for _, node := range pullRequestsForMe[j].GetLatestReviews().GetNodes() {
				if node.GetAuthor().GetLogin() == username {
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
	pullRequestsForMe := findPullRequestsForMe(allPullRequestsFromWatchedRepositories, r.Settings.Username)

	sortPullRequestsForMe(pullRequestsForMe, r.Logger, r.Settings.Username)

	r.Logger.Struct(pullRequestsForMe)

	r.pullRequests = mapGithubPullRequestsToApplicationPullRequests(pullRequestsForMe, r.Settings.Username)

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
	header := StyledHeader.Render("Pull requests")

	pullRequestStateToUI := map[string]string{
		PULL_REQUEST_AWAITING: StyledAwaiting.Render("review required"),
		PULL_REQUEST_APPROVED: StyledApproved.Render("approved"),
		PULL_REQUEST_REJECTED: StyledChangesRequested.Render("changes requested"),
		PULL_REQUEST_DRAFT:    StyledDraft.Render("draft"),
	}
	var pullRequestMessage string
	for _, pullRequest := range r.pullRequests {
		info, ok := pullRequestStateToUI[pullRequest.state]
		if !ok {
			r.Logger.Info(fmt.Sprintf("info does not exist for pull request state %v", pullRequest.state))
		}

		pullRequestMessage += fmt.Sprintf("• %v wants to merge \"%v\" (%v)\n", pullRequest.GetAuthor().GetLogin(), pullRequest.GetTitle(), info)
	}
	f := wordwrap.NewWriter(r.Window.Width - StyledMain.GetHorizontalPadding())
	f.Breakpoints = []rune{' '}
	_, err := f.Write([]byte(pullRequestMessage))
	if err != nil {
		r.Logger.Error(err)
	}

	return StyledMain.Render(lipgloss.JoinVertical(lipgloss.Left, header, f.String()))
}
