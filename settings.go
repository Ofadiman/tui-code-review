package main

import (
	"encoding/json"
	"os"
)

type Settings struct {
	GithubToken    string   `json:"github_token,omitempty"`
	Username       string   `json:"username,omitempty"`
	Repositories   []string `json:"repositories,omitempty"`
	ConfigFilePath string
	*Logger
}

func NewSettings(logger *Logger) *Settings {
	home, _ := os.UserHomeDir()

	return &Settings{
		ConfigFilePath: home + "/" + ".tui-code-review.json",
		Logger:         logger,
	}
}

func (r *Settings) Load() {
	_, err := os.Stat(r.ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.Save()
		} else {
			r.Logger.Info("could not stat configuration file")
			r.Logger.Error(err)
			panic(err)
		}
	}

	bytes, err := os.ReadFile(r.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, r)
	if err != nil {
		r.Logger.Info("could not unmarshal configuration file")
		r.Logger.Error(err)
		panic(err)
	}

	r.Logger.Struct(r)
}

func (r *Settings) Save() {
	bytes, err := json.Marshal(r)
	if err != nil {
		r.Logger.Info("could not stat configuration file")
		r.Logger.Error(err)
		panic(err)
	}

	err = os.WriteFile(r.ConfigFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (r *Settings) UpdateGitHubToken(token string) {
	r.GithubToken = token
	r.Save()
}

func (r *Settings) UpdateUsername(username string) {
	r.Username = username
	r.Save()
}

func (r *Settings) AddRepositoryUrl(repositoryUrl string) {
	r.Repositories = append(r.Repositories, repositoryUrl)
	r.Save()
}

func (r *Settings) DeleteRepositoryUrl(repositoryUrl string) {
	var updatedRepositories []string
	for _, url := range r.Repositories {
		if url == repositoryUrl {
			continue
		}

		updatedRepositories = append(updatedRepositories, url)
	}

	r.Repositories = updatedRepositories
}
