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

import os
import sys
import importlib
import traceback
import logging

LOG = logging.getLogger(__name__)

def execute_fsm(fsm_path, event, payload):
    """Execute an FSM

    :param fsm_path: A string, path to the fsm
    :param event: A string, the event to be tested (either start, stop,
        configure or state)
    :param payload: A dictionary, the input date for the method under testing

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a string with a possible error message
    """

    #check whether fsm exists
    if not os.path.isdir(fsm_path):
        return False, 'FSM does not exist with name \'' + fsm_path + '\'.'

    sup_ev = ['start', 'stop', 'configure', 'state']
    if event not in sup_ev:
        msg = '\'' + event + '\' event is not supported for fsms ' \
              '(' + str(sup_ev) + ').'
        return False, msg

    # extract fsm code module
    fsm_name = fsm_path.split('/')[-1]
    if '-' not in fsm_name:
        return False, 'FSM name was changed, cannot locate module'

    sys.stdout.write("\033[0;32m")
    print('\n#######output from test event#######')
    res, mes = _execute_sm(fsm_path, event, payload)
    print('####################################\n')
    sys.stdout.flush()
    sys.stdout.write("\033[0;0m")
    return res, mes


def execute_ssm(ssm_path, event, payload):
    """Execute an SSM

    :param ssm_path: A string, path to the ssm
    :param event: A string, the event to be tested (either task, configure,
        place or state)
    :param payload: A dictionary, the input date for the method under testing

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a string with a possible error message
    """

    #check whether fsm exists
    if not os.path.isdir(ssm_path):
        return False, 'SSM does not exist with name \'' + ssm_path + '\'.'

    sup_ev = ['task', 'placement', 'configure', 'state']
    if event not in sup_ev:
        msg = '\'' + event + '\' event is not supported for ssms ' \
              '(' + str(sup_ev) + ').'
        return False, msg

    # extract ssm code module
    ssm_name = ssm_path.split('/')[-1]
    if '-' not in ssm_name:
        return False, 'SSM name was changed, cannot locate module'

    sys.stdout.write("\033[0;32m")
    print('\n#######output from test event#######')
    res, mes = _execute_sm(ssm_path, event, payload)
    print('####################################\n')
    sys.stdout.flush()
    sys.stdout.write("\033[0;0m")
    return res, mes


def _execute_sm(path, event, payload, color_traceback=True):

    sm_name = path.split('/')[-1]
    module_name = sm_name.split('-')[0]
    module_dir = os.path.join(path, module_name)
    sm_type = sm_name.split('-')[-1].upper()

    # Execute the specific manager
    try:
        sys.path.insert(1, module_dir)
        module = importlib.import_module(module_name)
        obj = getattr(module, module_name + sm_type)
        sm = obj(connect_to_broker=False)
        met = getattr(sm, event + '_event')
        met(payload)
    except Exception as e:
        if color_traceback:
            sys.stdout.write("\033[1;31m")
            sys.stdout.flush()
        traceback.print_exc()

    return True, ''
