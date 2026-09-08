# Deployed broker: config baked in because Railway has no host to mount from, and no VOLUME because Railway declares its own.
FROM eclipse-mosquitto:2

COPY docker/mosquitto/mosquitto.conf /mosquitto/config/mosquitto.conf
COPY docker/mosquitto/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# 1883 plaintext for a private network, 8883 native TLS, 9001 websockets behind a TLS-terminating edge.
EXPOSE 1883 8883 9001

ENTRYPOINT ["/entrypoint.sh"]
