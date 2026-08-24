# The broker for a deployed environment.
#
# Compose can mount the config and the entrypoint straight off the host, so it
# runs the stock eclipse-mosquitto image. Railway builds from the repo and has
# no host to mount from, so the same two files are baked in here instead. Keep
# this in step with the volumes in docker-compose.yml.
FROM eclipse-mosquitto:2

COPY docker/mosquitto/mosquitto.conf /mosquitto/config/mosquitto.conf
COPY docker/mosquitto/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Dynamic security bootstraps on first boot and is written next to the config,
# so both the config and the persistence directory have to survive a restart.
VOLUME ["/mosquitto/config", "/mosquitto/data"]

EXPOSE 1883

ENTRYPOINT ["/entrypoint.sh"]
