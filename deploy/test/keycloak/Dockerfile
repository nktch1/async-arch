FROM quay.io/keycloak/keycloak:24.0.1
COPY providers /opt/keycloak/providers
RUN /opt/keycloak/bin/kc.sh build

ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start-dev"]