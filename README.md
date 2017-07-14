## caddy-docker

Do you find yourself constantly adding and removing Docker services? Are you tired of constantly creating and removing configuration files for your load balancer? caddy-docker is here to help.

By creating your containers with special tags, caddy-docker will automatically generate the appropriate configuration files for [Caddy](https://github.com/mholt/caddy) and reload the application. And because Caddy uses Let's Encrypt for free TLS certificates, all you have to do in order to have a TLS-enabled service up and running is merely start the container!
