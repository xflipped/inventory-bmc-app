# Redfish  Tools


## Build
```
docker build -t ghcr.io/foliagecp/discovery-bmc-app -f Dockerfile.discovery .

docker build -t ghcr.io/foliagecp/inventory-bmc-app -f Dockerfile.inventory .

docker build -t ghcr.io/foliagecp/led-bmc-app -f Dockerfile.led .

docker build -t ghcr.io/foliagecp/reset-bmc-app -f Dockerfile.reset .

docker build -t ghcr.io/foliagecp/subscribe-bmc-app -f Dockerfile.subscribe .
```

## Test
```
# add to foliage/docker-compose.yaml

  inventory-bmc:
    image: ghcr.io/foliagecp/inventory-bmc-app:${INVENTORY_BMC_VERSION:-latest}
    hostname: inventory-bmc
    profiles: ["app"]
    ports:
      - "31001:31001"
    depends_on:
      proxy:
        condition: service_healthy
      sfmanager:
        condition: service_healthy
      sfworker:
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
      test: curl --fail http://0.0.0.0:31001/readyz || exit 1
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s

  discovery-bmc:
    image: ghcr.io/foliagecp/discovery-bmc-app:${DISCOVERY_BMC_VERSION:-latest}
    hostname: discovery-bmc
    profiles: ["app"]
    ports:
      - "1900:1900/udp"
      - "31002:31002"
    depends_on:
      proxy:
        condition: service_healthy
      sfmanager:
        condition: service_healthy
      sfworker:
        condition: service_healthy
    environment:
      KAFKA_ADDR: ${KAFKA_ADDR}
      CMDB_ADDR: ${CMDB_ADDR}
      CMDB_PORT: ${CMDB_PORT}
      SSDP_MONITOR: ${SSDP_MONITOR}
    healthcheck:
      test: curl --fail http://0.0.0.0:31002/readyz || exit 1
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s

  led-bmc:
    image: ghcr.io/foliagecp/led-bmc-app:${LED_BMC_VERSION:-latest}
    hostname: led-bmc
    profiles: ["app"]
    ports:
      - "31003:31003"
    depends_on:
      inventory-bmc:
        condition: service_healthy
    networks:
      default:
        aliases:
          - led-bmc
    environment:
      KAFKA_ADDR: ${KAFKA_ADDR}
      CMDB_ADDR: ${CMDB_ADDR}
      CMDB_PORT: ${CMDB_PORT}
    healthcheck:
      test: curl --fail http://0.0.0.0:31003/readyz || exit 1
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s

  reset-bmc:
    image: ghcr.io/foliagecp/reset-bmc-app:${RESET_BMC_VERSION:-latest}
    hostname: reset-bmc
    profiles: ["app"]
    ports:
      - "31004:31004"
    depends_on:
      inventory-bmc:
        condition: service_healthy
    networks:
      default:
        aliases:
          - reset-bmc
    environment:
      KAFKA_ADDR: ${KAFKA_ADDR}
      CMDB_ADDR: ${CMDB_ADDR}
      CMDB_PORT: ${CMDB_PORT}
    healthcheck:
      test: curl --fail http://0.0.0.0:31004/readyz || exit 1
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s

  subscribe-bmc:
    image: ghcr.io/foliagecp/subscribe-bmc-app:${SUBSCRIBE_BMC_VERSION:-latest}
    hostname: subscribe-bmc
    profiles: ["app"]
    ports:
      - "31005:31005"
    depends_on:
      inventory-bmc:
        condition: service_healthy
    networks:
      default:
        aliases:
          - subscribe-bmc
    environment:
      KAFKA_ADDR: ${KAFKA_ADDR}
      CMDB_ADDR: ${CMDB_ADDR}
      CMDB_PORT: ${CMDB_PORT}
    healthcheck:
      test: curl --fail http://0.0.0.0:31005/readyz || exit 1
      interval: 10s
      timeout: 5s
      retries: 8
      start_period: 10s
```
