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
	Command        string
	ReturnCode     int
	TimeoutReached bool
	Pid            int
	Stdout         string
	Stderr         string
}

// Reaper execute a program with is parameters as a string, with a timeout limiting execution time
func Reaper(cmdline string, timeout int) Harvest {
	//cmdline = "sh -c " + cmdline
	cmdSplit := strings.Split(strings.TrimSpace(cmdline), " ")

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

	if err := cmd.Start(); err != nil {
		fmt.Println("Some error happened")
	}

	var h Harvest
	h.Command = cmdline
	if cmd.Process == nil {
		fmt.Fprintf(os.Stderr, "Process is nil, there is a problem\n")
		h.Pid = -1
		h.ReturnCode = 666
	} else {
		h.Pid = cmd.Process.Pid
		h.ReturnCode = 0
	}
	h.TimeoutReached = false

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
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

	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command \"%s\" returned an error = %s\n", h.Command, err.Error())
			if strings.HasPrefix(err.Error(), "exit status") {
				h.ReturnCode, _ = strconv.Atoi(err.Error()[12:])
			} else {
				h.ReturnCode = 666
			}
			h.Stderr = err.Error()
			if stderr.String() != "" {
				h.Stderr = h.Stderr + "\n" + stderr.String()
			}
		} else {
			h.ReturnCode = 0
			h.Stderr = stderr.String()
		}
		h.Stdout = stdout.String()
	}

	return h
}
