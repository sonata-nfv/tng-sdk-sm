# tng-sdk-sm

5GTANGO SDK Specific Manager develop and test framework

## Clone

This repository uses submodules. Clone it recursively to get all required content.

	`git clone --recurse-submodules https://github.com/sonata-nfv/tng-sdk-sm.git`

## Dependencies

* python3.x
* Docker
* golang

The specific managers for which we provide development and testing support in 
this repository are developed with python3.

The specific managers are embedded in docker containers. See https://docs.docker.com/install/
to install Docker CE.

The CLI tool is build using golang. See https://golang.org/doc/install to install.

## Installation

In order to use the CLI tool, first build to golang project

	`go build tng-sdk-sm`

and install it.

	`go install tng-sdk-sm`

In order to test the specific managers natively, your python env needs to be
configured. Execute

	`./env_config.sh`

The CLI tool needs to know where to find your version of this repository. Therefore,
on the root of this repository, set the following ENV paramater:

	`export TNG_SM_PWD=$pwd`

## Usage