package structs

// ConnectionPoints 
type VnfrConnectionPoints struct {
  Id string `yaml:"id"`
  Interface *VnfrInterface `yaml:"interface"`
  Type string `yaml:"type"`
  VirtualLinkReference string `yaml:"virtual_link_reference,omitempty"`
}

type VnfrInterface struct {
  Address string `yaml:address`
  HardwareAddress string `yaml:hardware_address`
  Netmask string `yaml:netmask`
}

// TheCoreSchemaForSONATANetworkFunctionRecords The core schema for SONATA network function records.
type Vnfr struct {
  ConnectionPoints []*VnfrConnectionPoints `yaml:"connection_points,omitempty"`
  DeploymentFavour string `yaml:"deployment_favour,omitempty"`
  DescriptorReference string `yaml:"descriptor_reference,omitempty"`
  DescriptorVersion string `yaml:"descriptor_version"`
  Id string `yaml:"id"`
  LifecycleEvents []*LifecycleEvents `yaml:"lifecycle_events,omitempty"`
  Logo string `yaml:"logo,omitempty"`
  ParentNs string `yaml:"parent_ns,omitempty"`
  Status string `yaml:"status"`
  Version string `yaml:"version,omitempty"`
  VirtualDeploymentUnits []*VnfrVirtualDeploymentUnits `yaml:"virtual_deployment_units"`
  VirtualLinks []*VirtualLinks `yaml:"virtual_links,omitempty"`
  VnfAddress []*VnfAddress `yaml:"vnf_address,omitempty"`
}

// VirtualDeploymentUnits 
type VnfrVirtualDeploymentUnits struct {
  Id string `yaml:"id"`
  MonitoringParameters []*Monitoring `yaml:"monitoring_parameters,omitempty"`
  NumberOfInstances int `yaml:"number_of_instances,omitempty"`
  ResourceRequirements *ResourceRequirements `yaml:"resource_requirements"`
  VduReference string `yaml:"vdu_reference,omitempty"`
  VmImage string `yaml:"vm_image,omitempty"`
  VnfcInstance []*VnfcInstance `yaml:"vnfc_instance,omitempty"`
}

// VnfAddress 
type VnfAddress struct {
  Address string `yaml:"address"`
}

// VnfcInstance 
type VnfcInstance struct {
  ConnectionPoints []*VnfrConnectionPoints `yaml:"connection_points,omitempty"`
  Id string `yaml:"id,omitempty"`
  VcId string `yaml:"vc_id,omitempty"`
  VimId string `yaml:"vim_id,omitempty"`
}


