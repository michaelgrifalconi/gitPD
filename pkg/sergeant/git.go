package sergeant

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"github.com/michaelgrifalconi/gitpd/pkg/configuration"
	"golang.org/x/oauth2"
)

func gitclone(cloneURL string, repoName string) {
	cmd := exec.Command("/usr/bin/git", "clone", cloneURL, repoName)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		// panic(err)
	}
}

// Moving cloning logic out of individual functions
func cloneRepo(r MyRepo, c *configuration.SeargentConf) {
	urlToClone := ""

	switch c.ScanPrivateReposOnly { //TODO: why???
	case false:
		urlToClone = *r.repo.CloneURL
	case true:
		urlToClone = *r.repo.SSHURL
	default:
		urlToClone = *r.repo.CloneURL
	}

	if c.EnterpriseURL != "" { //TODO: why???
		urlToClone = *r.repo.SSHURL
	}

	if !c.CloneForks && *r.repo.Fork {
		fmt.Println(*r.repo.Name + " is a fork and the cloneFork flag was set to false so moving on..")
	} else {
		// clone it
		fmt.Println(urlToClone)
		func(urlToClone string, directory string) {
			enqueueJob(func() {
				gitclone(urlToClone, directory)
			}, c)
		}(urlToClone, workDir+r.path)
	}

}
func cloneGist(g MyGist, c *configuration.SeargentConf) {
	var gisturl string
	if c.EnterpriseURL != "" { //TODO: WAT?
		d := strings.Split(*g.gist.GitPullURL, "/")[2]
		f := strings.Split(*g.gist.GitPullURL, "/")[4]
		gisturl = "git@" + d + ":gist/" + f
	} else {
		gisturl = *g.gist.GitPullURL
	}
	fmt.Println(gisturl)

	func(gisturl string, directory string) {
		enqueueJob(func() {
			gitclone(gisturl, directory)
		}, c)
	}(gisturl, workDir+g.path)
}
func listUserRepos(ctx context.Context, client *github.Client, user string, c *configuration.SeargentConf) []MyRepo {
	var uname string
	var userRepos []*github.Repository
	var storedOrgRepos []MyRepo

	var opt3 *github.RepositoryListOptions

	if c.ScanPrivateReposOnly {
		uname = ""
		opt3 = &github.RepositoryListOptions{
			Visibility:  "private",
			ListOptions: github.ListOptions{PerPage: 10},
		}
	} else {
		uname = user
		opt3 = &github.RepositoryListOptions{
			ListOptions: github.ListOptions{PerPage: 10},
		}
	}

	for {
		uRepos, resp, err := client.Repositories.List(ctx, uname, opt3)
		check(err)
		userRepos = append(userRepos, uRepos...) //adding to the userRepos array
		if resp.NextPage == 0 {
			break
		}
		opt3.Page = resp.NextPage
	}

	var r MyRepo
	for _, repo := range userRepos {
		r = MyRepo{repo: repo, path: "user/" + user + "/" + *repo.Name}
		storedOrgRepos = append(storedOrgRepos, r)
	}
	return storedOrgRepos
}

func listOrgRepos(ctx context.Context, client *github.Client, orgName string) []MyRepo {
	var storedOrgRepos []MyRepo
	var orgRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
		check(err)
		orgRepos = append(orgRepos, repos...) //adding to the repo array
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var r MyRepo
	for _, repo := range orgRepos {
		r = MyRepo{repo: repo, path: "org/" + orgName + "/" + *repo.Name}
		storedOrgRepos = append(storedOrgRepos, r)
	}
	return storedOrgRepos
}

func listUserGists(ctx context.Context, client *github.Client, user string, c *configuration.SeargentConf) []MyGist {
	var storedUserGist []MyGist

	var userGists []*github.Gist
	opt4 := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		uGists, resp, err := client.Gists.List(ctx, user, opt4)
		check(err)
		userGists = append(userGists, uGists...)
		if resp.NextPage == 0 {
			break
		}
		opt4.Page = resp.NextPage
	}

	var r MyGist
	for _, gist := range userGists {
		r = MyGist{gist: gist, path: "user/" + user + "/" + *gist.ID}
		storedUserGist = append(storedUserGist, r)
	}
	return storedUserGist

}

func listallusers(ctx context.Context, client *github.Client, org string) ([]*github.User, error) {
	Info("Listing users of the organization and their repositories and gists")
	var allUsers []*github.User
	opt2 := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		users, resp, err := client.Organizations.ListMembers(ctx, org, opt2)
		check(err)
		allUsers = append(allUsers, users...) //adding to the allUsers array
		if resp.NextPage == 0 {
			break
		}
		opt2.Page = resp.NextPage
	}

	return allUsers, nil
}

func authenticatetogit(ctx context.Context, token string, c *configuration.SeargentConf) (*github.Client, error) {
	var client *github.Client
	var err error
	fmt.Println("TOKEN: ", token)

	//Authenticating to Github using the token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	if c.EnterpriseURL == "" {
		client = github.NewClient(tc)
	} else if c.EnterpriseURL != "" {
		client, err = github.NewEnterpriseClient(c.EnterpriseURL, c.EnterpriseURL, tc)
		if err != nil {
			fmt.Printf("NewEnterpriseClient returned unexpected error: %v", err)
		}
	}
	return client, nil
}
