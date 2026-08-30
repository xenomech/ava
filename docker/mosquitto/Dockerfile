# Deployed broker: config baked in because Railway has no host to mount from, and no VOLUME because Railway declares its own.
FROM eclipse-mosquitto:2

COPY docker/mosquitto/mosquitto.conf /mosquitto/config/mosquitto.conf
COPY docker/mosquitto/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 1883

ENTRYPOINT ["/entrypoint.sh"]
