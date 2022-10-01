# FirstContributions App Backend

How to run

```sh
$ sudo make configure
$ make run
```


For mac and environments using docker compose, it is unable to connect to the docker network ip as docker is running in a VM.
Need to setup a reverse proxy to do local setup

1. Install nginx, refer config (./mac.nginx.conf)
2. change hosts config (/etc/hosts)

```
127.0.0.1 api.firstcontributions.com
127.0.0.1 explorer.firstcontributions.com
```
To run integration tests

```sh
$ sudo make configure
$ make itest
```
