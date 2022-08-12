package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/kevinklinger/consilio/libs/converters"
	validator "github.com/kevinklinger/consilio/libs/validators"
	"github.com/kevinklinger/consilio/model"
)

func (s *ConsilioRouter) handleCreateProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO create a new project in the DB and to return the ID
		fmt.Println("create project called")
		fmt.Printf("Bodysize %d", r.ContentLength)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Failed to read body %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(data))
	}
}

func (s *ConsilioRouter) handleUpdateProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//projectID := p.ByName("id")
		// TODO overwrite the config for the given project
		fmt.Println("update project called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Failed to read body %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var object []model.DynamicElement
		err = json.Unmarshal(data, &object)
		if err != nil {
			fmt.Printf("Failed to parse data %s\n", err.Error())
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if !validator.HasRequiredFields(object) {
			fmt.Printf("Not all required fields are supplied.\n")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		libvirt_obj, err := converters.ConvertToLibvirtConfig(object)
		if err != nil {
			fmt.Printf("Failed to parse objects: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		templ, err := template.New("libvirt.go.template").ParseFiles("./templates/libvirt/libvirt.go.template")
		if err != nil {
			fmt.Printf("Faild to generate template: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		f, err := os.Create("./templates/libvirt/test.go")
		if err != nil {
			fmt.Printf("Faild to create file to write template to: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = templ.Execute(f, libvirt_obj)
		if err != nil {
			fmt.Printf("Faild to execute template: %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Success :)")
	}
}

func (s *ConsilioRouter) handleGetProject() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//projectID := p.ByName("id")
		// TODO read in information for the project with the given ID

		fmt.Println("get project called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Failed to read body %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(data))

	}
}

func (s *ConsilioRouter) handleGetAllProjects() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO return all projects of the given user
		fmt.Println("get all projects called")
		fmt.Printf("Bodysize %d\n", r.ContentLength)

		data, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Failed to read body %s\n", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(data))
	}
}
