package hangman

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Harvest is the result of an execution done by the function Reaper
type Harvest struct {
	OriginalCommand string `json:"orignal_command"`
	ExecutedCommand string `json:"executed_command"`
	ReturnCode      int    `json:"return_code"`
	TimeoutReached  bool   `json:"timeout_reached"`
	Pid             int    `json:"pid"`
	Stdout          string `json:"stdout"`
	Stderr          string `json:"stderr"`
}

// Reaper execute a program with is parameters as a string, with a timeout limiting execution time
func Reaper(cmdline string, timeout uint32) Harvest {
	//cmdline = "sh -c " + cmdline
	expandedCmdline := os.ExpandEnv(cmdline)
	cmdSplit := strings.Split(strings.TrimSpace(expandedCmdline), " ")

	var cmd *exec.Cmd
	if len(cmdSplit) > 1 {
		cmd = exec.Command(cmdSplit[0], cmdSplit[1:]...)
	} else {
		cmd = exec.Command(cmdSplit[0])
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var h Harvest
	h.ExecutedCommand = expandedCmdline
	h.OriginalCommand = cmdline

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Can not execute command %s:  %v\n", h.ExecutedCommand, err)
		h.Pid = -1
		h.ReturnCode = 666
		h.TimeoutReached = false
		h.Stderr = err.Error()

		return h
	}

	h.Pid = cmd.Process.Pid
	h.ReturnCode = 0
	h.TimeoutReached = false

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
		close(done)
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Some error happened when trying to kill the process")
		}
		h.ReturnCode = 127
		h.TimeoutReached = true
		h.Stderr = stderr.String()
		h.Stdout = stdout.String()
		return h

	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command \"%s\" returned an error: %s\n", h.ExecutedCommand, err.Error())
			if strings.HasPrefix(err.Error(), "exit status") {
				h.ReturnCode, _ = strconv.Atoi(err.Error()[12:])
			} else {
				h.ReturnCode = 666
			}
		}

		h.Stderr = stderr.String()
		h.Stdout = stdout.String()
		return h
	}
}
