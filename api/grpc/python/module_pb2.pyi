import plugin_pb2 as _plugin_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class ModuleRegisterRequest(_message.Message):
    __slots__ = ["config", "default_secure"]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_SECURE_FIELD_NUMBER: _ClassVar[int]
    config: str
    default_secure: bool
    def __init__(self, config: _Optional[str] = ..., default_secure: bool = ...) -> None: ...

class ModuleRegisterResponse(_message.Message):
    __slots__ = ["module_type", "capabilities"]
    MODULE_TYPE_FIELD_NUMBER: _ClassVar[int]
    CAPABILITIES_FIELD_NUMBER: _ClassVar[int]
    module_type: str
    capabilities: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, module_type: _Optional[str] = ..., capabilities: _Optional[_Iterable[str]] = ...) -> None: ...

class TokenRequest(_message.Message):
    __slots__ = ["scopes", "expires_in"]
    SCOPES_FIELD_NUMBER: _ClassVar[int]
    EXPIRES_IN_FIELD_NUMBER: _ClassVar[int]
    scopes: _containers.RepeatedScalarFieldContainer[str]
    expires_in: int
    def __init__(self, scopes: _Optional[_Iterable[str]] = ..., expires_in: _Optional[int] = ...) -> None: ...

class TokenRefreshRequest(_message.Message):
    __slots__ = ["refresh_token"]
    REFRESH_TOKEN_FIELD_NUMBER: _ClassVar[int]
    refresh_token: str
    def __init__(self, refresh_token: _Optional[str] = ...) -> None: ...

class TokenResponse(_message.Message):
    __slots__ = ["access_token", "refresh_token", "expires_in", "token_type"]
    ACCESS_TOKEN_FIELD_NUMBER: _ClassVar[int]
    REFRESH_TOKEN_FIELD_NUMBER: _ClassVar[int]
    EXPIRES_IN_FIELD_NUMBER: _ClassVar[int]
    TOKEN_TYPE_FIELD_NUMBER: _ClassVar[int]
    access_token: str
    refresh_token: str
    expires_in: int
    token_type: str
    def __init__(self, access_token: _Optional[str] = ..., refresh_token: _Optional[str] = ..., expires_in: _Optional[int] = ..., token_type: _Optional[str] = ...) -> None: ...
