package github

import (
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Token        string
	Organization string
	BaseURL      string
}

type Organization struct {
	name   string
	client *github.Client
}

// Client configures and returns a fully initialized GithubClient
func (c *Config) Client() (interface{}, error) {
	var org Organization
	org.name = c.Organization
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	org.client = github.NewClient(tc)
	if c.BaseURL != "" {
		u, err := url.Parse(c.BaseURL)
		if err != nil {
			return nil, err
		}
		org.client.BaseURL = u
	}
	return &org, nil
}

func (o *Organization) Fork(owner, repository, organization string) error {
	var opt *github.RepositoryCreateForkOptions
	if organization != "" {
		opt = &github.RepositoryCreateForkOptions{Organization: organization}
	}

	_, _, err := o.client.Repositories.CreateFork(owner, repository, opt)
	if err != nil {
		return err
	}

	return nil
}
