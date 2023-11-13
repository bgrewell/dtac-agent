import json

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
    def register(self, args):
        raise NotImplementedError("this method must be implemented")
