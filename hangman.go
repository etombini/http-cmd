package hangman

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Harvest struct {
	Cmdline         string
	Return_code     int
	Timeout_reached bool
	Pid             int
	Stdout          string
	Stderr          string
}

func Reaper(cmdline string, timeout int) Harvest {
	//cmdline = "sh -c " + cmdline
	cmd_split := strings.Split(strings.TrimSpace(cmdline), " ")

	var cmd *exec.Cmd
	if len(cmd_split) > 1 {
		cmd = exec.Command(cmd_split[0], cmd_split[1:]...)
	} else {
		cmd = exec.Command(cmd_split[0])
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Some error happened")
	}

	var h Harvest
	h.Pid = cmd.Process.Pid
	h.Return_code = 0
	h.Timeout_reached = false

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Some error happened when trying to kill the process")
		}
		h.Return_code = 127
		h.Timeout_reached = true

	case err := <-done:
		if err != nil {
			fmt.Printf("process done with error = %v", err)
		}
	}

	h.Stdout = stdout.String()
	h.Stderr = stderr.String()

	return h
}
