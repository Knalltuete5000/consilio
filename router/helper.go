package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/pkg/errors"
)

func getURL(r *http.Request) string {
	if r.URL == nil {
		return ""
	}
	return r.URL.String()
}

func decodeBody(r *http.Request, result interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		log.Println(errors.WithStack(err))
		if err == io.EOF {
			err := errors.New("Request body not found")
			return err
		}
		return err
	}
	return nil
}

func cdktfDeploy(directory string, command string) {
	app := "/usr/bin/cdktf"
	arg0 := "deploy"
	arg1 := "--auto-approve"
	var cmd *exec.Cmd
	if command != "" {
		arg2 := "--app"
		cmd = exec.Command(app, arg0, arg1, arg2, command)
	} else {
		cmd = exec.Command(app, arg0, arg1)
	}
	cmd.Dir = directory
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run cdktf deploy: %s", err.Error())
	}
}

func cdktfDestroy(directory string, command string) {
	app := "/usr/bin/cdktf"
	arg0 := "destroy"
	arg1 := "--auto-approve"
	var cmd *exec.Cmd
	if command != "" {
		arg2 := "--app"
		cmd = exec.Command(app, arg0, arg1, arg2, command)
	} else {
		cmd = exec.Command(app, arg0, arg1)
	}
	cmd.Dir = directory
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run cdktf destroy: %s", err.Error())
	}
}
