from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class RecommendationRequest(_message.Message):
    __slots__ = ("user_id", "topic", "limit")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    TOPIC_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    topic: str
    limit: int
    def __init__(self, user_id: _Optional[str] = ..., topic: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class PostItem(_message.Message):
    __slots__ = ("post_id", "content", "author_id")
    POST_ID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_ID_FIELD_NUMBER: _ClassVar[int]
    post_id: str
    content: str
    author_id: str
    def __init__(self, post_id: _Optional[str] = ..., content: _Optional[str] = ..., author_id: _Optional[str] = ...) -> None: ...

class RecommendationResponse(_message.Message):
    __slots__ = ("posts",)
    POSTS_FIELD_NUMBER: _ClassVar[int]
    posts: _containers.RepeatedCompositeFieldContainer[PostItem]
    def __init__(self, posts: _Optional[_Iterable[_Union[PostItem, _Mapping]]] = ...) -> None: ...
