package main

import (
	"flag"

	"github.com/michaelgrifalconi/gitpd/pkg/configuration"
	"github.com/michaelgrifalconi/gitpd/pkg/sergeant"
)

func main() {

	c := configuration.Configuration{}
	flag.StringVar(&c.S.Target, "target", "", "Org/User name or Repo/Gist URL to target")
	flag.StringVar(&c.S.GitToken, "token", "", "Github Personal Access Token. This is required.")
	flag.StringVar(&c.S.OutputFile, "output", "results.txt", "Output file to save the results.")
	flag.BoolVar(&c.S.CloneForks, "cloneForks", false, "Option to clone org and user repos that are forks. Default is false")
	flag.BoolVar(&c.S.OrgOnly, "orgOnly", false, "Option to skip cloning user repo's when scanning an org. Default is false")
	flag.StringVar(&c.S.TeamName, "teamName", "", "Name of the Organization Team which has access to private repositories for scanning.")
	flag.BoolVar(&c.S.ScanPrivateReposOnly, "scanPrivateReposOnly", false, "Option to scan private repositories only. Default is false")
	flag.StringVar(&c.S.EnterpriseURL, "enterpriseURL", "", "Base URL of the Github Enterprise")
	flag.IntVar(&c.S.Threads, "threads", 10, "Amount of parallel threads")
	flag.BoolVar(&c.S.MergeOutput, "mergeOutput", false, "Merge the output files of all the tools used into one JSON file")
	flag.StringVar(&c.S.Blacklist, "blacklist", "", "Comma seperated values of Repos to Skip Scanning for")
	flag.BoolVar(&c.S.ScanOnly, "scanOnly", false, "Just scan, do not download. Please make sure to mount a volume with correct file structure.")       //TODO: improve docs about this
	flag.BoolVar(&c.S.DownloadOnly, "downloadOnly", false, "Just download, do not scan. Please make sure to mount a volume to retain downloaded data.") //TODO: improve docs about this

	kind := flag.String("kind", "", "Kind of target. [ORG|USER|REPO|GIST]")
	flag.Parse()

	switch *kind {
	case "ORG":
		c.S.Kind = configuration.ORG
	case "USER":
		c.S.Kind = configuration.USER
	case "REPO":
		c.S.Kind = configuration.REPO
	case "GIST":
		c.S.Kind = configuration.GIST
	}

	s := sergeant.Seargent{}
	s.Setup(&c.S)
	s.Investigate()
}

//TODO: review when to use ssh
