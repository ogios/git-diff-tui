package api

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

var (
	EMPTY_CMD_RESULT = []byte{}
	maxCmdBufSize    = 16 * 1024
	outbuf, errbuf   bytes.Buffer
)

func ExecCmd(args ...string) ([]byte, error) {
	defer func() {
		if outbuf.Len() > maxCmdBufSize {
			outbuf = bytes.Buffer{}
		} else {
			outbuf.Reset()
		}
		if errbuf.Len() > maxCmdBufSize {
			errbuf = bytes.Buffer{}
		} else {
			errbuf.Reset()
		}
	}()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	// o, err := cmd.Output()
	if err != nil {
		log.Printf("run command fatal: %v\n", cmd.Args)
		// return EMPTY_CMD_RESULT, fmt.Errorf("error executing command: %v\nerror: %v", cmd.Args, err)
		return errbuf.Bytes(), fmt.Errorf("error executing command: %v\nerror: %v", cmd.Args, err)
	}
	log.Printf("run command success: %v\n", cmd.Args)
	return outbuf.Bytes(), nil
}
