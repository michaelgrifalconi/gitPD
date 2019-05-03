package sergeant

import (
	"fmt"

	"github.com/michaelgrifalconi/gitpd/pkg/configuration"
)

func runTruffleGopher(dir string, c *configuration.SeargentConf) error {
	fmt.Println("Scanning:", dir)
	return nil
}

// Moving directory scanning logic out of individual functions
func scanDir(dir string, c *configuration.SeargentConf) error {

	func(dir string, c *configuration.SeargentConf) { //TODO: review
		enqueueJob(func() {
			runTruffleGopher(dir, c)
		}, c)
	}(dir, c)

	return nil
}
