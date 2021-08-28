using System.Threading.Tasks;
using MessagePack;
using NanoBus.Functions;

namespace Customers
{
#pragma warning disable CS8632 // The annotation for nullable reference types should only be used in code within a '#nullable' annotations context.
    [MessagePackObject]
    public class GetCustomerArgs
    {
        [Key("id")]
        public ulong Id { get; set; }
    }

    [MessagePackObject]
    public class Customer
    {
        [Key("id")]
        public ulong Id { get; set; }
        [Key("firstName")]
        public string FirstName { get; set; }
        [Key("middleName")]
        public string? MiddleName { get; set; }
        [Key("lastName")]
        public string LastName { get; set; }
        [Key("email")]
        public string Email { get; set; }
        [Key("address")]
        public Address? Address { get; set; }
    }

    [MessagePackObject]
    public class Address
    {
        [Key("line1")]
        public string Line1 { get; set; }
        [Key("line2")]
        public string? Line2 { get; set; }
        [Key("city")]
        public string City { get; set; }
        [Key("zip")]
        public string Zip { get; set; }
    }
#pragma warning restore CS8632 // The annotation for nullable reference types should only be used in code within a '#nullable' annotations context.

    public interface Outbound
    {
        public Task SaveCustomer(Customer customer);
        public Task CustomerCreated(Customer customer);
        public Task<Customer> FetchCustomer(ulong id);
    }

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

    public delegate Task<Customer> CreateCustomer(Customer customer);
    public delegate Task<Customer> GetCustomer(ulong id);

    public class InboundHandlers
    {
        public CreateCustomer CreateCustomer { get; set; }
        public GetCustomer GetCustomer { get; set; }

        public void Register(IHandlers handlers)
        {
            var codec = handlers.Codec();
            if (CreateCustomer != null)
            {
                handlers.RegisterHandler("/customers.v1.Inbound/createCustomer", async (input) =>
                {
                    var customer = codec.Decode<Customer>(input);
                    var result = await CreateCustomer(customer);
                    return codec.Encode(result);
                });
            }
            if (GetCustomer != null)
            {
                handlers.RegisterHandler("/customers.v1.Inbound/getCustomer", async (input) =>
                {
                    var args = codec.Decode<GetCustomerArgs>(input);
                    var result = await GetCustomer(args.Id);
                    return codec.Encode(result);
                });
            }
        }
    }
}
