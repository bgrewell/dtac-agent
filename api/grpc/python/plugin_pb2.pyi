from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class LogLevel(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    DEBUG: _ClassVar[LogLevel]
    INFO: _ClassVar[LogLevel]
    WARNING: _ClassVar[LogLevel]
    ERROR: _ClassVar[LogLevel]
    FATAL: _ClassVar[LogLevel]
DEBUG: LogLevel
INFO: LogLevel
WARNING: LogLevel
ERROR: LogLevel
FATAL: LogLevel

class EndpointRequestMessage(_message.Message):
    __slots__ = ["method", "request"]
    METHOD_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    method: str
    request: EndpointRequest
    def __init__(self, method: _Optional[str] = ..., request: _Optional[_Union[EndpointRequest, _Mapping]] = ...) -> None: ...

class EndpointResponseMessage(_message.Message):
    __slots__ = ["id", "response", "error"]
    ID_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    id: int
    response: EndpointResponse
    error: str
    def __init__(self, id: _Optional[int] = ..., response: _Optional[_Union[EndpointResponse, _Mapping]] = ..., error: _Optional[str] = ...) -> None: ...

class EndpointRequest(_message.Message):
    __slots__ = ["metadata", "headers", "parameters", "body"]
    class MetadataEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class HeadersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    class ParametersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    METADATA_FIELD_NUMBER: _ClassVar[int]
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    PARAMETERS_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    metadata: _containers.ScalarMap[str, str]
    headers: _containers.MessageMap[str, StringList]
    parameters: _containers.MessageMap[str, StringList]
    body: bytes
    def __init__(self, metadata: _Optional[_Mapping[str, str]] = ..., headers: _Optional[_Mapping[str, StringList]] = ..., parameters: _Optional[_Mapping[str, StringList]] = ..., body: _Optional[bytes] = ...) -> None: ...

class EndpointResponse(_message.Message):
    __slots__ = ["metadata", "headers", "parameters", "value"]
    class MetadataEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    class HeadersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    class ParametersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    METADATA_FIELD_NUMBER: _ClassVar[int]
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    PARAMETERS_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    metadata: _containers.ScalarMap[str, str]
    headers: _containers.MessageMap[str, StringList]
    parameters: _containers.MessageMap[str, StringList]
    value: bytes
    def __init__(self, metadata: _Optional[_Mapping[str, str]] = ..., headers: _Optional[_Mapping[str, StringList]] = ..., parameters: _Optional[_Mapping[str, StringList]] = ..., value: _Optional[bytes] = ...) -> None: ...

class StringList(_message.Message):
    __slots__ = ["values"]
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, values: _Optional[_Iterable[str]] = ...) -> None: ...

class RegisterRequest(_message.Message):
    __slots__ = ["config", "default_secure"]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_SECURE_FIELD_NUMBER: _ClassVar[int]
    config: str
    default_secure: bool
    def __init__(self, config: _Optional[str] = ..., default_secure: bool = ...) -> None: ...

class RegisterResponse(_message.Message):
    __slots__ = ["endpoints"]
    ENDPOINTS_FIELD_NUMBER: _ClassVar[int]
    endpoints: _containers.RepeatedCompositeFieldContainer[PluginEndpoint]
    def __init__(self, endpoints: _Optional[_Iterable[_Union[PluginEndpoint, _Mapping]]] = ...) -> None: ...

class PluginEndpoint(_message.Message):
    __slots__ = ["path", "action", "secure", "auth_group", "expected_metadata_schema", "expected_headers_schema", "expected_parameters_schema", "expected_body_schema", "expected_output_schema"]
    PATH_FIELD_NUMBER: _ClassVar[int]
    ACTION_FIELD_NUMBER: _ClassVar[int]
    SECURE_FIELD_NUMBER: _ClassVar[int]
    AUTH_GROUP_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_METADATA_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_HEADERS_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_PARAMETERS_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_BODY_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_OUTPUT_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    path: str
    action: str
    secure: bool
    auth_group: str
    expected_metadata_schema: str
    expected_headers_schema: str
    expected_parameters_schema: str
    expected_body_schema: str
    expected_output_schema: str
    def __init__(self, path: _Optional[str] = ..., action: _Optional[str] = ..., secure: bool = ..., auth_group: _Optional[str] = ..., expected_metadata_schema: _Optional[str] = ..., expected_headers_schema: _Optional[str] = ..., expected_parameters_schema: _Optional[str] = ..., expected_body_schema: _Optional[str] = ..., expected_output_schema: _Optional[str] = ...) -> None: ...

class LoggingArgs(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class LogField(_message.Message):
    __slots__ = ["key", "value"]
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: str
    value: str
    def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...

class LogMessage(_message.Message):
    __slots__ = ["level", "message", "fields"]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    level: LogLevel
    message: str
    fields: _containers.RepeatedCompositeFieldContainer[LogField]
    def __init__(self, level: _Optional[_Union[LogLevel, str]] = ..., message: _Optional[str] = ..., fields: _Optional[_Iterable[_Union[LogField, _Mapping]]] = ...) -> None: ...
