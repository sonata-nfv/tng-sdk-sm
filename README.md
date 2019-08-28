[![Join the chat at https://gitter.im/sonata-nfv/Lobby](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/sonata-nfv/Lobby) [![Build Status](https://jenkins.sonata-nfv.eu/buildStatus/icon?job=tng-sdk-sm-pipeline/master)](https://jenkins.sonata-nfv.eu/job/tng-sdk-sm-pipeline/job/master/)

<p align="center"><img src="https://github.com/sonata-nfv/tng-api-gtw/wiki/images/sonata-5gtango-logo-500px.png" /></p>

# tng-sdk-sm

This repository contains the `tng-sdk-sm` component that is part of the European H2020 project [5GTANGO](http://www.5gtango.eu) NFV SDK. The component is responsible to create and test [5GTANGO specific managers](https://github.com/sonata-nfv/son-mano-framework/wiki/Service-and-Function-Specific-Managers). These specific managers are a concept introduced by the [5GTANGO MANO Framework](https://github.com/sonata-nfv/son-mano-framework). A specific manager is associated to a single network service (service specific manager (SSM)) or VNF (function specific manager (FSM)), and contains instructions which customise the MANO Framework behaviour when dealing with this network service or VNF. These instructions are defined as code, and need to be formatted in accordance with the MANO Framework API. A specific manager needs to consume the MANO Framework RabbitMQ message bus, as it is used for all MANO communications. Furthermore, since all MANO Framework components are Docker containers, a specific manager also needs to be packaged as a Docker container. `tng-sdk-sm` was developed to assist the network service and VNF developers in satisfying all these requirements.

It is advised to read the entire [5GTANGO specific managers](https://github.com/sonata-nfv/son-mano-framework/wiki/Service-and-Function-Specific-Managers) documentation before continuing with this tool. It will provide the reader with a clear understanding of why we use specific managers and how they work.

## Documentation

Besides this README file, more documentation is available in the [wiki](https://github.com/sonata-nfv/tng-sdk-sm/wiki) belonging to this repository. Additional information on specific managers are available at

* [5GTANGO specific managers](https://github.com/sonata-nfv/son-mano-framework/wiki/Service-and-Function-Specific-Managers)
* [Communications Magazine paper](https://ieeexplore.ieee.org/abstract/document/8713806)

## Installation and Dependencies

This component is implemented in Python3. Its requirements are specified [here](https://github.com/sonata-nfv/tng-sdk-sm/blob/master/requirements.txt).

### Automated

The automated installation requires `pip3`.

```bash
$ pip3 install git+https://github.com/sonata-nfv/tng-sdk-sm
```

### Manual

```bash
$ git clone https://github.com/sonata-nfv/tng-sdk-sm
$ cd tng-sdk-sm
$ python3 setup.py install
```

## Usage

The `tng-sm` command runs the `tng-sdk-sm` tool locally from the command line. Details about all possible parameters can be shown using:

```bash
tng-sm -h
```

More details about the subcommands can be found in the [wiki](https://github.com/sonata-nfv/tng-sdk-sm/wiki/Usage).

## Development

To contribute to the development of this 5GTANGO component, you may use the very same development workflow as for any other 5GTANGO Github project. That is, you have to fork the repository and create pull requests.

## License

This 5GTANGO component is published under Apache 2.0 license. Please see the LICENSE file for more details.

---
#### Lead Developers

The following lead developers are responsible for this repository and have admin rights. They can, for example, merge pull requests.

- Thomas Soenen ([@tsoenen](https://github.com/tsoenen))

#### Feedback-Channel

* Please use the GitHub issues to report bugs.
