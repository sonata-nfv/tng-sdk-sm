#!/bin/bash

# Install setuptools
pip3 install setuptools || true

# Install son-mano-base package
root_dir=$(pwd)
cd son-mano-framework/son-mano-base
python3 setup.py develop

# Install specific manager package
cd $root_dir
cd base
python3 setup.py develop