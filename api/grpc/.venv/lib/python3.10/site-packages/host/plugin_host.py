from dtac_tools.plugins import PluginBase

class PluginHost:
    def serve(self) -> None:
        raise NotImplementedError

    def get_port(self) -> int:
        raise NotImplementedError

def new_plugin_host(plugin: PluginBase, cookie: str) -> PluginHost:
    # Add logic to return a concrete instance of PluginHost, like DefaultPluginHost
    pass
