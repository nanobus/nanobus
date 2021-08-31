using System.Threading.Tasks;
using MessagePack;

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
        [Key("state")]
        public string State { get; set; }
        [Key("zip")]
        public string Zip { get; set; }
    }
#pragma warning restore CS8632 // The annotation for nullable reference types should only be used in code within a '#nullable' annotations context.

    public delegate Task<Customer> CreateCustomer(Customer customer);
    public delegate Task<Customer> GetCustomer(ulong id);

    public class Inbound
    {
        public CreateCustomer CreateCustomer { get; set; }
        public GetCustomer GetCustomer { get; set; }
    }

    public interface Outbound
    {
        Task SaveCustomer(Customer customer);
        Task CustomerCreated(Customer customer);
        Task<Customer> FetchCustomer(ulong id);
    }
}
