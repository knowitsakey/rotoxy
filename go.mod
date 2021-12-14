module github.com/gtuk/rotating-tor-proxy

go 1.16

require (
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/cretz/bine v0.2.0
	github.com/eahydra/socks v0.0.0-20191219154456-071591e7aaf0
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)

replace github.com/eyedeekay/darkssh => /home/based/darkssh
