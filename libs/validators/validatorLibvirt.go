package validator

import (
	"github.com/kevinklinger/consilio/model"
	libvirt "github.com/kevinklinger/consilio/router/kvm"
)

func HasRequiredFields(dynamicElements []model.DynamicElement) bool {
	if len(dynamicElements) <= 0 {
		return true
	}

	requiredFields := libvirt.GetRequiredLibvirtFieldsByNames(getNamesOfElements((dynamicElements)))
	for _, dynamicElement := range requiredFields {
		toCheck := getDynamicModel(dynamicElement.Name, dynamicElements)
		if toCheck.Name == "" {
			return false
		}
		for _, field := range dynamicElement.Fields {
			fieldToCheck := getField(field.Name, toCheck.Fields)
			if fieldToCheck.Name == "" {
				return false
			}
		}
	}
	return true
}

func getNamesOfElements(dynElm []model.DynamicElement) []string {
	var names []string
	for _, elm := range dynElm {
		names = append(names, elm.Name)
	}
	return names
}

func getDynamicModel(name string, dynamicModels []model.DynamicElement) model.DynamicElement {
	for _, dyn := range dynamicModels {
		if dyn.Name == name {
			return dyn
		}
	}
	return model.DynamicElement{Name: ""}
}

func getField(name string, fieldTypes []model.FieldType) model.FieldType {
	for _, field := range fieldTypes {
		if field.Name == name {
			return field
		}
	}
	return model.FieldType{Name: ""}
}

func containsField(name string, fields []model.FieldType) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}
	return false
}
