using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Threading.Tasks;
using MessagePack;

namespace NanoBus.Functions
{
    public interface ICodec
    {
        byte[] Encode(object value);
        T Decode<T>(byte[] data);
    }

    public class MsgPackCodec : ICodec
    {
        public byte[] Encode(object value)
        {
            return MessagePackSerializer.Serialize(value);
        }

        public T Decode<T>(byte[] data)
        {
            return MessagePackSerializer.Deserialize<T>(data);
        }
    }

    public delegate Task<byte[]> Invoke(string operation, byte[] input);

    public class HTTPInvoker
    {
        private static HttpClient client = new HttpClient();

        private string baseURL;

        public HTTPInvoker(string baseURL)
        {
            this.baseURL = baseURL;
        }

        public async Task<byte[]> Invoke(string operation, byte[] input)
        {
            HttpContent content = new ByteArrayContent(input);
            var response = await client.PostAsync(baseURL + operation, content);
            var stream = await response.Content.ReadAsStreamAsync();
            var ms = new MemoryStream();
            stream.CopyTo(ms);
            return ms.ToArray();
        }
    }

    public class Invoker
    {
        private Invoke invoke;
        private ICodec codec;

        public Invoker(Invoke invoke, ICodec codec)
        {
            this.invoke = invoke;
            this.codec = codec;
        }

        public async Task Invoke(string operation, object value)
        {
            var data = codec.Encode(value);
            await invoke(operation, data);
        }

        public async Task<T> InvokeWithReturn<T>(string operation, object value)
        {
            var data = codec.Encode(value);
            var result = await invoke(operation, data);
            return codec.Decode<T>(result);
        }
    }

    public delegate Task<byte[]> Handler(byte[] input);

    public interface IHandlers
    {
        ICodec Codec();
        void RegisterHandler(string operation, Handler handler);
    }

    public class HTTPHandlers : IHandlers
    {
        private ICodec codec;
        private HttpListener listener;
        private Dictionary<string, Handler> handlers = new Dictionary<string, Handler>();

        public HTTPHandlers(ICodec codec)
        {
            this.codec = codec;
        }

        public ICodec Codec() { return codec; }

        public void RegisterHandler(string operation, Handler handler)
        {
            if (handler != null)
            {
                this.handlers[operation] = handler;
            }
        }

        private async Task HandleIncomingConnections()
        {
            bool runServer = true;

            // While a user hasn't visited the `shutdown` url, keep on handling requests
            while (runServer)
            {
                // Will wait here until we hear from a connection
                HttpListenerContext ctx = await listener.GetContextAsync();

                // Peel out the requests and response objects
                HttpListenerRequest req = ctx.Request;
                HttpListenerResponse resp = ctx.Response;

                if (req.HttpMethod != "POST")
                {
                    resp.Close();
                    continue;
                }

                Handler handler = handlers[req.Url.AbsolutePath];
                if (handler == null)
                {
                    resp.Close();
                    continue;
                }

                byte[] input;
                if (req.HasEntityBody)
                {
                    using (var memoryStream = new MemoryStream())
                    {
                        req.InputStream.CopyTo(memoryStream);
                        input = memoryStream.ToArray();
                    }
                }
                else
                {
                    input = new byte[] { };
                }

                byte[] output = await handler(input);

                // Write out to the response stream (asynchronously), then close it
                await resp.OutputStream.WriteAsync(output, 0, output.Length);
                resp.Close();
            }
        }

        public void Listen(int port, string hostname)
        {
            try
            {
                // Create a Http server and start listening for incoming connections
                listener = new HttpListener();
                string url = string.Format("http://{0}:{1}/", hostname, port);
                listener.Prefixes.Add(url);
                listener.Start();
                Console.WriteLine("Listening for connections on {0}", url);

                // Handle requests
                Task listenTask = HandleIncomingConnections();
                listenTask.GetAwaiter().GetResult();

                // Close the listener
                listener.Close();
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
            }
        }
    }
}
