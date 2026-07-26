ARG TARGETPLATFORM

FROM scratch
ARG TARGETPLATFORM
WORKDIR /opt/spikeball-league
ENTRYPOINT ["/opt/spikeball-league/service"]
EXPOSE 8080
COPY $TARGETPLATFORM/spikeball-league /opt/spikeball-league/service
