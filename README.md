## caddy-docker

[![GoDoc](https://godoc.org/github.com/nathan-osman/caddy-docker?status.svg)](https://godoc.org/github.com/nathan-osman/caddy-docker)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

Do you find yourself constantly adding and removing Docker services? Are you tired of constantly creating and removing configuration files for your load balancer? caddy-docker is here to help.

By creating your containers with special tags, caddy-docker will automatically generate the appropriate configuration files for [Caddy](https://github.com/mholt/caddy) and reload the application. And because Caddy uses Let's Encrypt for free TLS certificates, all you have to do in order to have a TLS-enabled service up and running is merely start the container!
