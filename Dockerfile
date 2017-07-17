FROM scratch
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

# Add the binary
ADD dist/caddy-docker /usr/local/bin/

# Add the root CAs
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/

# Expose ports 80 and 443
EXPOSE 80 443

# Create a volume for the TLS files
VOLUME /var/lib/caddy-docker

# Tell Caddy to use the volume
ENV CADDYPATH=/var/lib/caddy-docker

# No arguments are needed for running the app
ENTRYPOINT ["/usr/local/bin/caddy-docker"]
