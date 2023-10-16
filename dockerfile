FROM golang:1.21.0 as build

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

FROM quay.io/keycloak/keycloak:22.0.4 as builder

ENV KC_DB=postgres

WORKDIR /opt/keycloak

RUN keytool -genkeypair -storepass password -storetype PKCS12 -keyalg RSA -keysize 2048 -dname "CN=server" -alias server -ext "SAN:c=DNS:localhost,IP:192.168.3.6" -keystore conf/server.keystore
RUN /opt/keycloak/bin/kc.sh build

## Deploy
FROM busybox

## Copiare i certificati per connettersi a mongodb atlas
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /app/main /opt/main

COPY --from=builder /opt/keycloak/ /opt/keycloak/

ENV KC_DB=fxurerfl
ENV KC_DB_USERNAME=fxurerfl
ENV KC_DB_PASSWORD=R71RUnRpWkRRvpm3gbhUJ3pmlq2FVKoP
ENV KC_HOSTNAME=flora.db.elephantsql.com

ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]

EXPOSE 8000
CMD ["/opt/main"]