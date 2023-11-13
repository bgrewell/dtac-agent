from dtac_tools.host.plugin_pb2_grpc import PluginServiceServicer

class PluginHost(PluginServiceServicer):
    def serve(self) -> None:
        raise NotImplementedError

    def get_port(self) -> int:
        raise NotImplementedError