package main

type Help struct {
	Shortcut    string
	Description string
	Display     string
}

var helpQuit = Help{
	Shortcut:    "ctrl+q",
	Description: "Quit program",
	Display:     "Ctrl + Q",
}

var helpUpdateGithubToken = Help{
	Shortcut:    "ctrl+t",
	Description: "Update GitHub Token",
	Display:     "Ctrl + T",
}

var helpAddGitHubRepositoryUrl = Help{
	Shortcut:    "ctrl+r",
	Description: "Add GitHub repository URL",
	Display:     "Ctrl + R",
}

var helpOpenGitHubRepositoryUrl = Help{
	Shortcut:    "enter",
	Description: "Open selected GitHub repository url",
	Display:     "Enter",
}

var helpDeleteGitHubRepositoryUrl = Help{
	Shortcut:    "delete",
	Description: "Delete selected GitHub repository url",
	Display:     "Enter",
}

var helpDown = Help{
	Shortcut:    "j",
	Description: "Down",
	Display:     "J",
}

var helpUp = Help{
	Shortcut:    "k",
	Description: "Up",
	Display:     "K",
}

var helpEscape = Help{
	Shortcut:    "esc",
	Description: "Close current window",
	Display:     "Escape",
}

var helpSwitchToSettingsScreen = Help{
	Shortcut:    "ctrl+s",
	Description: "Switch to settings screen",
	Display:     "Ctrl + S",
}

var helpSwitchToPullRequestsScreen = Help{
	Shortcut:    "ctrl+p",
	Description: "Switch to pull requests screen",
	Display:     "Ctrl + P",
}
