FROM scratch
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

# Add the binary
ADD dist/caddy-docker /usr/local/bin/

# Expose ports 80 and 443
EXPOSE 80 443

# No arguments are needed for running the app
CMD "caddy-docker"
