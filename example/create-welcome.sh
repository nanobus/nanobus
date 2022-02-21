nanogen new go welcome \
    description="Dapr Welcome App" \
    version=1.0.0 \
    module="github.com/nanobus/nanobus/example/welcome" \
    package="welcome"

cd welcome

cat > schema.apex <<EOF
namespace "welcome.v1"

role Inbound {
  greetCustomer{customer: Customer}: void
}

role Outbound {
  sendEmail(email: string, message: string): void
}

type Customer {
  firstName: string
  lastName: string
  email: string @email
}
EOF

cat > bus.yaml <<EOF
compute:
  type: mux
  with: {}

subscriptions:

  #- pubsub: pubsub
  #  topic: test_topic
  #  function: receiveMessage

resiliency:
  retries:
    pubsub:
      policy: constant
      duration: 3s

  circuitBreakers:
    pubsub:
      maxRequests: 2
      timeout: 30s

inbound:

  receiveMessage:
    summary: Receives the customer message
    actions:
      - summary: Routing messages
        name: route
        with:
          routes:
            - when: input.type == 'customer.created'
              then:
                - summary: Send to route
                  name: invoke
                  with:
                    function: welcome.v1.Inbound/greetCustomer
                    input: input.data
                  retry: pubsub
                  circuitBreaker: pubsub

outbound:

  'welcome.v1.Outbound':
    sendEmail:
      summary: Sends an email to the customer
      actions:
        - summary: Pretend to send an email by logging the message
          name: log
          with:
            format: "Sending email to %s with message %q"
            args:
              - input.email
              - input.message
EOF

cat > dapr.sh <<EOF
# !/bin/sh
INVOKE_BASE_URL=http://localhost:8001 dapr run -d ../components --app-id welcome --app-port 32322 --dapr-http-port 3501 -- nanobus --http-listen-addr :8082 --rest-listen-addr :8092 --bus-listen-addr :32322 bus.yaml
EOF

chmod +x dapr.sh

cat > start.sh <<EOF
# !/bin/sh
go run cmd/main.go --http-listen-addr :8001 --outbound-base-uri http://localhost:32322/outbound
EOF

chmod +x start.sh