# go-dynect

This is a fork of [go-dynect](https://github.com/nesv/go-dynect) 0.6.0 that add some new types to request dynect API for DSF resources.

It also adds a mutex to prevent client to make concurrent requests using the same session, as this is not supported by the Dynect API.
