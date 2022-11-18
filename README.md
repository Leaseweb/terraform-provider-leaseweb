Terraform Leaseweb Provider
---------------------------

A Terraform provider to manage Leaseweb resources.


Requirements
------------

Terraform 0.12.0 or later is needed to use this plugin.

Go 1.18 or later is needed to build the plugin.


Setup for development
---------------------

This setup uses docker so you do not need go (or any of the build tools) on
your workstation.
1. You need `docker` and `docker compose`.
2. Git clone this repository and `cd` into it.
3. Run `docker compose build`
4. Run `docker compose up -d`


Building the plugin
-------------------

To build the plugin (in docker):

    docker compose exec --env GOOS=$GOOS --env GOARCH=$GOARCH backend go build -o terraform-provider-leaseweb


Using the plugin locally
------------------------

Move the plugin into the `~/.terraform.d/plugins/terraform.local/local/leaseweb/$VERSION/$GOOS_$GOARCH/` directory.

It can then be used the following way in your config:

    terraform {
      required_providers {
        leaseweb = {
          source = "terraform.local/local/leaseweb"
          version = "0.1.2"
        }
      }
    }
