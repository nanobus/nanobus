using System.Threading.Tasks;
using NanoBus.Functions;

namespace Customers
{
    public class Inbound
    {
        private Outbound outbound;

        public Inbound(Outbound outbound)
        {
            this.outbound = outbound;
        }

        public async Task<Customer> CreateCustomer(Customer customer)
        {
            await outbound.SaveCustomer(customer);
            await outbound.CustomerCreated(customer);

            return customer;
        }

        public Task<Customer> GetCustomer(ulong id)
        {
            return outbound.FetchCustomer(id);
        }
    }

    public class Program
    {
        public static void Main(string[] args)
        {
            var (handlers, invoker, start) = HTTP.Initialize();

            var outbound = new OutboundImpl(invoker);
            var inbound = new Inbound(outbound);

            new InboundHandlers
            {
                CreateCustomer = inbound.CreateCustomer,
                GetCustomer = inbound.GetCustomer,
            }.Register(handlers);

            start();
        }
    }
}
