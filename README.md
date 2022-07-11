terraform-provider-leaseweb
---------------------------

A leaseweb terraform provider


setup for development
---------------------

This setup uses docker so you do not need go (or any of the build tools) on
your workstation.

1. You need `docker` and `docker compose`.
2. Git clone this repository and `cd` into it.
3. Run `docker compose build`
4. Run `docker compose up -d`


building the plugin
-------------------

To build the plugin (in docker):

    docker compose exec --env GOOS=$GOOS --env GOARCH=$GOARCH backend go build -o terraform-provider-leaseweb

Now you can move the plugin into the `~/.terraform.d/plugins/` directorym (see
the `Makefile` for details) and you are ready to go.
