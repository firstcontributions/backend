# FirstContributions App Backend

## Overview

This repository contains the backend service for the First Contributions app. It handles API requests and supports integration with the frontend.

## Setup & Run

```sh
sudo make configure
make run
```

## macOS / Docker Setup Notes

On macOS and environments using Docker Compose, the application may not be able to connect to Docker network IPs because Docker runs inside a virtual machine.

To resolve this, set up a reverse proxy using Nginx.

### Steps

1. Install Nginx and refer to the configuration file:
   `./mac.nginx.conf`

2. Update your `/etc/hosts` file:

```text
127.0.0.1 api.firstcontributions.com
127.0.0.1 explorer.firstcontributions.com
```

## Running Integration Tests

```sh
sudo make configure
make itest
```

