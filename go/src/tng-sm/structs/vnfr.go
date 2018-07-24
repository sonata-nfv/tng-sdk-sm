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


