"""
Copyright (c) 2015 SONATA-NFV, 2017 5GTANGO
ALL RIGHTS RESERVED.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Neither the name of the SONATA-NFV, 5GTANGO
nor the names of its contributors may be used to endorse or promote
products derived from this software without specific prior written
permission.

This work has been performed in the framework of the SONATA project,
funded by the European Commission under Grant number 671517 through
the Horizon 2020 and 5G-PPP programmes. The authors would like to
acknowledge the contributions of their colleagues of the SONATA
partner consortium (www.sonata-nfv.eu).

This work has been performed in the framework of the 5GTANGO project,
funded by the European Commission under Grant number 761493 through
the Horizon 2020 and 5G-PPP programmes. The authors would like to
acknowledge the contributions of their colleagues of the 5GTANGO
partner consortium (www.5gtango.eu).
"""
import logging
import yaml
import threading
import os
import re
from sonmanobase import messaging


logging.basicConfig(level=logging.INFO)
LOG = logging.getLogger("son-sm-base")
LOG.setLevel(logging.DEBUG)
logging.getLogger("son-mano-base:messaging").setLevel(logging.INFO)


class smbase(object):

    def __init__(self,
                 sm_name= None,
                 sm_version= None,
                 description=None):
        """
        :param specific_manager_type: specifies the type of specific manager that could be either fsm or ssm.
        :param service_name: the name of the service that this specific manager belongs to.
        :param function_name: the name of the function that this specific manager belongs to, will be null in SSM case
        :param specific_manager_name: the actual name of specific manager (e.g., scaling, placement)
        :param id_number: the specific manager id number which is used to distinguish between multiple SSM/FSM
        that are created for the same objective (e.g., scaling with algorithm 1 and 2)
        :param updated_version: specifies whether this SM is developed to update a current version or not,should be
        filled either by 'true' or 'false'
        :param version: version
        :param description: description
        """

        #Populating SSM-FSM fileds
        self.sm_name = sm_name
        self.sm_version = sm_version
        self.description = description
        self.uuid = None
        self.sfuuid = None

        LOG.info("Starting specific manager with name: " + self.sm_name)

        # # create and initialize broker connection
        # self.manoconn = messaging.ManoBrokerRequestResponseConnection(self.specific_manager_id)

        # self.wait_for_event = threading.Event()
        # self.wait_for_event.clear()

        # self.registration()

    def registration(self):

        """
        Send a register request to the Specific Manager registry.
        """
        LOG.info('Sending registration request...')

        if 'sf_uuid' in os.environ:
            self.sfuuid = os.environ['sf_uuid']
        else:
            self.sfuuid = ''

        message = {'specific_manager_type': self.specific_manager_type,
                   'service_name': self.service_name,
                   'function_name': self.function_name,
                   'specific_manager_name': self.specific_manager_name,
                   'specific_manager_id': self.specific_manager_id,
                   'update_version': self.update_version,
                   'version': self.version,
                   'description': self.description,
                   "sf_uuid": self.sfuuid}

        self.manoconn.call_async(self._on_registration_response,
                                 'specific.manager.registry.ssm.registration',
                                 yaml.dump(message))

        self.waitForRegistration()

    def _on_registration_response(self, ch, method, props, response):

        response = yaml.load(response)
        if response['status'] != "registered":
            LOG.error("{0} registration failed. Exit".format(self.specific_manager_id))
        else:
            self.uuid = response['uuid']
            LOG.info("{0} registered with uuid:{1}".format(self.specific_manager_id, self.uuid))

            # release the registration thread
            self.wait_for_event.set()

            # jump to on_registration_ok()
            self.on_registration_ok()

    def waitForRegistration(self):
        if not self.wait_for_event.wait(20):
            LOG.error("Registration response not received.")

    def on_registration_ok(self):

        """
        To be overwritten by subclasses
        """
        LOG.info("Received registration ok event.")