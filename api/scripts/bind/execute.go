package bind

import (
	"os/exec"
	"io/ioutil"
	"github.com/e154/smart-home/api/log"
	"strings"
	"fmt"
)


// exec.command "sh", "-c", "echo stdout; echo 1>&2 stderr"
func Execute(name string, arg ...string) {

	go func() {
		log.Infof("Execute command: %s %s", name, strings.Trim(fmt.Sprint(arg), "[]"))

		// https://golang.org/pkg/os/exec/#example_Cmd_Start
		cmd := exec.Command(name, arg...)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Error(err.Error())
			return
		}

		if err := cmd.Start(); err != nil {
			log.Error(err.Error())
			return
		}

		slurp, _ := ioutil.ReadAll(stderr)
		log.Infof("Result: %s", string(slurp))

		if err := cmd.Wait(); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	return
}
