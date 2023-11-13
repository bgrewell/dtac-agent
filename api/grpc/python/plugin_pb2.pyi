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

class PluginRequest(_message.Message):
    __slots__ = ["method", "input_args"]
    METHOD_FIELD_NUMBER: _ClassVar[int]
    INPUT_ARGS_FIELD_NUMBER: _ClassVar[int]
    method: str
    input_args: InputArgs
    def __init__(self, method: _Optional[str] = ..., input_args: _Optional[_Union[InputArgs, _Mapping]] = ...) -> None: ...

class PluginResponse(_message.Message):
    __slots__ = ["id", "return_val", "error"]
    ID_FIELD_NUMBER: _ClassVar[int]
    RETURN_VAL_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    id: int
    return_val: ReturnVal
    error: str
    def __init__(self, id: _Optional[int] = ..., return_val: _Optional[_Union[ReturnVal, _Mapping]] = ..., error: _Optional[str] = ...) -> None: ...

class StringList(_message.Message):
    __slots__ = ["values"]
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, values: _Optional[_Iterable[str]] = ...) -> None: ...

class InputArgs(_message.Message):
    __slots__ = ["headers", "params", "body"]
    class HeadersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    class ParamsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    headers: _containers.MessageMap[str, StringList]
    params: _containers.MessageMap[str, StringList]
    body: bytes
    def __init__(self, headers: _Optional[_Mapping[str, StringList]] = ..., params: _Optional[_Mapping[str, StringList]] = ..., body: _Optional[bytes] = ...) -> None: ...

class ReturnVal(_message.Message):
    __slots__ = ["headers", "params", "value"]
    class HeadersEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    class ParamsEntry(_message.Message):
        __slots__ = ["key", "value"]
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: StringList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...
    HEADERS_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    headers: _containers.MessageMap[str, StringList]
    params: _containers.MessageMap[str, StringList]
    value: str
    def __init__(self, headers: _Optional[_Mapping[str, StringList]] = ..., params: _Optional[_Mapping[str, StringList]] = ..., value: _Optional[str] = ...) -> None: ...

class RegisterArgs(_message.Message):
    __slots__ = ["config", "default_secure"]
    CONFIG_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_SECURE_FIELD_NUMBER: _ClassVar[int]
    config: str
    default_secure: bool
    def __init__(self, config: _Optional[str] = ..., default_secure: bool = ...) -> None: ...

class RegisterReply(_message.Message):
    __slots__ = ["endpoints"]
    ENDPOINTS_FIELD_NUMBER: _ClassVar[int]
    endpoints: _containers.RepeatedCompositeFieldContainer[PluginEndpoint]
    def __init__(self, endpoints: _Optional[_Iterable[_Union[PluginEndpoint, _Mapping]]] = ...) -> None: ...

class PluginEndpoint(_message.Message):
    __slots__ = ["path", "action", "uses_auth", "expected_args", "expected_body", "expected_output"]
    PATH_FIELD_NUMBER: _ClassVar[int]
    ACTION_FIELD_NUMBER: _ClassVar[int]
    USES_AUTH_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_ARGS_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_BODY_FIELD_NUMBER: _ClassVar[int]
    EXPECTED_OUTPUT_FIELD_NUMBER: _ClassVar[int]
    path: str
    action: str
    uses_auth: bool
    expected_args: str
    expected_body: str
    expected_output: str
    def __init__(self, path: _Optional[str] = ..., action: _Optional[str] = ..., uses_auth: bool = ..., expected_args: _Optional[str] = ..., expected_body: _Optional[str] = ..., expected_output: _Optional[str] = ...) -> None: ...

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
