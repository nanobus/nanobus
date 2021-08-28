using System.Threading.Tasks;

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
}