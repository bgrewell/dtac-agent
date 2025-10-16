package structs

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type VirtualMachineID string

func (v *VirtualMachineID) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*v = VirtualMachineID(str)
		return nil
	}

	// Try unmarshaling as number
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*v = VirtualMachineID(strconv.Itoa(num))
		return nil
	}

	return fmt.Errorf("VirtualMachineID: unable to unmarshal %s", string(data))
}

// Machine is the struct for a machine
type Machine struct {
	AddressTTL                   string                   `json:"address_ttl" yaml:"address_ttl"`
	Architecture                 string                   `json:"architecture" yaml:"architecture"`
	BiosBootMethod               string                   `json:"bios_boot_method" yaml:"bios_boot_method"`
	CPUCount                     int                      `json:"cpu_count" yaml:"cpu_count"`
	CPUSpeed                     int                      `json:"cpu_speed" yaml:"cpu_speed"`
	CPUTestStatus                int                      `json:"cpu_test_status" yaml:"cpu_test_status"`
	CPUTestStatusName            string                   `json:"cpu_test_status_name" yaml:"cpu_test_status_name"`
	CommissioningStatus          int                      `json:"commissioning_status" yaml:"commissioning_status"`
	CommissioningStatusName      string                   `json:"commissioning_status_name" yaml:"commissioning_status_name"`
	CurrentCommissioningResultID int                      `json:"current_commissioning_result_id" yaml:"current_commissioning_result_id"`
	CurrentInstallationResultID  int                      `json:"current_installation_result_id" yaml:"current_installation_result_id"`
	CurrentTestingResultID       int                      `json:"current_testing_result_id" yaml:"current_testing_result_id"`
	Description                  string                   `json:"description" yaml:"description"`
	DistroSeries                 string                   `json:"distro_series" yaml:"distro_series"`
	DisableIpv4                  bool                     `json:"disable_ipv4" yaml:"disable_ipv4"`
	Fqdn                         string                   `json:"fqdn" yaml:"fqdn"`
	HardwareUUID                 string                   `json:"hardware_uuid" yaml:"hardware_uuid"`
	Hostname                     string                   `json:"hostname" yaml:"hostname"`
	HweKernel                    string                   `json:"hwe_kernel" yaml:"hwe_kernel"`
	InterfaceTestStatus          int                      `json:"interface_test_status" yaml:"interface_test_status"`
	InterfaceTestStatusName      string                   `json:"interface_test_status_name" yaml:"interface_test_status_name"`
	LastSync                     string                   `json:"last_sync" yaml:"last_sync"`
	Locked                       bool                     `json:"locked" yaml:"locked"`
	Memory                       int                      `json:"memory" yaml:"memory"`
	MemoryTestStatus             int                      `json:"memory_test_status" yaml:"memory_test_status"`
	MemoryTestStatusName         string                   `json:"memory_test_status_name" yaml:"memory_test_status_name"`
	MinHweKernel                 string                   `json:"min_hwe_kernel" yaml:"min_hwe_kernel"`
	Netboot                      bool                     `json:"netboot" yaml:"netboot"`
	NetworkTestStatus            int                      `json:"network_test_status" yaml:"network_test_status"`
	NetworkTestStatusName        string                   `json:"network_test_status_name" yaml:"network_test_status_name"`
	NodeType                     int                      `json:"node_type" yaml:"node_type"`
	NodeTypeName                 string                   `json:"node_type_name" yaml:"node_type_name"`
	NextSync                     string                   `json:"next_sync" yaml:"next_sync"`
	OtherTestStatus              int                      `json:"other_test_status" yaml:"other_test_status"`
	OtherTestStatusName          string                   `json:"other_test_status_name" yaml:"other_test_status_name"`
	Osystem                      string                   `json:"osystem" yaml:"osystem"`
	Owner                        string                   `json:"owner" yaml:"owner"`
	Pod                          string                   `json:"pod" yaml:"pod"`
	PowerState                   string                   `json:"power_state" yaml:"power_state"`
	PowerType                    string                   `json:"power_type" yaml:"power_type"`
	ResourceURI                  string                   `json:"resource_uri" yaml:"resource_uri"`
	Status                       int                      `json:"status" yaml:"status"`
	StatusAction                 string                   `json:"status_action" yaml:"status_action"`
	StatusMessage                string                   `json:"status_message" yaml:"status_message"`
	StatusName                   string                   `json:"status_name" yaml:"status_name"`
	Storage                      float64                  `json:"storage" yaml:"storage"`
	StorageTestStatus            int                      `json:"storage_test_status" yaml:"storage_test_status"`
	SwapSize                     string                   `json:"swap_size" yaml:"swap_size"`
	SyncInterval                 string                   `json:"sync_interval" yaml:"sync_interval"`
	SystemID                     string                   `json:"system_id" yaml:"system_id"`
	TestingStatus                int                      `json:"testing_status" yaml:"testing_status"`
	TestingStatusName            string                   `json:"testing_status_name" yaml:"testing_status_name"`
	VirtualmachineID             VirtualMachineID         `json:"virtualmachine_id" yaml:"virtualmachine_id"`
	InterfaceSet                 []InterfaceStruct        `json:"interface_set" yaml:"interface_set"`
	DefaultGateways              DefaultGatewayStruct     `json:"default_gateways" yaml:"default_gateways"`
	HardwareInfo                 map[string]interface{}   `json:"hardware_info" yaml:"hardware_info"`
	BootInterface                map[string]interface{}   `json:"boot_interface" yaml:"boot_interface"`
	Zone                         map[string]interface{}   `json:"zone" yaml:"zone"`
	BlockdeviceSet               []map[string]interface{} `json:"blockdevice_set" yaml:"blockdevice_set"`
	Bcaches                      []map[string]interface{} `json:"bcaches" yaml:"bcaches"`
	WorkloadAnnotations          map[string]interface{}   `json:"workload_annotations" yaml:"workload_annotations"`
	NumanodeSet                  []map[string]interface{} `json:"numanode_set" yaml:"numanode_set"`
	Raids                        []map[string]interface{} `json:"raids" yaml:"raids"`
	SpecialFilesystems           []map[string]interface{} `json:"special_filesystems" yaml:"special_filesystems"`
	CacheSets                    []map[string]interface{} `json:"cache_sets" yaml:"cache_sets"`
	Domain                       map[string]interface{}   `json:"domain" yaml:"domain"`
	Pool                         map[string]interface{}   `json:"pool" yaml:"pool"`
	VirtualblockdeviceSet        []map[string]interface{} `json:"virtualblockdevice_set" yaml:"virtualblockdevice_set"`
	OwnerData                    map[string]interface{}   `json:"owner_data" yaml:"owner_data"`
	PhysicalblockdeviceSet       []map[string]interface{} `json:"physicalblockdevice_set" yaml:"physicalblockdevice_set"`
	IPAddresses                  []string                 `json:"ip_addresses" yaml:"ip_addresses"`
	VolumeGroups                 []map[string]interface{} `json:"volume_groups" yaml:"volume_groups"`
	BootDisk                     map[string]interface{}   `json:"boot_disk" yaml:"boot_disk"`
	TagNames                     []string                 `json:"tag_names" yaml:"tag_names"`
}
