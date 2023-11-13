import json

class InputArgs:
    def __init__(self, headers=None, params=None, body=None):
        if headers is None:
            headers = {}
        if params is None:
            params = {}
        self.headers = headers
        self.params = params
        self.body = body

    def to_dict(self):
        return {
            "headers": self.headers,
            "params": self.params,
            "body": self.body.decode('utf-8') if self.body else None
        }

    def to_json(self):
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str):
        data = json.loads(json_str)
        body = data.get("body")
        if body is not None:
            body = body.encode('utf-8')
        return cls(
            headers=data.get("headers"),
            params=data.get("params"),
            body=body
        )

    @classmethod
    def from_dict(cls, data):
        body = data.get("body")
        if body is not None:
            body = body.encode('utf-8')
        return cls(
            headers=data.get("headers"),
            params=data.get("params"),
            body=body
        )

class ReturnVal:
    def __init__(self, headers=None, params=None, value=None):
        if headers is None:
            headers = {}
        if params is None:
            params = {}
        self.headers = headers
        self.params = params
        self.value = value

    def to_dict(self):
        return {
            "headers": self.headers,
            "params": self.params,
            "value": self.value
        }

    def to_json(self):
        return json.dumps(self.to_dict())


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