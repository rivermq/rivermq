[![Build Status](https://travis-ci.org/rivermq/rivermq.svg?branch=master)](https://travis-ci.org/rivermq/rivermq) [![codecov](https://codecov.io/gh/rivermq/rivermq/branch/master/graph/badge.svg)](https://codecov.io/gh/rivermq/rivermq)


RiverMQ
========

WebHook based messaging system

RiverMQ will provide a WebHook based asynchronous messaging solution for distributed applications.

Clients will register subscriptions to a message "type" with a RiverMQ Server instance via a HTTP POST.  This subscription will include a "callback url" which the client expects to receive messages on.  When a separate client sends a message of a matching "type" to a RiverMQ server instance, that message will be sent, via a HTTP POST, to the "callback url" of all subscribed clients.  Based on the response code received by RiverMQ, message redelivery attempts will be retried for up to 30 mintes.  After 30 minutes RiverMQ will cease attempts to send the message to the client.

RiverMQ is inspired by [subpub](https://github.com/PearsonEducation/subpub)


Goals
-----

1. Provide a single executable with minimal configuration
1. Simple horizontal scaling
1. Automatic discovery and configuration of RiverMQ nodes via [Serf](https://www.serf.io/)
1. Secured communication between RiverMQ nodes with [ZeroMQ](http://zeromq.org/)
1. Message and Subscription storage with [MongoDB](https://mongodb.org)
1. Administration UI with [Angular.js](https://angularjs.org/), and some charting library.
1. Metrics with [Prometheus](https://prometheus.io/)
1. Allow data storage to be customizable with a plugin solution similar in design to [Docker Plugins](https://docs.docker.com/engine/extend/plugin_api/)



Development
-----------
Clone the repo and inspect the Makefile for running tests and building.

Test coverage result file concatination is done using [gover](https://github.com/modocache/gover).  This is a required dependency which must be installed using the following:
```bash
go get github.com/modocache/gover
```

Developed using [Atom](https://atom.io/) [configured for Go development](http://marcio.io/2015/07/supercharging-atom-editor-for-go-development).
