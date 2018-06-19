# tng-sdk-sm

5GTANGO SDK Specific Manager develop and test framework. The 5GTANGO SP supports
Service (SSM) and Function (FSM) Specific Managers. These managers are created by
Network Service and VNF developers to customise the MANO functionality of the SP.
This CLI-tool aims to help develop and test such specific managers. 

## Clone

This repository uses submodules. Clone it recursively to get all required content.

	git clone --recurse-submodules https://github.com/sonata-nfv/tng-sdk-sm.git

## Dependencies

* python3.x
* Docker
* golang

The specific managers for which we provide development and testing support in this 
repository are developed with python3.

The specific managers are embedded in docker containers. See https://docs.docker.com/install/
to install Docker CE.

The CLI tool is build using golang. See https://golang.org/doc/install to install.

## Installation

Set the GOPATH ENV parameter from the root of this repository

	export GOPATH=$(pwd)/go

Get the golang dependencies

	go get gopkg.in/yaml.v2
	go get github.com/nu7hatch/gouuid
	go get github.com/fatih/color

Build the project

	go build tng-sdk-sm

Install the project

	go install tng-sdk-sm

In order to test the specific managers natively, your python environment needs to be
configured. Execute

	./env_config.sh

The CLI tool needs to know where to find your version of this repository. Therefore,
on the root of this repository, set the following ENV parameter:

	export TNG_SM_PWD=$(pwd)

## Usage

	
	usage: tng-sm [--version] [--help]

	These are the subcommands for tng-sm:

	    new            Create a new specific manager
	    delete         Delete an existing specific manager
	    execute        Execute an event of a specific manager
	    generate       Generate artefacts to be used when executing specific managers

	usage: tng-sm new <specific manager name>

	    --path         Path where new specific manager should be stored
	    --type         Type of specific manager to be created: "ssm" or "fsm"

	usage: tng-sm delete <specific manager name>

	    --path         Path where specific manager can be found

	usage: tng-sm execute <specific manager name>

	    --path         Path where specific manager can be found
	    --event        Event that needs to be executed: "start", "stop" or "configure"
	    --payload      Payload for the execution

	usage: tng-sm generate <name output file>

	    --type         Type of payload to be generated: "vnfr" or "nsr"
	    --descriptor   File that serves as input for generation, should be a vnfd or nsd
    	

## Example

To create a new FSM, run

	tng-sm new --type fsm foo

This creates the directory `foo-fsm`. The actual FSM code needs to be added at `/foo-fsm/foo/foo.py`.
The `start_event` of an FSM is executed just after VNF instantiation. At this stage, it is often
required to SSH into the VNF and perform some config. An example of a `start_event` method could be

	
    def start_event(self, content):
	    """
	    This method handles a start event.
	    """
	    # Extract the mgmt ip of the VNF from the VNFR
	    vnfc = content['vnfr']['virtual_deployment_units'][0]['vnfc_instance'][0]
	    mgmt_ip = vnfc['connection_points'][0]['interface']['address']
	    # Initiate SSH connection with the VM
	    ssh_client = Client(mgmt_ip,username='foo',password='bar')
	    # Execute command on remote
	    ssh_client.sendCommand("mkdir foo")
	    # Dummy content
	    response = {'status': 'completed'}
	    return response
	

To test this `start_event` method, the `content` input parameter is required as if the FSM will receive it
at runtime. It can be generated based on the VNFD of the VNF with

	tng-sm generate --type vnfr --descriptor foobar.yml output.yml

`/examples/generated/output.yml` was generated based on the descriptor at `/examples/vnfd/foobar.yml`. Some fields in the generated file 
need to be extended with runtime information of the VM that will be used to test the FSM. These fields are 
marked as `< ... >`.

Next, the `start_event` method can be tested by running

	tng-sm execute --event start --payload output.yml foo-fsm

The CLI will return the output of the event.

To delete this FSM, use

	tng-sm delete foo-fsm

