package hangman

import (
	"bytes"
	"fmt"
	"os/exec"
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
	h.Pid = cmd.Process.Pid
	h.ReturnCode = 0
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

	case err := <-done:
		if err != nil {
			fmt.Printf("process done with error = %v", err)
		}
	}

	h.Stdout = stdout.String()
	h.Stderr = stderr.String()

	return h
}
