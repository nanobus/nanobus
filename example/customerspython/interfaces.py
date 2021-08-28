# Python 3 server example
from http.server import BaseHTTPRequestHandler, HTTPServer
from requests.api import request
import msgpack
from typing import Callable, Dict, Tuple, Type, TypeVar, Generic
from dataclasses import asdict
from dacite import from_dict
import requests

hostName = "localhost"
serverPort = 9000

T = TypeVar('T')


class Codec(Generic[T]):
    def encode(self, value: T) -> bytes:
        pass

    def decode(self, data: bytes, data_class: Type[T]) -> T:
        pass


class MsgPackCodec(Codec):
    def encode(self, value: T) -> bytes:
        return msgpack.packb(asdict(value))

    def decode(self, data: bytes, data_class: Type[T]) -> T:
        return from_dict(data_class, msgpack.unpackb(data))


class HTTPInvoker:
    def __init__(self, baseURL: str):
        self.baseURL = baseURL

    def invoke(self, operation: str, input: bytes) -> bytes:
        resp = requests.post(self.baseURL + operation, data=input)
        return resp.content


class Invoker(Generic[T]):
    def __init__(self, call: Callable[[str, bytes], bytes], codec: Codec):
        self.call = call
        self.codec = codec

    def invoke(self, operation: str, value: any):
        data = self.codec.encode(value)
        self.call(operation, data)

    def invokeWithReturn(self, operation: str, value: any, data_class: Type[T]):
        data = self.codec.encode(value)
        result = self.call(operation, data)
        return self.codec.decode(result, data_class)


class Handlers:
    def __init__(self, codec: Codec):
        self.handlers = {}  # Dict[str, Callable[[bytes], bytes]] = {}
        self.codec = codec

    def register_handler(self, operation: str, handler: Callable[[bytes], bytes]):
        if handler != None:
            self.handlers[operation] = handler

    def register_handlers(self, handlers: Dict[str, Callable[[bytes], bytes]]):
        for operation in handlers:
            self.register_handler(operation, handlers[operation])

    def get_handler(self, operation: str):
        return self.handlers[operation]


class HTTPRequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        handler = self.server.handlers.get_handler(self.path)
        if handler == None:
            self.send_response(404)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(bytes("Not Found"))
            return

        # <--- Gets the size of data
        content_length = int(self.headers['Content-Length'])
        # <--- Gets the data itself
        post_data = self.rfile.read(content_length)

        result = handler(post_data)

        self.send_response(200)
        self.send_header("Content-type", "application/msgpack")
        self.end_headers()
        self.wfile.write(result)


def http_initialize() -> Tuple[Handlers, Invoker, Callable[[], None]]:
    codec = MsgPackCodec()
    handlers = Handlers(codec)
    webServer = HTTPServer((hostName, serverPort), HTTPRequestHandler)
    webServer.handlers = handlers

    invoke = HTTPInvoker("http://localhost:32321/outbound").invoke
    invoker = Invoker(invoke, codec)

    def start():
        print("Server started http://%s:%s" % (hostName, serverPort))
        try:
            webServer.serve_forever()
        except KeyboardInterrupt:
            pass

        webServer.server_close()
        print("Server stopped.")

    return (handlers, invoker, start)
