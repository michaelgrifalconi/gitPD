package configuration

const (
	ORG  = 1
	USER = 2
	REPO = 3
	GIST = 4
)

type SeargentConf struct {
	Kind                 int
	Target               string
	GitToken             string
	OutputFile           string
	CloneForks           bool
	OrgOnly              bool
	TeamName             string //TODO: use it?
	ScanPrivateReposOnly bool   //TODO: this affect only user repos and single repos, not org
	EnterpriseURL        string
	Threads              int
	MergeOutput          bool //TODO: use this?
	Blacklist            string
	ExecutionQueue       chan bool
	ScanOnly             bool
	DownloadOnly         bool
}

type Configuration struct {
	S SeargentConf
}
