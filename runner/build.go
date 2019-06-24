package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func build() (string, bool) {
	beforeBuild()
	go afterBuild()

	buildLog("Building...")

	cmd := exec.Command("go", "build", "-o", buildPath(), root())

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
	hasBuild <- true

	return "", true
}

func beforeBuild() {
	buildLog("BeforeBuilding...")
	for _, cmdstr := range cmdBeforeBuild() {
		cmdstr = strings.Trim(cmdstr, " ")
		buildLog("  BeforeBuild %s", cmdstr)
		cmdarr := strings.Split(cmdstr, " ")

		output, err := exec.Command(cmdarr[0], cmdarr[1:]...).Output()
		// err := cmd.Run()
		if err != nil {
			fatal(err)
		}
		if len(output) != 0 {
			buildLog("  BeforeBuild out:\n%s", output)
		}
	}
}

var hasBuild = make(chan bool)

func afterBuild() {
	buildLog("AfterBuilding...")
	select {
	case <-hasBuild:
		{
			for _, cmdstr := range cmdAfterBuild() {
				if cmdstr == "" {
					continue
				}
				buildLog("  AfterBuild %s", cmdstr)
				cmdstr = strings.Trim(cmdstr, " ")
				cmdarr := strings.Split(cmdstr, " ")
				output, err := exec.Command(cmdarr[0], cmdarr[1:]...).Output()
				if err != nil {
					fatal(err)
				}
				if len(output) != 0 {
					buildLog("  AfterBuild out:\n%s", output)
				}
			}
		}
	}
}
