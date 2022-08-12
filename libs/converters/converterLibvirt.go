package converters

/*
We receive a number of dynamic elements which contain our basic information but we do not know the
dependencies to other resources. We try to find out these dependencies here, before we generate
the source code for the cdktf via a template.

Following are possible dependencies for the resources (a depends on b)

Pool (no dependencies) ✅

CloudInit -> Pool ✅

CoreosIgnition -> Pool ✅

Volume -> Pool ✅
	   -> Volume (not itself, base_volume) ✅

Network (no dependencies) ✅

Domain -> Volume (boot device or disk) ✅
	   -> Network ❌
	   -> CloudInit ✅
	   -> CoreosIgnition ✅

check for cyclic dependencies ✅
	Volume -> same Volume ✅
*/

import (
	"errors"
	"fmt"

	"github.com/kevinklinger/consilio/model"
)

type LibvirtConfig struct {
	Domain    *LibvirtDomainConfig // for a typical terraform provider this should be arrays instead of a single object
	Pool      *LibvirtPoolConfig
	Network   *LibvirtNetworkConfig
	Volume    *LibvirtVolumeConfig
	CloudInit *LibvirtCloudInitConfig
	Ignition  *LibvirtIgnitionConfig
}

type LibvirtDomainConfig struct {
	dependsOn []string
	Domain    model.DynamicElement
}

type LibvirtPoolConfig struct {
	dependsOn []string
	Pool      model.DynamicElement
}

type LibvirtNetworkConfig struct {
	dependsOn []string
	Network   model.DynamicElement
}

type LibvirtVolumeConfig struct {
	dependsOn []string
	Volume    model.DynamicElement
}

type LibvirtCloudInitConfig struct {
	dependsOn []string
	CloudInit model.DynamicElement
}

type LibvirtIgnitionConfig struct {
	dependsOn []string
	Ignition  model.DynamicElement
}

func ConvertToLibvirtConfig(dynamicElements []model.DynamicElement) (LibvirtConfig, error) {
	var configs LibvirtConfig
	for _, model := range dynamicElements {
		switch model.Name {
		case "libvirt_cloudinit_disk":
			configs.CloudInit.CloudInit = model
			configs.CloudInit.dependsOn = checkDependenciesForCloudInitConfig(model.Fields, dynamicElements)
		case "libvirt_domain":
			configs.Domain.Domain = model
			configs.Domain.dependsOn = checkDependenciesForDomainConfig(model.Fields, dynamicElements)
		case "libvirt_ignition":
			configs.Ignition.Ignition = model
			configs.Ignition.dependsOn = checkDependenciesForIgnitionConfig(model.Fields, dynamicElements)
		case "libvirt_network":
			configs.Network.Network = model
			configs.Network.dependsOn = checkDependenciesForNetworkConfig(model.Fields, dynamicElements)
		case "libvirt_pool":
			configs.Pool.Pool = model
			configs.Pool.dependsOn = checkDependenciesForNetworkConfig(model.Fields, dynamicElements)
		case "libvirt_volume":
			configs.Volume.Volume = model
			configs.Volume.dependsOn = checkDependenciesForVolumeConfig(model.Fields, dynamicElements)
		default:
		}
	}

	er := checkCyclicDependencies(configs)
	if er != nil {
		return configs, er
	}

	return configs, nil
}

func checkCyclicDependencies(config LibvirtConfig) error {
	// volume does not refer to itself
	if config.Volume != nil {
		var name string = GetName(config.Volume.Volume.Fields)
		for _, dep := range config.Volume.dependsOn {
			if dep == name {
				return errors.New("volume depends on itself")
			}
		}
	}

	return nil
}

func checkDependenciesForCloudInitConfig(cloudinit []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// check for pool (field pool)
	// get pool
	var pool_name = fmt.Sprint(GetFieldValue(cloudinit, "name"))

	for _, m := range models {
		if m.Name == "libvirt_pool" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == pool_name {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		}
	}

	return deps
}

func checkDependenciesForDomainConfig(domain []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// check for volume
	//	fields:
	//		disk (multiple) -> requires volume id (we do not have a volume id to compare to),
	///		kernel
	// check for cloudinit  (field  cloudinit)
	// check for coreosignition (field coreos_ignition)
	// check for network (field network_interface (multiple))
	//	Not possible to connect the networkinterface to a domain
	//	Web does not contain any field to reference a network
	var coreos_ignition = fmt.Sprint(GetFieldValue(domain, "coreos_ignition"))
	var kernel = fmt.Sprint(GetFieldValue(domain, "kernel"))
	var cloud_init_disk = fmt.Sprint(GetFieldValue(domain, "cloudinit"))

	for _, m := range models {
		if m.Name == "libvirt_ignition" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == coreos_ignition {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		} else if m.Name == "libvirt_volume" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == kernel {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		} else if m.Name == "libvirt_cloudinit_disk" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == cloud_init_disk {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		} /*else if m.Name == "libvirt_network" {
			for _, f := range m.Fields {
				if f.Name == "network_interface" {
					for _, sf := range *f.Subfields {

					}
				}
			}
		}*/
	}

	return deps
}

func checkDependenciesForIgnitionConfig(ignition []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// check for pool (field pool)
	var pool_name = fmt.Sprint(GetFieldValue(ignition, "pool"))

	for _, m := range models {
		if m.Name == "libvirt_pool" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == pool_name {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		}
	}
	return deps
}

func checkDependenciesForNetworkConfig(network []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// done, no dependencies
	return deps
}

func checkDependenciesForPoolConfig(pool []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// done, no dependencies
	return deps
}

func checkDependenciesForVolumeConfig(volume []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// check for Pool (field pool, base_volume_pool)
	// check for other volumes as base volumes (fields base_volume, base_volume_id (we do not have any ids))
	var base_volume = fmt.Sprint(GetFieldValue(volume, "base_volume_name"))
	var base_volume_pool = fmt.Sprint(GetFieldValue(volume, "base_volume_pool"))
	var pool = fmt.Sprint(GetFieldValue(volume, "pool"))

	for _, m := range models {
		if m.Name == "libvirt_pool" {
			for _, f := range m.Fields {
				if f.Name == "name" && (f.Value == base_volume_pool || f.Value == pool) {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		} else if m.Name == "libvirt_volume" {
			for _, f := range m.Fields {
				if f.Name == "name" && f.Value == base_volume {
					deps = append(deps, fmt.Sprint(f.Value))
				}
			}
		}
	}
	return deps
}

func GetFieldValue(fields []model.FieldType, fieldName string) interface{} {
	for _, f := range fields {
		if f.Name == fieldName {
			return f.Value
		}
	}
	return nil
}

func GetNames(model []model.FieldType) []string {
	var names []string
	for _, f := range model {
		if f.Name == "name" {
			names = append(names, fmt.Sprint(f.Value))
		}
	}
	return names
}

func GetName(model []model.FieldType) string {
	for _, f := range model {
		if f.Name == "name" {
			return fmt.Sprint(f.Value)
		}
	}

	return ""
}
