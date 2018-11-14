package bind

import (
	"os/exec"
	"strings"
	"fmt"
	"bytes"
)

type Response struct {
	Out string
	Err string
}

// IC.Execute "sh", "-c", "echo stdout; echo 1>&2 stderr"
func ExecuteSync(name string, arg ...string) (r *Response) {

	r = &Response{}

	log.Infof("Execute [SYNC] command: %s %s", name, strings.Trim(fmt.Sprint(arg), "[]"))

	// https://golang.org/pkg/os/exec/#example_Cmd_Start
	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		r.Err = err.Error()
		return
	}

	if err := cmd.Wait(); err != nil {
		r.Err = err.Error()
		return
	}

	r.Out = strings.TrimSuffix(stdout.String(), "\n")
	r.Err = strings.TrimSuffix(stderr.String(), "\n")

	return
}

func ExecuteAsync(name string, arg ...string) (r *Response) {

	r = &Response{}

	log.Infof("Execute [ASYNC] command: %s %s", name, strings.Trim(fmt.Sprint(arg), "[]"))

	go func() {
		b, err := exec.Command(name, arg...).Output()
		if err != nil {
			r.Err = err.Error()
			return
		}

		log.Infof("%s", b)
	}()

	return
}
