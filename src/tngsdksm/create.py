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

import shutil
import os
import logging

LOG = logging.getLogger(__name__)

def create_specific_manager(sm_type, name, dest_path=os.getcwd()):
    """Create a new specific manager

    :param sm_type: A string, fsm or ssm
    :param name: A string, the name of the specific manager
    :param dest_path: A string, the target directory for the manager

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a string with a possible error message
    """
    if sm_type not in ['fsm', 'ssm']:
        return False, 'sm_type should be either \'fsm\' or \'ssm\''

    mod_path = os.path.abspath(os.path.join(__file__,os.pardir))
    res_path = os.path.join(dest_path, name + '-' + sm_type)

    if os.path.isdir(res_path):
        return False, 'Specific manager already exists'

    shutil.copytree(os.path.join(mod_path,'templates'), res_path)

    os.rename(os.path.join(res_path, 'template'), os.path.join(res_path, name))

    if sm_type == 'fsm':
        os.rename(os.path.join(res_path, name, 'fsm_template.py'),
                  os.path.join(res_path, name, name + '.py'))
        os.remove(os.path.join(res_path, name, 'ssm_template.py'))

    if sm_type == 'ssm':
        os.rename(os.path.join(res_path, name, 'ssm_template.py'),
                  os.path.join(res_path, name, name + '.py'))
        os.remove(os.path.join(res_path, name, 'fsm_template.py'))

    for dirName, subdirList, fileList in os.walk(os.path.join(res_path)):
        if 'pycache' in dirName:
            shutil.rmtree(dirName)
            continue
        for file in fileList:
            file_in = os.path.join(dirName, file)
            file_out = os.path.join(dirName, "tmp.txt")
            with open(file_in, "r") as fin:
                with open(file_out, "w") as fout:
                    for line in fin:
                        nl = line.replace("<name>", name).replace("<type>",
                                                                  sm_type)
                        fout.write(nl)
            os.rename(file_out, file_in)

    return True, ''


