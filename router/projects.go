package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"fmt"

	"io"

	"encoding/json"

	"github.com/kevinklinger/consilio/model"
)

func (s *ConsilioRouter) handleCreateProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO create a new project in the DB and to return the ID
		fmt.Println("create project called")
		fmt.Printf("Bodysize %d", r.ContentLength)
		byt, errRead := io.ReadAll(r.Body)
		if errRead == nil {
			fmt.Println(byt)
		} else {
			fmt.Printf("Failed to read body %s\n", errRead.Error())
		}
	}
}

func (s *ConsilioRouter) handleUpdateProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//projectID := p.ByName("id")
		// TODO overwrite the config for the given project
		fmt.Println("update project called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)
		var data model.Content
		byt, errRead := io.ReadAll(r.Body)
		if errRead == nil {
			fmt.Println(byt)
			errPars := json.Unmarshal(byt, &data)
			if errPars == nil {
				fmt.Println("Success parsing")
			} else {
				fmt.Printf("Failed to pars data: %s\n", errPars.Error())
			}
		} else {
			fmt.Printf("Failed to read body %s\n", errRead.Error())
		}
	}
}

func (s *ConsilioRouter) handleGetProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//projectID := p.ByName("id")
		// TODO read in information for the project with the given ID
		content := p.ByName("content")
		fmt.Println("get project called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)
		fmt.Println(content)
		byt, errRead := io.ReadAll(r.Body)
		if errRead == nil {
			fmt.Println(byt)
		} else {
			fmt.Printf("Failed to read body %s\n", errRead.Error())
		}
	}
}

func (s *ConsilioRouter) handleGetAllProjects() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO return all projects of the given user
		fmt.Println("get all projects called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)
		byt, errRead := io.ReadAll(r.Body)
		if errRead == nil {
			fmt.Println(byt)
		} else {
			fmt.Printf("Failed to read body %s\n", errRead.Error())
		}
	}
}
