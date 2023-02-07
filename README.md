# Redfish  Tools


## Build
```
docker build -t ghcr.io/foliagecp/discovery-bmc-app -f Dockerfile.discovery .

docker build -t ghcr.io/foliagecp/inventory-bmc-app -f Dockerfile.inventory .
```

## Test
```
# add to foliage/docker-compose.yaml

  inventory-bmc:
    image: ghcr.io/xflipped/inventory-bmc-app:${INVENTORY_BMC_VERSION:-latest}
    hostname: inventory-bmc
    ports:
      - "31001:31001"
    depends_on:
      proxy:
        condition: service_healthy
    networks:
      default:
        aliases:
          - inventory-bmc
    environment:
      KAFKA_ADDR: ${KAFKA_ADDR}
      CMDB_ADDR: ${CMDB_ADDR}
      CMDB_PORT: ${CMDB_PORT}
    healthcheck:
      test: "nc -z localhost 31001"
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s

  discovery-bmc:
    image: ghcr.io/foliagecp/discovery-bmc-app:${DISCOVERY_BMC_VERSION:-latest}
    hostname: discovery-bmc
    network_mode: "host"
    ports:
      - "1900:1900/udp"
    depends_on:
      inventory-bmc:
        condition: service_healthy
    environment:
      KAFKA_ADDR: 127.0.0.1:9094
      CMDB_ADDR: 127.0.0.1
      CMDB_PORT: 31415
```
