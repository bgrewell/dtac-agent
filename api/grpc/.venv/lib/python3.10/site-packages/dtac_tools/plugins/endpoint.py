import json

class PluginEndpoint:
    def __init__(self, function, path, action, uses_auth, expected_args=None, expected_body=None, expected_output=None):
        self.function = function
        self.function_name = function.__name__
        self.path = path
        self.action = action
        self.uses_auth = uses_auth
        self.expected_args = expected_args
        self.expected_body = expected_body
        self.expected_output = expected_output

    def to_dict(self):
        return {
            "function_name": self.function_name,
            "path": self.path,
            "action": self.action,
            "uses_auth": self.uses_auth,
            "expected_args": self.expected_args,
            "expected_body": self.expected_body,
            "expected_output": self.expected_output
        }

    def to_json(self):
        return json.dumps(self, default=lambda o: {k: v for k, v in o.__dict__.items() if v is not None and not callable(v)}, indent=4)