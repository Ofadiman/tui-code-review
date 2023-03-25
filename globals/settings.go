package globals

import (
	"encoding/json"
	"os"
)

type settings struct {
	GithubToken    string   `json:"github_token,omitempty"`
	Repositories   []string `json:"repositories,omitempty"`
	ConfigFilePath string
	*logger
}

func newSettings(logger *logger) *settings {
	home, _ := os.UserHomeDir()

	return &settings{
		ConfigFilePath: home + "/" + ".tui-code-review.json",
		logger:         logger,
	}
}

func (r *settings) Load() {
	_, err := os.Stat(r.ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.Save()
		} else {
			r.logger.Info("could not stat configuration file")
			r.logger.Error(err)
			panic(err)
		}
	}

	bytes, err := os.ReadFile(r.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, r)
	if err != nil {
		r.logger.Info("could not unmarshal configuration file")
		r.logger.Error(err)
		panic(err)
	}

	r.logger.Struct(r)
}

func (r *settings) Save() {
	bytes, err := json.Marshal(r)
	if err != nil {
		r.logger.Info("could not stat configuration file")
		r.logger.Error(err)
		panic(err)
	}

	err = os.WriteFile(r.ConfigFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (r *settings) UpdateGitHubToken(token string) {
	r.GithubToken = token
	r.Save()
}

func (r *settings) AddRepositoryUrl(repositoryUrl string) {
	r.Repositories = append(r.Repositories, repositoryUrl)
	r.Save()
}

func (r *settings) DeleteRepositoryUrl(repositoryUrl string) {
	var updatedRepositories []string
	for _, url := range r.Repositories {
		if url == repositoryUrl {
			continue
		}

		updatedRepositories = append(updatedRepositories, url)
	}

	r.Repositories = updatedRepositories
}
