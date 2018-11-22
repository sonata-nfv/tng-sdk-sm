// Copyright (c) 2015 SONATA-NFV, 2017 5GTANGO
// ALL RIGHTS RESERVED.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Neither the name of the SONATA-NFV, 5GTANGO
// nor the names of its contributors may be used to endorse or promote
// products derived from this software without specific prior written
// permission.

// This work has been performed in the framework of the SONATA project,
// funded by the European Commission under Grant number 671517 through
// the Horizon 2020 and 5G-PPP programmes. The authors would like to
// acknowledge the contributions of their colleagues of the SONATA
// partner consortium (www.sonata-nfv.eu).

// This work has been performed in the framework of the 5GTANGO project,
// funded by the European Commission under Grant number 761493 through
// the Horizon 2020 and 5G-PPP programmes. The authors would like to
// acknowledge the contributions of their colleagues of the 5GTANGO
// partner consortium (www.5gtango.eu).

package structs

// Cpu All the requirements and parameters related to the (virtual) CPU.
type Cpu struct {
  CpuClockSpeed string `yaml:"cpu_clock_speed,omitempty"`
  CpuModel string `yaml:"cpu_model,omitempty"`
  CpuSupportAccelerator string `yaml:"cpu_support_accelerator,omitempty"`
  Vcpus int `yaml:"vcpus"`
}

// Events 
type Events struct {
  Restart *Restart `yaml:"restart,omitempty"`
  ScaleIn *ScaleIn `yaml:"scale-in,omitempty"`
  ScaleOut *ScaleOut `yaml:"scale-out,omitempty"`
  Start *Start `yaml:"start,omitempty"`
  Stop *Stop `yaml:"stop,omitempty"`
}

// HypervisorParameters The requirements and parameters of a (potential) hyperviser that operates the VDU VM.
type HypervisorParameters struct {
  Type string `yaml:"type,omitempty"`
  Version string `yaml:"version,omitempty"`
}

// LifecycleEvents 
type LifecycleEvents struct {
  Authentication string `yaml:"authentication,omitempty"`
  AuthenticationType string `yaml:"authentication_type,omitempty"`
  AuthenticationUsername string `yaml:"authentication_username,omitempty"`
  Driver string `yaml:"driver,omitempty"`
  Events *Events `yaml:"events,omitempty"`
  FlavorIdRef string `yaml:"flavor_id_ref,omitempty"`
  VnfContainer string `yaml:"vnf_container,omitempty"`
}

// Memory 
type Memory struct {
  LargePagesRequired bool `yaml:"large_pages_required,omitempty"`
  NumaAllocationPolicy string `yaml:"numa_allocation_policy,omitempty"`
  Size float64 `yaml:"size"`
  SizeUnit string  `yaml:"size_unit,omitempty"`
}

// Monitoring 
type Monitoring struct {
  Command string `yaml:"command,omitempty"`
  Frequency float64 `yaml:"frequency,omitempty"`
  FrequencyUnit string `yaml:"frequency_unit,omitempty"`
  Name string `yaml:"name"`
  Unit string `yaml:"unit"`
}

// Network 
type Network struct {
  DataProcessingAccelerationLibrary string `yaml:"data_processing_acceleration_library,omitempty"`
  NetworkInterfaceBandwidth float64 `yaml:"network_interface_bandwidth,omitempty"`
  NetworkInterfaceBandwidthUnit string `yaml:"network_interface_bandwidth_unit,omitempty"`
  NetworkInterfaceCardCapabilities *NetworkInterfaceCardCapabilities `yaml:"network_interface_card_capabilities,omitempty"`
}

// NetworkInterfaceCardCapabilities Additional NIC capabilities:
type NetworkInterfaceCardCapabilities struct {
  Mirroring bool `yaml:"mirroring,omitempty"`
  SRIOV bool `yaml:"SR-IOV,omitempty"`
}

// Pcie The PCIe parameters of the platform.
type Pcie struct {
  DevicePassThrough bool `yaml:"device_pass_through,omitempty"`
  SRIOV bool `yaml:"SR-IOV,omitempty"`
}

// ResourceRequirementsements.
type ResourceRequirements struct {
  Cpu *Cpu `yaml:"cpu"`
  HypervisorParameters *HypervisorParameters `yaml:"hypervisor_parameters,omitempty"`
  Memory *Memory `yaml:"memory"`
  Network *Network `yaml:"network,omitempty"`
  Pcie *Pcie `yaml:"pcie,omitempty"`
  Storage *Storage `yaml:"storage,omitempty"`
  VswitchCapabilities *VswitchCapabilities `yaml:"vswitch_capabilities,omitempty"`
}

// Restart 
type Restart struct {
  Command string `yaml:"command,omitempty"`
  TemplateFile string `yaml:"template_file,omitempty"`
  TemplateFileFormat string `yaml:"template_file_format,omitempty"`
}

// ScaleIn 
type ScaleIn struct {
  Command string `yaml:"command,omitempty"`
  TemplateFile string `yaml:"template_file,omitempty"`
  TemplateFileFormat string `yaml:"template_file_format,omitempty"`
}

// ScaleInOut The scale-in/scale-out parameters.
type ScaleInOut struct {
  Maximum int `yaml:"maximum,omitempty"`
  Minimum int `yaml:"minimum,omitempty"`
}

// ScaleOut 
type ScaleOut struct {
  Command string `yaml:"command,omitempty"`
  TemplateFile string `yaml:"template_file,omitempty"`
  TemplateFileFormat string `yaml:"template_file_format,omitempty"`
}

// Start 
type Start struct {
  Command string `yaml:"command,omitempty"`
  TemplateFile string `yaml:"template_file,omitempty"`
  TemplateFileFormat string `yaml:"template_file_format,omitempty"`
}

// Stop 
type Stop struct {
  Command string `yaml:"command,omitempty"`
  TemplateFile string `yaml:"template_file,omitempty"`
  TemplateFileFormat string `yaml:"template_file_format,omitempty"`
}

// Storage 
type Storage struct {
  Persistence bool `yaml:"persistence,omitempty"`
  Size float64 `yaml:"size"`
  SizeUnit string `yaml:"size_unit,omitempty"`
}

// VirtualLinks 
type VirtualLinks struct {
  Access bool `yaml:"access,omitempty"`
  ConnectionPointsReference []string `yaml:"connection_points_reference"`
  ConnectivityType string `yaml:"connectivity_type"`
  Dhcp bool `yaml:"dhcp,omitempty"`
  ExternalAccess bool `yaml:"external_access,omitempty"`
  Id string `yaml:"id"`
  LeafRequirement string `yaml:"leaf_requirement,omitempty"`
  Qos string `yaml:"qos,omitempty"`
  RootRequirement string `yaml:"root_requirement,omitempty"`
}

// VswitchCapabilities 
type VswitchCapabilities struct {
  OverlayTunnel string `yaml:"overlay_tunnel,omitempty"`
  Type string `yaml:"type,omitempty"`
  Version string `yaml:"version,omitempty"`
}

//StartStop
type StartStop struct {
  Vnfr *Vnfr `yaml:vnfr`
  Vnfd *Vnfd `yaml:vnfd`
}

// Configure
type Configure struct {
  Vnfrs []*Vnfr `yaml:vnfrs`
  Vnfds []*Vnfd `yaml:vnfds`
}