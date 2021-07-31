# !/bin/sh
#INVOKE_BASE_URL=http://localhost:8000 dapr run -d ../components --app-id customers --app-port 32321 --dapr-http-port 3500 -- nanobus --http-listen-addr :8081 --rest-listen-addr :8091 --bus-listen-addr localhost:32321 bus.yaml
INVOKE_BASE_URL=http://localhost:8000 nanobus --components-path ../components --metrics-port 9091 --app-id customers --placement-host-address localhost:50005 --http-listen-addr :8081 --rest-listen-addr :8091 --bus-listen-addr localhost:32321 bus.yaml
