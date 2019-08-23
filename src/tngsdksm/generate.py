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

import uuid
import yaml
import shutil
import os
import logging
import tngsdk.package as tngpkg

LOG = logging.getLogger(__name__)

def generate_vnfr(vnfd, name='vnfr', dest_path=os.getcwd()):
    """Generate a vnfr based on a vnfd

    :param vnfd: A dictionary, the vnfd
    :param name: A string, the filename of the resulting vnfr
    :param dest_path: A string, the target directory for the vnfr

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a dictionary with the vnfr
    """

    vnfr = {}
    vnfr['id'] = str(uuid.uuid4())
    vnfr['virtual_links'] = vnfd.get('virtual_links')
    vnfr['descriptor_version'] = vnfd.get('descriptor_version')
    vnfr['descriptor_reference'] = str(uuid.uuid4())
    vnfr["status"] = "normal operation"
    vnfr["version"] = "1"

    if 'virtual_deployment_units' in vnfd.keys():
        vnfr['virtual_deployment_units'] = []
        vdus = vnfr['virtual_deployment_units']

        for d_vdu in vnfd['virtual_deployment_units']:
            r_vdu = {}
            r_vdu['id'] = d_vdu.get('id')
            r_vdu['number_of_instances'] = 1
            r_vdu['resource_requirements'] = d_vdu.get('resource_requirements')
            r_vdu['vdu_reference'] = d_vdu.get('vdu_reference')
            r_vdu['vm_image'] = d_vdu.get('vm_image')
            vnfc = {}
            vnfc['id'] = "0"
            vnfc['vim_id'] = str(uuid.uuid4())
            r_cps = []
            for d_cp in d_vdu['connection_points']:
                r_cp = {}
                r_cp['type'] = d_cp.get('type')
                r_cp['id'] = d_cp.get('id')
                r_cp['interface'] = {}
                r_cp['interface']['address'] = '<ENTER IP ADDRESS>'
                r_cp['interface']['hardwareaddress'] = '<ENTER MAC ADDRESS>'
                r_cp['interface']['netmask'] = '<ENTER NETMASK>'
                r_cps.append(r_cp)
            vnfc['connection_points'] = r_cps
            r_vdu['vnfc_instance'] = [vnfc]
            vdus.append(r_vdu)

    if 'cloudnative_deployment_units' in vnfd.keys():
        vnfr['cloudnative_deployment_units'] = []
        cdus = vnfr['cloudnative_deployment_units']

        for d_cdu in vnfd['cloudnative_deployment_units']:
            r_cdu = {}
            r_cdu['id'] = d_cdu.get('id')
            r_cdu['image'] = d_cdu.get('image')
            r_cdu['number_of_instances'] = 1
            r_cdu['vim_id'] = str(uuid.uuid4())
            r_cdu['load_balancer_ip'] = {'floating_ip': '<ENTER IP HERE>',
                                         'internal_ip': '<ENTER IP HERE>'}
            r_cps = []
            for d_cp in d_cdu['connection_points']:
                r_cp = d_cp
                r_cp['type'] = 'serviceendpoint'
                r_cps.append(r_cp)
            r_cdu['connection_points'] = r_cps
            cdus.append(r_cdu)           

    with open(os.path.join(dest_path, name + '.yml'), 'w') as fout:
        yaml.dump(vnfr, fout, default_flow_style=False)

    return True, vnfr

def generate_nsr(nsd, vnfrs, name='nsr', dest_path=os.getcwd()):
    """Generate an nsr based on a vnfd

    :param nsd: A dictionary, the nsd
    :param vnfrs: a list of dictionaries, the vnfrs
    :param name: A string, the filename of the resulting nsr
    :param dest_path: A string, the target directory for the nsr

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a dictionary with the nsr
    """

    nsr = {}
    nsr['descriptor_reference'] = str(uuid.uuid4())
    nsr['descriptor_version'] = nsd.get('descriptor_version')
    nsr['instance_name'] = 'foobar'
    nsr['network_functions'] = []
    for vnfr in vnfrs:
        nsr['network_functions'].append({'vnfr_id': vnfr['id']})
    nsr['status'] = 'normal operation'
    nsr['id'] = str(uuid.uuid4())
    nsr['version'] = '1'
    nsr['virtual_links'] = nsd['virtual_links']

    with open(os.path.join(dest_path, name + '.yml'), 'w') as fout:
        yaml.dump(nsr, fout, default_flow_style=False)

    return True, nsr

def generate_all(pkg_path, dest_path=os.getcwd()):
    """Generate a set of nsr and vnfrs based on a service package

    :param pkg_path: A string, a path to the package
    :param dest_path: A string, the target directory for the set of files

    :returns: a tuple with two elements. [0] is a bool with the result, 
        [1] a string with a possible error message
    """
    # Unpackage the package
    pkg_proj = _unpack(pkg_path, '_tmp')
    nsds_f, vnfds_f = _obtain_descriptor_paths(pkg_proj)

    # Open artefacts
    nsd = yaml.load(open(nsds_f[0], 'r'))

    vnfds = []
    for vnfd_f in vnfds_f:
        vnfds.append(yaml.load(open(vnfd_f, 'r')))

    # Create resulting directory of files
    dir_path = os.path.join(dest_path, nsd['name'])
    if os.path.isdir(dir_path):
        shutil.rmtree('_tmp', ignore_errors=True)
        return False, 'directory already exists'
    os.makedirs(dir_path)

    vnfrs = []
    for vnfd in vnfds:
        res, vnfr = generate_vnfr(vnfd,
                                  'vnfr' + str(vnfds.index(vnfd)),
                                  dir_path)
        vnfrs.append(vnfr)

    generate_nsr(nsd, vnfrs, name='nsr', dest_path=dir_path)

    shutil.rmtree('_tmp', ignore_errors=True)

    return True, ''

def _unpack(pkg_path, proj_path):
    """
    Wraps the tng-sdk-package unpacking functionality.
    """
    args = [
        "--unpackage", pkg_path,
        "--output", proj_path,
        "--store-backend", "TangoProjectFilesystemBackend",
        "--format", "eu.5gtango",
        "--quiet",
        "--offline",
        "--skip-validation",
        "--loglevel", "ERROR"
    ]
    r = tngpkg.run(args)
    if r.error is not None:
        raise BaseException("Can't read package {}: {}"
                            .format(pkg_path, r.error))
    # return the full path to the project
    proj_path = r.metadata.get("_storage_location")
#    LOG.debug("Unpacked {} to {}".format(pkg_path, proj_path))
    return proj_path

def _obtain_descriptor_paths(proj_path):
    """
    Return descriptor paths for unpackaged package folder
    """

    proj = yaml.load(open(proj_path + '/project.yml', 'r'))

    nsds = []
    vnfds = []

    for file in proj['files']:
        if 'eu.5gtango' in file['tags']:
            if file['type'] == 'application/vnd.5gtango.nsd':
                nsds.append(proj_path + '/' + file['path'])
            elif file['type'] == 'application/vnd.5gtango.vnfd':
                vnfds.append(proj_path + '/' + file['path'])
            else:
                pass

    return nsds, vnfds
