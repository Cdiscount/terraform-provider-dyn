Terraform Provider
==================

A Terraform provider for Dynect forked from [official archived terraform provider](https://github.com/hashicorp/terraform-provider-dyn).

It support Dynect record creation, like the original project, and it add support for traffic director creation.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) (works with 0.15.x)
-	[Go](https://golang.org/doc/install) (works with 1.16) (to build the provider plugin)
-   make


Developing the Provider
---------------------------

To compile the provider, run `make build`. This will build the provider.

```sh
$ make build
$ ./terraform-provider-dyn
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to install the provider locally, run `make install`.
