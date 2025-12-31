# Description

Automation for starting/stopping Docker Compose stacks on an Asustor ADM NAS.
Runs any Docker Compose stack found in `/volmain/.@docker_compose/stacks/`.

## Usage

Run under root user in ADM:

```bash
curl -sSL https://raw.githubusercontent.com/kreigan/adm-docker-loader/main/setup.sh | sh
```

## Configuration

Edit the `/volmain/.@docker_compose/.env` and `/volmain/.@docker_compose/config.conf` files to customize behavior.
`.env` file contains environment variables for Docker Compose stacks. The loader script
reads this file for each `up` and `down` operation. 

`config.conf` allows to modify configuration options for the loader script.
