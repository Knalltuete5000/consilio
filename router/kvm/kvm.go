package libvirt

import (
	libvirtProvider "github.com/dmacvicar/terraform-provider-libvirt/libvirt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kevinklinger/consilio/libs"
	"github.com/kevinklinger/consilio/model"
)

var (
	provider = libvirtProvider.Provider().(*schema.Provider)
)

// GetLibvirtFields returns a list of elements the terraform provider for libvirt supports
func GetLibvirtFields() []model.DynamicElement {
	var result []model.DynamicElement

	for rscName, rscAttr := range provider.ResourcesMap {
		fields := libs.ExtractFields(rscAttr)
		if len(fields) != 0 {
			result = append(result, model.DynamicElement{
				Name:   rscName,
				Fields: fields,
			})
		}
	}

	return result
}

func GetRequiredLibvirtFields() []model.DynamicElement {
	var result []model.DynamicElement

	for rscName, rscAttr := range provider.ResourcesMap {
		fields := libs.ExtractRequiredFields(rscAttr)
		if len(fields) != 0 {
			result = append(result, model.DynamicElement{
				Name:   rscName,
				Fields: fields,
			})
		}
	}
	return result
}

func GetRequiredLibvirtFieldsByNames(names []string) []model.DynamicElement {
	var result []model.DynamicElement

	for rscName, rscAttr := range provider.ResourcesMap {
		if libs.Contains(rscName, names) {
			fields := libs.ExtractRequiredFields(rscAttr)
			if len(fields) != 0 {
				result = append(result, model.DynamicElement{
					Name:   rscName,
					Fields: fields,
				})
			}
		}
	}
	return result
}

func ToLibvirtFields(models []model.DynamicElement) map[string]*schema.Resource {
	var resourceMap map[string]*schema.Resource
	for _, model := range models {
		var resource schema.Resource
		for _, field := range model.Fields {
			resource.Schema[field.Name] = CreateSchema(field)
		}
	}
	return resourceMap
}

func CreateSchema(field model.FieldType) *schema.Schema {
	var s *schema.Schema

	return s
}

func HasSubFields(field model.FieldType) bool {
	return field.Subfields != nil
}
