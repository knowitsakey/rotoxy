This is a proxy server using multiple instances of tor in order to increase bandwidth for html requests over the tor network.
Traffic is routed through a hidden service such that outgoing and incoming traffic appears to be coming from the remote ip address
The client's connection follows this path:

client <-> tor entry node <-> hidden service <-> vps <-> server
