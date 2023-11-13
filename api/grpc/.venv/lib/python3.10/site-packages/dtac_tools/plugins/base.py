import json
from pydantic import BaseModel
from typing import TypeVar, Type
from ..host.helpers.debug_sender import DebugSender

# Base Config class, this is used to help deserialize configuration options for plugins
class BaseConfig(BaseModel):
    pass


class PluginShared:
    @staticmethod
    def name() -> str:
        return "UnnamedPlugin"
    
    @staticmethod
    def route_root() -> str:
        return ""
    
    @staticmethod
    def serialize(v: object) -> str:
        return json.dumps(v)

class PluginBase(PluginShared):

    T = TypeVar("T", bound=BaseConfig)

    def enable_debugging(self):
        self.debug_sender = DebugSender(5678)
        self.debug(f"Plugin debugging enabled on port {5678}")
        self.debug(f"Plugin Name: {self.plugin.name()}")

    def debug(self, message):
        if self.debug_sender is not None:
            self.debug_sender.write(message + "\n")

    def load_config(self, config_json: str, config_class: Type[T]) -> T:
        return config_class.parse_raw(config_json)

    def register(self, args):
        raise NotImplementedError("this method must be implemented")
