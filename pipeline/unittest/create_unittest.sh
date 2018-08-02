#!/bin/bash
set -e

echo "-------------"
echo "Create an FSM"
echo "-------------"

export TNG_SM_PWD=`pwd`; ./go/bin/tng-sm new --type fsm test

echo "----------------------"
echo "Evaluate if FSM exists"
echo "----------------------"

if [ ! -d "test-fsm" ]; then
	echo "FSM directory not existing"
	exit 1
fi

echo "--------------------------"
echo "Evaluate if FSM is correct"
echo "--------------------------"

if [ ! -e "test-fsm/test/ssh.py" ]; then
	echo "FSM not correct"
	exit 1
fi

echo "-------------"
echo "Create an SSM"
echo "-------------"

export TNG_SM_PWD=`pwd`; ./go/bin/tng-sm new --type ssm test

echo "----------------------"
echo "Evaluate if SSM exists"
echo "----------------------"

if [ ! -d "test-ssm" ]; then
	echo "SSM directory not existing"
	exit 1
fi

echo "--------------------------"
echo "Evaluate if SSM is correct"
echo "--------------------------"

if [ -e "test-ssm/test/ssh.py" ]; then
	echo "SSM not correct"
	exit 1
fi
