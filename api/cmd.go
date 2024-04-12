package api

import (
	"fmt"
	"log"
	"os/exec"
)

func ExecCmd(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	o, err := cmd.Output()
	if err != nil {
		log.Printf("run command fatal: %v\n", cmd.Args)
		return "", fmt.Errorf("error executing command: %v\nerror: %v", cmd.Args, err)
	}
	log.Printf("run command success: %v\n", cmd.Args)
	return string(o), nil
}
