package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

var PULL_REQUESTS_HELP = []Help{helpUp, helpDown, helpSwitchToSettingsScreen}

type PullRequestsScreen struct {
	*Window
	*Settings
	*Logger
	*GithubApi
	pullRequests             []*PullRequest
	SelectedPullRequestIndex int
}

type PullRequest struct {
	*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest
	order int
}

func NewPullRequestsScreen(globalState *Window, settings *Settings, logger *Logger, githubApi *GithubApi) *PullRequestsScreen {
	return &PullRequestsScreen{
		Window:    globalState,
		Settings:  settings,
		Logger:    logger,
		GithubApi: githubApi,
	}
}

const (
	PULL_REQUEST_AWAITING  = 1
	PULL_REQUEST_REJECTED  = 2
	PULL_REQUEST_COMMENTED = 3
	PULL_REQUEST_APPROVED  = 4
	PULL_REQUEST_DRAFT     = 5
)

func mapGithubPullRequestsToApplicationPullRequests(githubPullRequests []*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest, user string) []*PullRequest {
	var applicationPullRequests []*PullRequest
	for _, githubPullRequest := range githubPullRequests {
		pullRequest := PullRequest{
			getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequest: githubPullRequest,
		}

		if githubPullRequest.GetIsDraft() == true {
			pullRequest.order = PULL_REQUEST_DRAFT
		} else {
			for _, latestReviews := range githubPullRequest.GetLatestReviews().GetNodes() {
				if latestReviews.GetAuthor().GetLogin() == user {
					if latestReviews.GetState() == PullRequestReviewStateApproved {
						pullRequest.order = PULL_REQUEST_APPROVED
					} else if latestReviews.GetState() == PullRequestReviewStateChangesRequested {
						pullRequest.order = PULL_REQUEST_REJECTED
					} else if latestReviews.GetState() == PullRequestReviewStateCommented {
						pullRequest.order = PULL_REQUEST_COMMENTED
					}
				}
			}

			for _, reviewRequest := range githubPullRequest.GetReviewRequests().GetNodes() {
				requestedReviewer, ok := reviewRequest.GetRequestedReviewer().(*getRepositoryInfoRepositoryPullRequestsPullRequestConnectionNodesPullRequestReviewRequestsReviewRequestConnectionNodesReviewRequestRequestedReviewerUser)
				if ok {
					if requestedReviewer.GetLogin() == user {
						pullRequest.order = PULL_REQUEST_AWAITING
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

func sortPullRequestsForMe(pullRequestsForMe []*PullRequest, logger *Logger, username string) {
	sort.Slice(pullRequestsForMe, func(i, j int) bool {
		if pullRequestsForMe[i].order == pullRequestsForMe[j].order {
			return pullRequestsForMe[i].GetCreatedAt().After(pullRequestsForMe[j].GetCreatedAt())
		}

		if pullRequestsForMe[i].order > pullRequestsForMe[j].order {
			return false
		}

		return true
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

	r.pullRequests = mapGithubPullRequestsToApplicationPullRequests(pullRequestsForMe, r.Settings.Username)

	sortPullRequestsForMe(r.pullRequests, r.Logger, r.Settings.Username)

	r.Logger.Struct(pullRequestsForMe)

	return nil
}

func (r *PullRequestsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case helpDown.Shortcut:
				{
					if r.SelectedPullRequestIndex == len(r.pullRequests)-1 {
						r.SelectedPullRequestIndex = 0
					} else {
						r.SelectedPullRequestIndex = r.SelectedPullRequestIndex + 1
					}
				}
			case helpUp.Shortcut:
				{
					if r.SelectedPullRequestIndex == 0 {
						r.SelectedPullRequestIndex = len(r.pullRequests) - 1
					} else {
						r.SelectedPullRequestIndex = r.SelectedPullRequestIndex - 1
					}
				}
			case helpOpenPullRequest.Shortcut:
				{
					selectedPullRequest := r.pullRequests[r.SelectedPullRequestIndex]
					var err error
					switch runtime.GOOS {
					case "linux":
						{
							err = exec.Command("xdg-open", selectedPullRequest.GetUrl()).Start()
						}
					case "windows":
						{
							err = exec.Command("rundll32", "url.dll,FileProtocolHandler", selectedPullRequest.GetUrl()).Start()
						}
					case "darwin":
						{
							err = exec.Command("open", selectedPullRequest.GetUrl()).Start()
						}
					default:
						err = fmt.Errorf("unsupported platform %v", runtime.GOOS)
					}

					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	return r, nil
}

func (r *PullRequestsScreen) View() string {
	header := StyledHeader.Render("Pull requests")

	pullRequestStateToUI := map[int]string{
		PULL_REQUEST_AWAITING:  StyledAwaiting.Render("review required"),
		PULL_REQUEST_APPROVED:  StyledApproved.Render("approved"),
		PULL_REQUEST_REJECTED:  StyledChangesRequested.Render("changes requested"),
		PULL_REQUEST_DRAFT:     StyledDraft.Render("draft"),
		PULL_REQUEST_COMMENTED: StyledCommented.Render("commented"),
	}
	var pullRequestMessage string
	for i, pullRequest := range r.pullRequests {
		info, ok := pullRequestStateToUI[pullRequest.order]
		if !ok {
			r.Logger.Info(fmt.Sprintf("info does not exist for pull request state %v", pullRequest.order))
		}

		if i == r.SelectedPullRequestIndex {
			pullRequestMessage += StyledUnderline.Render(fmt.Sprintf("• %v wants to merge \"%v\"", pullRequest.GetAuthor().GetLogin(), pullRequest.GetTitle())) + " (" + info + ")\n"
		} else {
			pullRequestMessage += fmt.Sprintf("• %v wants to merge \"%v\" (%v)\n", pullRequest.GetAuthor().GetLogin(), pullRequest.GetTitle(), info)
		}
	}
	pullRequestsWrapper := wordwrap.NewWriter(r.Window.Width - StyledMain.GetHorizontalPadding())
	pullRequestsWrapper.Breakpoints = []rune{' '}
	_, err := pullRequestsWrapper.Write([]byte(pullRequestMessage))
	if err != nil {
		r.Logger.Error(err)
	}

	helpString := ""
	for _, help := range PULL_REQUESTS_HELP {
		helpString += lipgloss.JoinHorizontal(lipgloss.Left, StyledHelpShortcut.Render(help.Display), " ", StyledHelpDescription.Render(help.Description), "   ")
	}
	helpWrapper := wordwrap.NewWriter(r.Window.Width - StyledMain.GetHorizontalPadding())
	helpWrapper.Breakpoints = []rune{' '}
	_, err = helpWrapper.Write([]byte(helpString))
	if err != nil {
		r.Logger.Error(err)
	}

	return StyledMain.Render(lipgloss.JoinVertical(lipgloss.Left, header, pullRequestsWrapper.String(), helpWrapper.String()))
}
