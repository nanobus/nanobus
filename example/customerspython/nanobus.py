import asyncio
from logging import WARNING
from typing import Awaitable, Callable, Dict, NoReturn, Type, TypeVar, Generic
import msgpack
from dataclasses import asdict
import dacite
from aiohttp import web, ClientSession
import uvicorn


T = TypeVar('T')

session: ClientSession = None


class Codec(Generic[T]):
    def encode(self, value: T) -> bytes:
        pass

    def decode(self, data: bytes, data_class: Type[T]) -> T:
        pass


class MsgPackCodec(Codec):
    def encode(self, value: T) -> bytes:
        return msgpack.packb(asdict(value))

    def decode(self, data: bytes, data_class: Type[T]) -> T:
        return dacite.from_dict(data_class, msgpack.unpackb(data))


class HTTPInvoker:
    def __init__(self, baseURL: str):
        self.baseURL = baseURL

    async def invoke(self, namespace: str, operation: str, input: bytes) -> bytes:
        async with session.post(self.baseURL + '/' + namespace + '/' + operation, data=input) as resp:
            return await resp.read()


class Invoker(Generic[T]):
    def __init__(self, call: Callable[[str, str, bytes], bytes], codec: Codec):
        self.call = call
        self.codec = codec

    async def invoke(self, namespace: str, operation: str, value: any):
        data = self.codec.encode(value)
        await self.call(namespace, operation, data)

    async def invoke_with_return(self, namespace: str, operation: str, value: any, data_class: Type[T]):
        data = self.codec.encode(value)
        result = await self.call(namespace, operation, data)
        return self.codec.decode(result, data_class)


class Handlers:
    def __init__(self, codec: Codec):
        self.handlers: Dict[str, Callable[[bytes],
                                          Awaitable[bytes]]] = {}
        self.codec = codec

    def register_handler(self, namespace: str, operation: str, handler: Callable[[bytes], Awaitable[bytes]]):
        if handler != None:
            self.handlers[namespace + '/' + operation] = handler

    def get_handler(self, namespace: str, operation: str):
        return self.handlers.get(namespace + '/' + operation, operation)


class AIOHTTPServer:
    def __init__(self, handlers: Handlers):
        global session
        session = ClientSession()
        self.handlers = handlers
        app = web.Application()
        app.add_routes([web.post('/{namespace}/{operation}', self.handle)])
        app.cleanup_ctx.append(self.client_session_ctx)
        self.app = app

    async def handle(self, request):
        namespace = request.match_info['namespace']
        operation = request.match_info['operation']
        handler = self.handlers.get_handler(namespace, operation)
        if handler == None:
            return web.HTTPNotFound()

        post_data = await request.read()
        result = await handler(post_data)

        return web.Response(body=result, content_type="application/msgpack")

    async def client_session_ctx(self, app: web.Application) -> NoReturn:
        app['client_session'] = session
        yield
        await session.close()

    def run(self, host: str, port: int):
        web.run_app(self.app, host=host, port=port)


class UvicornServer:
    def __init__(self, handlers: Handlers):
        global session
        self.handlers = handlers

    async def read_body(self, receive):
        """
        Read and return the entire body from an incoming ASGI message.
        """
        body = b''
        more_body = True

        while more_body:
            message = await receive()
            body += message.get('body', b'')
            more_body = message.get('more_body', False)

        return body

    async def __call__(self, scope, receive, send):
        parts = scope["path"].split('/')
        if (len(parts) < 3):
            await send({
                'type': 'http.response.start',
                'status': 400,
                'headers': [
                    [b'content-type', b'text/plain'],
                ]
            })
            await send({
                'type': 'http.response.body',
                'body': b'Bad request',
            })
            return

        namespace = parts[1]
        operation = parts[2]

        handler = self.handlers.get_handler(namespace, operation)
        if handler == None:
            await send({
                'type': 'http.response.start',
                'status': 404,
                'headers': [
                    [b'content-type', b'text/plain'],
                ]
            })
            await send({
                'type': 'http.response.body',
                'body': b'Not found',
            })
            return

        body = await self.read_body(receive)
        result = await handler(body)

        await send({
            'type': 'http.response.start',
            'status': 200,
            'headers': [
                [b'content-type', b'application/msgpack'],
            ]
        })
        await send({
            'type': 'http.response.body',
            'body': result,
        })

    async def _serve(self, config: uvicorn.Config):
        global session
        session = ClientSession()
        server = uvicorn.Server(config)
        await server.serve()
        await session.close()

    def run(self, host: str, port: int):
        loop = asyncio.new_event_loop()
        config = uvicorn.Config(app=self, host=host,
                                port=port, loop=loop, log_level=WARNING)
        loop.run_until_complete(self._serve(config))
