# Rotoxy NG

A rotating tor proxy service that leverages darkssh to start a number of connections to tor hidden services, exposing them under a single proxy. Ideal usage is putting tor traffic through a web browser using proxychains.
The purpose of this project is to evade web traffic fingerprinting attacks at the client.
```bash
http => proxychains => ssh client => tor => ssh server => vpn client => vpn server
```
Note: this project is experimental, use at your own risk.

the vpn at the other end can be set up using the following iptables rules: 



### Prerequisites
In order to use the tool you need have Tor installed on the machine

### Usage
Download the latest release from github
open three terminals, set http	127.0.0.1 8080 in proxychains.conf
```bash
./rotating-tor-proxy
./proxychains4-daemon
./proxychains4 -f /home/based/proxychains-ng/src/proxychains.conf firefox
```
### Docker
```bash
docker run -p 8080:8080 gtuk/rotoxy:latest # Run with default parameters
docker run -p 8088:8088 gtuk/rotoxy:latest --tors 1 --port 8080 --circuitInterval 30 # Run with custom parameters
```

### TODOS
* Tests
* Better documentation
