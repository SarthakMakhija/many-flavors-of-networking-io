### Many flavors of networking IO

[![build](https:github.com/SarthakMakhija/many-flavors-of-networking-io/actions/workflows/build.yml/badge.svg)](https:github.com/SarthakMakhija/many-flavors-of-networking-io/actions/workflows/build.yml)

This repository is a reference implementation for the article titled "Many flavors of networking IO". 
It has the following implementations of TCP servers:

1. **Single Thread Blocking IO**

`TCPServer` implements "Single thread blocking IO" pattern. The implementation of `TCPServer`:

- runs a continuous loop in a single goroutine (/main goroutine).
- a new instance of `IncomingTCPConnection` is created for every new connection.
- The incoming TCP connection is handled in the same main goroutine.
- This pattern involves **blocking IO** to read from the incoming connection.


2. **Multi Thread Blocking IO**

`TCPServer` implements "Multi thread blocking IO" pattern. The implementation of `TCPServer`:

- runs a continuous loop in a single goroutine (/main goroutine).
- a new instance of IncomingTCPConnection is created for every new connection.
- The incoming TCP connection is handled in new goroutine.
- This pattern involves **goroutine per connection** and **blocking IO** to read from the incoming connection.

3. **Non-blocking with Busy Wait**

`TCPServer` implements "Non-Blocking with Busy-Wait" pattern. The implementation of `TCPServer`:

- runs a continuous loop in a single goroutine (/main goroutine).
- it marks the server file descriptor non-blocking, this means any IO operations on this file descriptor will not block. However, the file descriptor can be polled.
- an incoming connection is represented by its own file descriptor: `connectionFd`.
- `connectionFd` is also marked non-blocking.
- a new client is created (for the incoming `connectionFd`) which handles the connection by performing **busy-wait or polling**.
- all the IO operations are **non-blocking** .

4. **Single Thread Event loop** (using `KQueue`)

`TCPServer` implements "Single thread Non-Blocking with event loop" pattern. It starts an event loop which:

- runs in its own goroutine.
- polls the `KQueue` for events on the subscribed file descriptors.
- if the polled event's file descriptor is same as the server's file descriptor: a new client is accepted,
- else: an existing client for the file descriptor is run.
- all the IO operations are **non-blocking** .

*The article is yet to be written.
