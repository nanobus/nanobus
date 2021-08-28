using System.Threading.Tasks;
using NanoBus.Functions;

namespace Customers
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var (handlers, invoker, start) = Server.Initialize();

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
