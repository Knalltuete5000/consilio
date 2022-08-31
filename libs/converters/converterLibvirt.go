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

	helper "github.com/kevinklinger/consilio/libs"
	"github.com/kevinklinger/consilio/model"
)

type LibvirtConfig struct {
	Domain    LibvirtDomainConfig // for a typical terraform provider this should be arrays instead of a single object
	Pool      LibvirtPoolConfig
	Network   LibvirtNetworkConfig
	Volume    LibvirtVolumeConfig
	CloudInit LibvirtCloudInitConfig
	Ignition  LibvirtIgnitionConfig
}

type LibvirtDomainConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	Domain         model.DynamicElement
}

type LibvirtPoolConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	Pool           model.DynamicElement
}

type LibvirtNetworkConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	Network        model.DynamicElement
}

type LibvirtVolumeConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	Volume         model.DynamicElement
}

type LibvirtCloudInitConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	CloudInit      model.DynamicElement
}

type LibvirtIgnitionConfig struct {
	Name           string
	CreateVariable bool
	DependsOn      []string
	Ignition       model.DynamicElement
}

func ConvertToLibvirtConfig(dynamicElements []model.DynamicElement) (LibvirtConfig, error) {
	var configs LibvirtConfig
	for _, model := range dynamicElements {
		switch model.Name {
		case "libvirt_cloudinit_disk":
			configs.CloudInit.CloudInit = model
			configs.CloudInit.DependsOn = checkDependenciesForCloudInitConfig(model.Fields, dynamicElements)
			configs.CloudInit.Name = GetName(model.Fields)
			configs.CloudInit.CreateVariable = false
		case "libvirt_domain":
			configs.Domain.Domain = model
			configs.Domain.DependsOn = checkDependenciesForDomainConfig(model.Fields, dynamicElements)
			configs.Domain.Name = GetName(model.Fields)
			configs.Domain.CreateVariable = false // no other resource can depend on the domain
		case "libvirt_ignition":
			configs.Ignition.Ignition = model
			configs.Ignition.DependsOn = checkDependenciesForIgnitionConfig(model.Fields, dynamicElements)
			configs.Ignition.Name = GetName(model.Fields)
			configs.Ignition.CreateVariable = false
		case "libvirt_network":
			configs.Network.Network = model
			configs.Network.DependsOn = checkDependenciesForNetworkConfig(model.Fields, dynamicElements)
			configs.Network.Name = GetName(model.Fields)
			configs.Network.CreateVariable = false
		case "libvirt_pool":
			configs.Pool.Pool = model
			configs.Pool.DependsOn = checkDependenciesForPoolConfig(model.Fields, dynamicElements)
			configs.Pool.Name = GetName(model.Fields)
			configs.Pool.CreateVariable = false
		case "libvirt_volume":
			configs.Volume.Volume = model
			configs.Volume.DependsOn = checkDependenciesForVolumeConfig(model.Fields, dynamicElements)
			configs.Volume.Name = GetName(model.Fields)
			configs.Volume.CreateVariable = false
		default:
		}
	}

	er := checkCyclicDependencies(configs)
	if er != nil {
		return configs, er
	}

	checkForVariableCreation(&configs)

	return configs, nil
}

func checkCyclicDependencies(config LibvirtConfig) error {
	// volume does not refer to itself
	// FIXME We only check on the name of the volume
	// If a pool has the same name as the volume, it is recognized as a cyclice dependency
	var name string = GetName(config.Volume.Volume.Fields)
	for _, dep := range config.Volume.DependsOn {
		if dep == name {
			return errors.New("volume depends on itself")
		}
	}

	return nil
}

func checkForVariableCreation(config *LibvirtConfig) {
	// Pool
	var pool_deps []string
	pool_deps = append(pool_deps, config.Volume.DependsOn...)
	pool_deps = append(pool_deps, config.CloudInit.DependsOn...)
	pool_deps = append(pool_deps, config.Ignition.DependsOn...)
	if helper.Contains(config.Pool.Name, pool_deps) {
		config.Pool.CreateVariable = true
	}

	// Volume
	var volume_deps []string
	volume_deps = append(volume_deps, config.Domain.DependsOn...)
	volume_deps = append(volume_deps, config.Volume.DependsOn...) // in our case this could not happen because of cyclic dependencies but if we had multiple volumes this can be possible
	if helper.Contains(config.Volume.Name, volume_deps) {
		config.Volume.CreateVariable = true
	}

	// CloudInit
	var cloud_init_deps []string
	cloud_init_deps = append(cloud_init_deps, config.Domain.DependsOn...)
	if helper.Contains(config.CloudInit.Name, cloud_init_deps) {
		config.CloudInit.CreateVariable = true
	}

	// Ignition
	var ignition_deps []string
	ignition_deps = append(ignition_deps, config.Domain.DependsOn...)
	if helper.Contains(config.Ignition.Name, ignition_deps) {
		config.Ignition.CreateVariable = true
	}

	// Network
	var network_deps []string
	network_deps = append(network_deps, config.Domain.DependsOn...)
	if helper.Contains(config.Network.Name, network_deps) {
		config.Network.CreateVariable = true
	}
}

func checkDependenciesForCloudInitConfig(cloudinit []model.FieldType, models []model.DynamicElement) []string {
	var deps []string
	// check for pool (field pool)
	// get pool
	var pool_name = fmt.Sprint(GetFieldValue(cloudinit, "pool"))

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
