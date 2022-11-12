# tests

- [jmeter](./jmeter/) : Jmeter testcases for REST APIs.
- [shell](./shell/) : Shell testcases for `mantleNode & modules`.
- [curl](./curl/) : Curl script testcases for REST APIs.
- [scripts](./scripts/) : Scripts folders for setting up node & modules.

# Setting Up Node

- make sure you installed go & `assetNode & assetClient`.
- after installing, please setup systemd services:
  - [mantle-node.service](./scripts/mantle-node.service) : mantle-node service file, investigate this file and edit accordingly.
  - [mantle-client.service](./scripts/mantle-client.service) : mantle-client service file, investigate this file and edit accordingly.
- run `systemctl enable mantle-client && systemctl enable mantle-node` 
- now run [service.sh](./scripts/service.sh)
