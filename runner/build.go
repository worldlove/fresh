package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func build() (string, bool) {
	buildLog("Building...")
	build := "go build -o" + buildPath() + root()
	cmdList := append(cmdBeforeBuild(), build)
	cmdList = append(cmdList, cmdAfterBuild()...)
	for _, cmdstr := range cmdList {
		cmd := exec.Command(cmdstr)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fatal(err)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fatal(err)
		}

		err = cmd.Start()
		if err != nil {
			fatal(err)
		}

		io.Copy(os.Stdout, stdout)
		errBuf, _ := ioutil.ReadAll(stderr)

		err = cmd.Wait()
		if err != nil {
			return string(errBuf), false
		}

	}
	return "", true
}
