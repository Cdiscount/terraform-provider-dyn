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
Terraform documentation can be generated using `make docs`.


Terraform resources are in `dyn/`.
Originally the provider uses `go-dynect` package, but it misses some features for our needs. As this is an old project, we just embbeded it in this repository, but it may be moved in it's own repository again in the future.

### About test coverage

This fork do not add tests, every resource has been tested manually.

To our knowledge, there is not dynect sandbox API to make good integration tests. We only kept original acceptance tests.

A good improvement could be to add a mock HTTP server to add tests to this provider.
