package model

type ProviderConfiguration struct {
	Uri string
}

type LibvirtVolume struct {
	Name           string `json:""`
	Pool           string //`default:"default"`
	Source         string //`default:""`
	Size           int
	BaseVolumeId   string
	BaseVolumeName string
	BaseVolumePool string
}

type LibvirtPool struct {
	Name string
	Type string
	Path string `default:"./"`
}

type LibvirtNetworkMode string

const (
	None   LibvirtNetworkMode = "none"
	Nat    LibvirtNetworkMode = "nat"
	Route  LibvirtNetworkMode = "route"
	Open   LibvirtNetworkMode = "open"
	Bridge LibvirtNetworkMode = "bridge"
)

type LibvirtDNS struct {
	LocalyOnly bool
	Forwarders []string
	SRVS       []string
	Hosts      []string
}

type LibvirtDHCP struct {
	Enabled bool
}

type LibvirtDNSMasqOptions struct {
	Options []LibvirtDNSmasqOption
}

type LibvirtDNSmasqOption struct {
	OptionName  string
	OptionValue string
}

type LibvirtNetwork struct {
	Name           string
	Domain         string
	Addresses      []string
	Mode           LibvirtNetworkMode
	Bridge         string
	MTU            int
	Autostart      bool
	Routes         []string
	DNS            LibvirtDNS
	DHCP           LibvirtDHCP
	DNSMasqOptions LibvirtDNSMasqOptions
	XSLT           string
}

type LibvirtCPU struct {
	Mode string
}

type LibvirtDomainDiskArguments struct {
	Blockdevice string
	File        string
	SCSI        bool
	URL         string
	VolumeId    string
	WWN         string
}

type LibvirtDomainNetworkInterfaceArguments struct {
	Addresses    []string
	Bridge       string
	Hostname     string
	Mac          string
	Macvtap      string
	NetworkId    string
	NetworkName  string
	Passthrough  string
	VEPA         string
	WaitForLease bool
}

type console struct {
	Type          string
	TargetPort    string
	TargetType    string
	SourcePath    string
	SourceHost    string
	SourceService string
}

type LibvirtDomainFilesystemArguments struct {
	Source     string
	Target     string
	Accessmode string
	ReadOnly   bool
}

type LibrivtDomainBootDeviceArguments struct {
	Devs []string
}

type LibvirtDomainTPMArguments struct {
	BackendDevicePath       string
	BackendEncryptionSecret string
	BackendPersistentState  bool
	BackendType             string
	BackendVersion          string
	Model                   string
}

type LibvirtDomain struct {
	Name             string
	Description      string
	CPU              LibvirtCPU
	VCPU             int
	Memory           int
	Running          bool
	Disk             LibvirtDomainDiskArguments
	NetworkInterface []LibvirtDomainNetworkInterfaceArguments
	Cloudinit        string
	Autostart        bool
	Filesystem       []LibvirtDomainFilesystemArguments
	Coreos_ignition  string
	Fw_cfg_Name      string //`default:"opt/com.coreos/config"`
	Arch             string
	Machine          string
	BootDevice       []LibrivtDomainBootDeviceArguments
	Emulator         string
	QEMU_Agent       bool
	TPM              LibvirtDomainTPMArguments
	Kernel           string
	Initrd           string
	Cmdline          []map[string]interface{}
}

type LibvirtIgnition struct {
	Name    string
	Pool    string
	Content string
}

type Cloudinit struct {
	Name          string
	Pool          string
	UserData      string
	MetaData      string
	NetworkConfig string
}
