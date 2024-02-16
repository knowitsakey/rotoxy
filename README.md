This is a proxy server using multiple instances of tor in order to increase bandwidth for html requests over the tor network. Html requests are split between multiple proxies in a round-robin approach
Traffic is routed through one or more hidden services such that outgoing and incoming traffic appears to be coming from the remote ip address
The client's connection follows the below path:

client <-> tor entry node <-> hidden service <-> vps <-> server
