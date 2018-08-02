#!/bin/bash
set -e

echo "-------------"
echo "Create an FSM"
echo "-------------"

export TNG_SM_PWD=`pwd`; ./go/bin/tng-sm new --type fsm test2

echo "---------------------"
echo "Delete the FSM exists"
echo "---------------------"

./go/bin/tng-sm delete test2-fsm

echo "-----------------------"
echo "Check if FSM is deleted"
echo "-----------------------"


if [ -d "test2-fsm" ]; then
	echo "FSM directory not deleted"
	exit 1
fi

echo "-------------"
echo "Create an SSM"
echo "-------------"

export TNG_SM_PWD=`pwd`; ./go/bin/tng-sm new --type ssm test2

echo "---------------------"
echo "Delete the SSM exists"
echo "---------------------"

./go/bin/tng-sm delete test2-ssm

echo "-----------------------"
echo "Check if SSM is deleted"
echo "-----------------------"


if [ -d "test2-ssm" ]; then
	echo "SSM directory not deleted"
	exit 1
fi
