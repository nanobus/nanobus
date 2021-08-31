using System;
using System.Threading.Tasks;
using MessagePack;
using NanoBus.Functions;

namespace Customers
{
    public class OutboundImpl : Outbound
    {
        private Invoker invoker;

        public OutboundImpl(Invoker invoker)
        {
            this.invoker = invoker;
        }

        public async Task SaveCustomer(Customer customer)
        {
            await invoker.Invoke("/customers.v1.Outbound/saveCustomer", customer);
        }

        public async Task CustomerCreated(Customer customer)
        {
            await invoker.Invoke("/customers.v1.Outbound/customerCreated", customer);
        }

        public async Task<Customer> FetchCustomer(ulong id)
        {
            var args = new GetCustomerArgs
            {
                Id = id,
            };
            return await invoker.InvokeWithReturn<Customer>("/customers.v1.Outbound/fetchCustomer", args);
        }
    }


    public delegate void Starter();

    public class Adapter
    {
        private HTTPHandlers handlers;
        private ICodec codec;
        private Invoker invoker;

        private static string GetEnvironmentVariable(string name, string defaultValue)
            => Environment.GetEnvironmentVariable(name) is string v && v.Length > 0 ? v : defaultValue;

        public Adapter()
        {
            var outboundBaseURL = GetEnvironmentVariable("OUTBOUND_BASE_URL", "http://localhost:32321/outbound");
            this.codec = new MsgPackCodec();
            var invoke = new HTTPInvoker(outboundBaseURL);
            this.invoker = new Invoker(invoke.Invoke, this.codec);
            this.handlers = new HTTPHandlers(this.codec);
        }

        public void Start()
        {
            var host = GetEnvironmentVariable("HOST", "localhost");
            var port = int.Parse(GetEnvironmentVariable("PORT", "9000"));
            this.handlers.Listen(port, host);
        }

        public void RegisterInboundHanders(Inbound handlers)
        {
            if (handlers.CreateCustomer != null)
            {
                this.handlers.RegisterHandler("/customers.v1.Inbound/createCustomer", async (input) =>
                {
                    var customer = this.codec.Decode<Customer>(input);
                    var result = await handlers.CreateCustomer(customer);
                    return codec.Encode(result);
                });
            }
            if (handlers.GetCustomer != null)
            {
                this.handlers.RegisterHandler("/customers.v1.Inbound/getCustomer", async (input) =>
                {
                    var args = this.codec.Decode<GetCustomerArgs>(input);
                    var result = await handlers.GetCustomer(args.Id);
                    return codec.Encode(result);
                });
            }
        }

        public Outbound NewOutbound()
        {
            return new OutboundImpl(this.invoker);
        }
    }
}
