namespace "welcome.v1"
  @version("0.1.0")

role Inbound @events @uses([Outbound]) {
  greetCustomer{customer: Customer}: void
}

role Outbound @provider {
  sendEmail(email: string, message: string): void
}

type Customer {
  firstName: string
  lastName: string
  email: string @email
}