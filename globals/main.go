package globals

func NewGlobals() *Globals {
	window := newWindow()
	logger := newLogger()
	settings := newSettings(logger)
	githubApi := newGithubApi(settings.GithubToken)

	settings.Load()

	return &Globals{
		githubApi: githubApi,
		logger:    logger,
		settings:  settings,
		window:    window,
	}
}

type Globals struct {
	*githubApi
	*logger
	*settings
	*window
}
