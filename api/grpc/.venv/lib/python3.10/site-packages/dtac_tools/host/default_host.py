# default_host.py
import os
import sys
import json
import grpc
import traceback
import dtac_tools.host.plugin_pb2
import dtac_tools.host.plugin_pb2_grpc

from grpc import ssl_server_credentials
from urllib.parse import quote
from concurrent import futures
from .helpers.debug_sender import DebugSender
from .helpers.encryptor import RpcEncryptor
from .plugin_host import PluginHost
from .helpers.network import get_unused_tcp_port
from ..plugins.types import InputArgs
from .plugin_pb2 import RegisterArgs, RegisterReply, PluginEndpoint


class DefaultPluginHost(PluginHost):
    def __init__(self, plugin, debug=False, debug_port=5678):
        self.plugin = plugin
        self.rpc_proto = "grpc"
        self.proto = "tcp"
        self.ip = "127.0.0.1"
        self.interface_version = 'plug_api_1.0'
        self.port = None
        self.route_map = {}
        self.encryptor = RpcEncryptor.new_encryptor()
        if debug:
            self.debug_sender = DebugSender(debug_port)
            self.debug(f"PluginHost debugging enabled on port {debug_port}")
            self.debug(f"Plugin Name: {self.plugin.name()}")

    def debug(self, message):
        if self.debug_sender is not None:
            self.debug_sender.write(message + "\n")

    def Register(self, request, context):
        self.debug("Register called but not implemented")
        try :
            params = {}
            if request.config is not None:
                self.debug(f"config: {request.config}")
                params["config"] = request.config

            if request.default_secure is not None:
                self.debug(f"default_secure: {request.default_secure}")
                params["default_secure"] = request.default_secure

            response = self.plugin.register(params)

            return_eps = []
            for endpoint in response:
                api_ep = PluginEndpoint(
                    path = endpoint.path,
                    action = endpoint.action,
                    uses_auth = endpoint.uses_auth,
                    expected_args = "",
                    expected_body = "",
                    expected_output = "",
                )
                return_eps.append(api_ep)

            self.debug(f"response: {response}")
            return_val = RegisterReply(
                endpoints = return_eps,
            )
            return return_val
        except Exception as ex:
            self.debug(f"Exception: {ex}")
            self.debug(traceback.format_exc())


    def Call(self, request, context):
        self.debug("Call called but not implemented")

    def handle_request(self, data):
        try:
            request = json.loads(data)
            self.debug(f"request: {request}")
            if request["method"] == f"{self.plugin.name()}.Register":
                self.debug(f"Calling {request['method']} - Params: {request['params']}")
                endpoints = self.plugin.register(request["params"])
                self.debug(f"endpoints: {endpoints}")
                for endpoint in endpoints:
                    self.route_map[endpoint.function_name] = endpoint
                response = {
                    "id": request["id"],
                    "result": {
                        "Endpoints": [e.to_dict() for e in endpoints],
                    },
                    "error": None
                }
                output = json.dumps(response)
                self.debug(f"output: {output}")
                return output
            else:
                self.debug(f"Calling {request['method']} - Params: {request['params']}")
                key = request["method"].replace(f"{self.plugin.name()}.", "")
                if key in self.route_map:
                    f = self.route_map[key].function
                    input_args = InputArgs.from_dict(request["params"][0])
                    response = f(input_args)
                    output = json.dumps(response.to_dict())
                    self.debug(f"output: {output}")
                    return output
                else:
                    response = {
                        "id": request["id"],
                        "result": None,
                        "error": f"Unknown method: {request['method']}"
                    }
                    output = json.dumps(response)
                    self.debug(f"output: {output}")
                    return output
        except Exception as ex:
            self.debug(f"Exception: {ex}")
            self.debug(traceback.format_exc())

    def serve(self):

        env_cookie = os.getenv("DTAC_PLUGINS")
        if env_cookie is None:
            print('============================ WARNING ============================')
            print('This is a DTAC plugin and is not designed to be executed directly')
            print('Please use the DTAC agent to load this plugin')
            print('==================================================================')
            sys.exit(-1)

        self.port = get_unused_tcp_port()

        options = [
            f'enc={quote(self.encryptor.key_string())}',
        ]

        print(f"CONNECT{{{{{self.plugin.name()}:{self.plugin.route_root()}:{self.rpc_proto}:{self.proto}:{self.ip}:{self.port}:{self.interface_version}:[{','.join(options)}]}}}}")
        sys.stdout.flush()


        # Check for certificate and key files passed via ENV variables
        cert = os.getenv("DTAC_TLS_CERT")
        key = os.getenv("DTAC_TLS_KEY")


        # Create a gRPC server
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

        # Check if both certificate and key are provided
        if cert and key:
            # Convert certificate and key strings into bytes
            certificate_chain = cert.encode('utf-8')
            private_key = key.encode('utf-8')

            # Create server SSL credentials
            server_credentials = ssl_server_credentials(
                [(private_key, certificate_chain)]
            )

            # Add secure port using credentials
            server.add_secure_port(f"[::]:{self.port}", server_credentials)
        else:
            # Add insecure port if no TLS credentials are provided
            server.add_insecure_port(f"[::]:{self.port}")

        dtac_tools.host.plugin_pb2_grpc.add_PluginServiceServicer_to_server(self, server)
        server.add_insecure_port(f"[::]:{self.port}")
        server.start()
        server.wait_for_termination()

    def get_port(self):
        return self.port