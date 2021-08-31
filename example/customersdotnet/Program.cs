namespace Customers
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var adapter = new Adapter();
            var outbound = adapter.NewOutbound();
            var service = new Service(outbound);

            adapter.RegisterInboundHanders(new Inbound
            {
                CreateCustomer = service.CreateCustomer,
                GetCustomer = service.GetCustomer,
            });

            adapter.Start();
        }
    }
}
