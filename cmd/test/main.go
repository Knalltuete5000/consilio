package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kevinklinger/consilio/model"
)

func main() {
	content, err := ioutil.ReadFile("F:\\Downloads\\libvirtConfig.json")
	if err == nil {
		var data []model.DynamicElement
		errPars := json.Unmarshal(content, &data)
		if errPars == nil {
			fmt.Println("Success")
		} else {
			fmt.Printf("Failed to pars objects: %s\n", errPars.Error())
		}
	} else {
		fmt.Printf("Failed to read file: %s\n", err.Error())
	}
}
