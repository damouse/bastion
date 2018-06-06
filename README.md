# Go Bastion

An SSH tunnel proxy server that manages the connections in go instead of managing the connections on the host. 

Setup: 

```
go get 
```


Setting up a reverse proxy by hand (two terminal windows): 

```
$ ssh -R 19999:localhost:22 ubuntu@lhr-vpn
$ ssh damouse@localhost -p 19999
```


Links:
- https://www.howtoforge.com/reverse-ssh-tunneling
- [Keep alive](http://mirko.windhoff.net/how_to/make_a_reverse_ssh_tunnel)

## Sketchpad

Current intended function is running `proxy.go` as a docker container with a well-known IP address, have devices set up a reverse proxy to it, then ssh to the host with a command like this:

```
ssh localhost -p 2222 -l foo
```

where the passed username (foo) represents a device serial number or name. The only missing piece of the puzzle is mapping device names to their reverse-proxied ports, including informing the server that mappings change. If the server supported reverse proxying out of the box it would be easy, but I don't know that I can do that. 
