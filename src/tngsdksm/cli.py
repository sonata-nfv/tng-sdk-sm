# Copyright (c) 2015 SONATA-NFV, 2017 5GTANGO
# ALL RIGHTS RESERVED.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Neither the name of the SONATA-NFV, 5GTANGO
# nor the names of its contributors may be used to endorse or promote
# products derived from this software without specific prior written
# permission.
#
# This work has been performed in the framework of the SONATA project,
# funded by the European Commission under Grant number 671517 through
# the Horizon 2020 and 5G-PPP programmes. The authors would like to
# acknowledge the contributions of their colleagues of the SONATA
# partner consortium (www.sonata-nfv.eu).
#
# This work has been performed in the framework of the 5GTANGO project,
# funded by the European Commission under Grant number 761493 through
# the Horizon 2020 and 5G-PPP programmes. The authors would like to
# acknowledge the contributions of their colleagues of the 5GTANGO
# partner consortium (www.5gtango.eu).

import sys
import argparse
import tngsdksm
import yaml
import os
import logging

LOG = logging.getLogger(__name__)

def run(args=None):
    """
    Entry point.
    Can get a list of args that are then used as input for the CLI arg parser.
    """

    if args is None:
        args = sys.argv[1:]
    parsed_args = parse_args(args)
    return dispatch(parsed_args)


def dispatch(args):
    """
    post process the arguments and link them to specific actions
    """

    # Handle the verbose argument
    init_logger(args.verbose)

    # abort if no subcommand is provided
    if args.subparser_name is None:
        print("Missing subcommand. Type tng-sm -h")
        exit(1)

    # new subcommand
    if args.subparser_name == 'new':
        res, mes = tngsdksm.create_specific_manager(args.type, args.name)
        if not res:
            print(mes)
        exit(not res)

    # generate subcommand
    if args.subparser_name == 'generate':
        if args.file.split('.')[-1] in ['yml', 'yaml']:
            res, mes = tngsdksm.generate_vnfr(yaml.load(open(args.file, 'r')))
            if not res:
                print(mes)
            exit(not res)
        elif args.file.split('.')[-1] in ['tgo']:
            res, mes = tngsdksm.generate_all(args.file)
            if not res:
                print(mes)
            exit(not res)
        else:
            print("Provided file is neither a vnfd nor a package.")
            exit(1)

    # execute subcommand
    if args.subparser_name == 'execute':
        if args.payload.split('.')[-1] not in ['yml', 'yaml']:
            print("Provided payload is not a yaml file.")
            exit(1)
        try:
          content = yaml.load(open(args.payload, 'r'))
        except:
            print("Couldn't open " + args.payload + ". Does it exist?")
            exit(1)
        args.sm = args.sm.strip('/')
        if args.sm.split('-')[-1] == 'fsm':
            res, mes = tngsdksm.execute_fsm(args.sm, args.event, content)
            if not res:
                print(mes)
            exit(not res)
        elif args.sm.split('-')[-1] == 'ssm':
            res, mes = tngsdksm.execute_ssm(args.sm, args.event, content)
            if not res:
                print(mes)
            exit(not res)
        else:
            print(args.sm + ' is not a valid SSM or FSM.')
            exit(1)

    return


def parse_args(args):
    """
    This method parses the arguments provided with the cli command
    """
    parser = argparse.ArgumentParser(description="5GTANGO tng-sm tool")

    parser.add_argument('-v',
                        '--verbose',
                        dest='verbose',
                        action="store_true",
                        default=False,
                        help='Verbose output')

    subparsers = parser.add_subparsers(description='',
                                       dest='subparser_name')

    parser_new = subparsers.add_parser('new',
                                       help='creating new specific managers')
    parser_gen = subparsers.add_parser('generate',
                                       help='generate payloads for testing')
    parser_exe = subparsers.add_parser('execute',
                                       help='execute testing events')

    # new sub arguments
    help_mes = 'Provide the type of specific manager, should be either.' \
               '\'ssm\' or \'fsm\'.'
    parser_new.add_argument('-t',
                              '--type',
                              required=True,
                              metavar="TYPE",
                              choices=['ssm', 'fsm'],
                              help=help_mes)

    parser_new.add_argument('-n',
                              '--name',
                              required=True,
                              metavar="NAME",
                              help='provide the name of specific manager')

    # generate sub arguments
    help_mes = 'Provide a file that serves as input for generated files. ' \
               'Should be either a yaml file with a vnfd, or a .tgo package.'
    parser_gen.add_argument('-f',
                              '--file',
                              required=True,
                              metavar="FILE",
                              help=help_mes)

    # execute sub arguments
    help_mes = 'The specific manager event that needs testing. One of ' \
               '[\'placement\', \'configure\', \'task\', \'state\'] for ssms, ' \
               '[\'start\', \'stop\', \'configure\', \'state\'] for fsms.'
    choices = ['placement', 'configure', 'task', 'start', 'stop', 'state']
    parser_exe.add_argument('-e',
                              '--event',
                              required=True,
                              metavar="EVENT",
                              choices=choices,
                              help=help_mes)

    help_mes = 'The payload that serves as input for the event that needs ' \
               'testing. Should be a yaml file.'
    parser_exe.add_argument('-p',
                              '--payload',
                              required=True,
                              metavar="PAYLOAD",
                              help=help_mes)

    help_mes = 'Specific manager subject to the test.'
    parser_exe.add_argument('-s',
                              '--sm',
                              required=True,
                              metavar="SPECIFIC MANAGER",
                              help=help_mes)

    return parser.parse_args(args)


def init_logger(verbose):
    """
    Configure logging
    """
    if verbose:
        level = logging.DEBUG
        logging.getLogger("tngsdksm").setLevel(level)
        logging.getLogger(__name__).setLevel(level)
        logging.basicConfig(level=level)

    return
