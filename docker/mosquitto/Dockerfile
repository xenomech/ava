# The broker for a deployed environment.
#
# Compose can mount the config and the entrypoint straight off the host, so it
# runs the stock eclipse-mosquitto image. Railway builds from the repo and has
# no host to mount from, so the same two files are baked in here instead. Keep
# this in step with the volumes in docker-compose.yml.
#
# No VOLUME instruction. Railway attaches one volume per service and declares it
# in .railway/railway.ts; an anonymous Dockerfile volume it cannot satisfy makes
# the deploy fail before it produces a single line of build log.
FROM eclipse-mosquitto:2

COPY docker/mosquitto/mosquitto.conf /mosquitto/config/mosquitto.conf
COPY docker/mosquitto/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 1883

ENTRYPOINT ["/entrypoint.sh"]
