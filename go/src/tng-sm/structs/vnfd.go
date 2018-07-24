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

// AssuranceParameters 
type AssuranceParameters struct {
  Formula string `yaml:"formula,omitempty"`
  Id string `yaml:"id,omitempty"`
  Penalty *Penalty `yaml:"penalty,omitempty"`
  RelId string `yaml:"rel_id,omitempty"`
  Unit string `yaml:"unit,omitempty"`
  Value int `yaml:"value,omitempty"`
  Violation []*Violation `yaml:"violation,omitempty"`
}

// DeploymentFlavours 
type DeploymentFlavours struct {
  AssuranceParameters []*AssuranceParameters `yaml:"assurance_parameters,omitempty"`
  Constraint string `yaml:"constraint,omitempty"`
  FlavourKey string `yaml:"flavour_key,omitempty"`
  Id string `yaml:"id,omitempty"`
  VduReference []string `yaml:"vdu_reference,omitempty"`
  VlinkReference []string `yaml:"vlink_reference,omitempty"`
}


// FunctionSpecificManagers An FSM object of this VNF. FSMs are always Docker containers.
type FunctionSpecificManagers struct {
  Description string `yaml:"description,omitempty"`
  Id string `yaml:"id"`
  Image string `yaml:"image"`
  ImageMd5 string `yaml:"image_md5,omitempty"`
  Options []*Options `yaml:"options,omitempty"`
  ResourceRequirements *ResourceRequirements `yaml:"resource_requirements,omitempty"`
}

// MonitoringRules 
type MonitoringRules struct {
  Condition string `yaml:"condition"`
  Description string `yaml:"description,omitempty"`
  Duration float64 `yaml:"duration"`
  DurationUnit string `yaml:"duration_unit,omitempty"`
  Name string `yaml:"name"`
  Notification []*Notification `yaml:"notification"`
}

// Notification 
type Notification struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
}

// Options A key-value parameter object.
type Options struct {
  Key string `yaml:"key"`
  Value string `yaml:"value"`
}

// Penalty 
type Penalty struct {
  Expression int `yaml:"expression,omitempty"`
  Type string `yaml:"type,omitempty"`
  Unit string `yaml:"unit,omitempty"`
  Validity string `yaml:"validity,omitempty"`
}

// Violation 
type Violation struct {
  BreachesCount int `yaml:"breaches_count,omitempty"`
  Interval int `yaml:"interval,omitempty"`
}

// VirtualDeploymentUnits 
type VnfdVirtualDeploymentUnits struct {
  ConnectionPoints []*VnfdConnectionPoints `yaml:"connection_points,omitempty"`
  Description string `yaml:"description,omitempty"`
  Id string `yaml:"id"`
  MonitoringParameters []*Monitoring `yaml:"monitoring_parameters,omitempty"`
  ResourceRequirements *ResourceRequirements `yaml:"resource_requirements"`
  ScaleInOut *ScaleInOut `yaml:"scale_in_out,omitempty"`
  UserData string `yaml:"user_data,omitempty"`
  VmImage string `yaml:"vm_image,omitempty"`
  VmImageFormat string  `yaml:"vm_image_format,omitempty"`
  VmImageMd5 string `yaml:"vm_image_md5,omitempty"`
}

// VNFD The core schema for SONATA network function descriptors.
type Vnfd struct {
  Author string `yaml:author`
  ConnectionPoints []*VnfdConnectionPoints `yaml: "connection_points`
  DeploymentFlavours []*DeploymentFlavours `yaml: "deployment_flavours`
  Description string `yaml:description`
  DescriptorVersion string `yaml:descriptor_version`
  FunctionSpecificManagers []*FunctionSpecificManagers `yaml:"function_specific_managers,omitempty"`
  Licenses []string `yaml:"licenses,omitempty"`
  LifecycleEvents []*LifecycleEvents `yaml:"lifecycle_events,omitempty"`
  Logo string `yaml:"logo,omitempty"`
  MonitoringRules []*MonitoringRules `yaml:"monitoring_rules,omitempty"`
  Name string `yaml:"name"`
  Vendor string `yaml:"vendor"`
  Version string `yaml:"version"`
  VirtualDeploymentUnits []*VnfdVirtualDeploymentUnits `yaml:"virtual_deployment_units"`
  VirtualLinks []*VirtualLinks `yaml:"virtual_links,omitempty"`
}

// ConnectionPoints 
type VnfdConnectionPoints struct {
  Id string `yaml:id`
  Interface string `yaml:interface`
  Type string `yaml:type`
  VirtualLinkReference string
}
