# !/bin/sh
nanobus --components-path ../components --metrics-port 19091 --app-id customers --placement-host-address localhost:50005 --http-listen-addr :18081 --rest-listen-addr :8091 --bus-listen-addr localhost:32321 bus.yaml
