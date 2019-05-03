package sergeant

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/michaelgrifalconi/gitpd/pkg/configuration"
)

const workDir = "/tmp/repos/"
const resultDir = "/tmp/results/"

var ctx context.Context
var client *github.Client

type Seargent struct {
	config *configuration.SeargentConf
}
type MyRepo struct {
	repo *github.Repository
	path string
}
type MyGist struct {
	gist *github.Gist
	path string
}

func (s *Seargent) Setup(c *configuration.SeargentConf) {
	s.config = c
}

func enqueueJob(item func(), c *configuration.SeargentConf) {
	c.ExecutionQueue <- true
	go func() {
		item()
		<-c.ExecutionQueue
	}()
}

// Info Function to show colored text
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func check(e error) {
	if e != nil {
		panic(e)
	} else if _, ok := e.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
	} else if _, ok := e.(*github.AcceptedError); ok {
		log.Println("scheduled on GitHub side")
	}
}

func makeDirectories() error {
	os.MkdirAll(workDir+"org", 0700)
	//	os.MkdirAll(workdir+"team", 0700)
	os.MkdirAll(workDir+"users", 0700)

	return nil
}

func (s *Seargent) Investigate() {

	s.config.ExecutionQueue = make(chan bool, s.config.Threads)

	ctx = context.Background()

	//authN
	var err error
	client, err = authenticatetogit(ctx, s.config.GitToken, s.config)
	check(err)

	//Creating some temp directories to store repos & results. These will be deleted in the end
	err = makeDirectories()
	check(err)

	//By now, we either have the org, user, repoURL or the gistURL. The program flow changes accordingly..

	switch s.config.Kind {
	case configuration.ORG:
		scanOrg(s.config)
	case configuration.REPO:
		//scanRepo() TODO: implement repo only?
	case configuration.USER:
		scanUser(s.config)
	case configuration.GIST:
		//scanGist() TODO: implement gist only?
	default:
		log.Fatalln("Failed to parse -kind flag")
	}
}

func scanOrg(c *configuration.SeargentConf) {
	Info("Listing org's repos..")
	orgRepos := listOrgRepos(ctx, client, c.Target)
	userRepos := []MyRepo{}
	userGists := []MyGist{}
	if !c.OrgOnly {
		Info("Listing org's users..")
		allUsers, err := listallusers(ctx, client, c.Target)
		check(err)

		Info("Listing user's repos and gists..")
		for _, user := range allUsers {
			userRepos = append(userRepos, listUserRepos(ctx, client, *user.Login, c)...)
			userGists = append(userGists, listUserGists(ctx, client, *user.Login, c)...)
		}
	}
	totalRepos := append(orgRepos, userRepos...)

	if !c.ScanOnly {
		//cloning all the repos of the org
		Info("Cloning repositories..")
		for _, r := range totalRepos {
			if strings.Contains(c.Blacklist, *r.repo.Name) {
				fmt.Println("Repo " + *r.repo.Name + " is in the repo blacklist, moving on..")
			} else {
				cloneRepo(r, c) //TODO:continue
			}
		}
		Info("Done cloning repos.")

		Info("Cloning users gists..")
		for _, g := range userGists {
			cloneGist(g, c)
		}
		Info("Done cloning gists.")

	}
	if !c.DownloadOnly {
		Info("Scanning ..")
		for _, r := range totalRepos {
			err := scanDir(r.path, c)
			check(err)
		}
		for _, g := range userGists {
			err := scanDir(g.path, c)
			check(err)
		}
		Info("Finished scanning")
	}
}

func scanUser(c *configuration.SeargentConf) {
	userRepos := []MyRepo{}
	userGists := []MyGist{}
	Info("Listing user's repos and gists..")

	userRepos = listUserRepos(ctx, client, c.Target, c)
	userGists = listUserGists(ctx, client, c.Target, c)

	if !c.ScanOnly {
		Info("Cloning repositories..")
		for _, r := range userRepos {
			if strings.Contains(c.Blacklist, *r.repo.Name) {
				fmt.Println("Repo " + *r.repo.Name + " is in the repo blacklist, moving on..")
			} else {
				cloneRepo(r, c)
			}
		}
		Info("Done cloning repos.")

		Info("Cloning users gists..")
		for _, g := range userGists {
			cloneGist(g, c)
		}
		Info("Done cloning gists.")
	}
	if !c.DownloadOnly {
		Info("Scanning..")
		for _, r := range userRepos {
			err := scanDir(r.path, c)
			check(err)
		}
		for _, g := range userGists {
			err := scanDir(g.path, c)
			check(err)
		}
		Info("Finished scanning")
	}
}
