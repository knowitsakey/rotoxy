package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/cretz/bine/tor"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type TorProxy struct {
	Ctx             *tor.Tor
	ControlPort     *int
	ProxyPort       *int
	CircuitInterval *int
}

type TorProxy1 struct {
	Ctx             *tor.Tor
	Client          *ssh.Client
	Socks5s         *socks5.Server
	ControlPort     *int
	ProxyPort       *int
	CircuitInterval *int
	DoubleProxyPort *int
	//Onionproxy      *socks5.Server
}

func CreateSimpleSshProxy(circuitinterval int, hsaddr string) (*TorProxy1, error) {

	var err error
	torProxy1 := &TorProxy1{}
	//ctx := context.Background()
	//torCtx, err := tor.Start(ctx, &tor.StartConf{
	//	ExtraArgs: extraArgs,
	//	//NoAutoSocksPort: true,
	//	EnableNetwork: true,
	//})

	sshConf := &ssh.ClientConfig{
		User:            "based",
		Auth:            []ssh.AuthMethod{ssh.Password("lab")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//tp, err := NewTorGate(conf1.ProxyAddress)

	if err != nil {
		return nil, err
	}
	//conn1, err := tp.DialTor(hsaddr + ":22")
	torProxy1.Client, err = ssh.Dial("tcp", "127.0.0.1:22", sshConf)
	if err != nil {
		return nil, err
	}

	//d, err := torCtx.Dialer(ctx, conf1)
	//conn1, err := d.DialContext(ctx, "tcp", hsaddr)
	//if err != nil {
	//	return nil, err
	//}
	//c, chans, reqs, err := ssh.NewClientConn(conn1, "127.0.0.1:22", sshConf)
	//if err != nil {
	//	return nil, err
	//}
	////torProxy1.client, err := &ssh.NewClient(c, chans, reqs)
	////client1 := &ssh.NewClient(c, chans, reqs
	//torProxy1.Client = ssh.NewClient(c, chans, reqs)
	////client1 := ssh.NewClient(c, chans, reqs)
	//fmt.Println("Connected to .onion successfully!")
	//
	//defer client1.Close()

	fmt.Println("connected to ssh server")

	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return torProxy1.Client.Dial(network, addr)
		},
	}

	torProxy1.Socks5s, err = socks5.New(conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	port, err := GetFreePort()
	torProxy1.DoubleProxyPort = &port
	socks5Address := "127.0.0.1:" + strconv.Itoa(port)
	fmt.Println("we try kreate socks5 server at "+socks5Address, err)
	go func() {
		if err := torProxy1.Socks5s.ListenAndServe("tcp", socks5Address); err != nil {
			fmt.Println("failed to create socks5 server", err)
		}
	}()
	fmt.Println("kreated u a socks serva at "+socks5Address, err)

	if err != nil {
		return nil, err
	}
	//	torProxy.
	return torProxy1, nil

}

//func launchSSHServ(conn net.Conn) (s *socks5.Server) {
//
//}

func feedConn(torProxy1 TorProxy1, hsaddr string) (*TorProxy1, error) {
	var err error
	socks5Address := "127.0.0.1:" + strconv.Itoa(*torProxy1.DoubleProxyPort)

	sshConf := &ssh.ClientConfig{
		User:            "based",
		Auth:            []ssh.AuthMethod{ssh.Password("lab")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	proxyAddress := "127.0.0.1:" + strconv.Itoa(*torProxy1.ProxyPort)
	tp, err := NewTorGate(proxyAddress)

	if err != nil {
		return nil, err
	}
	conn1, err := tp.DialTor(hsaddr + ":22")
	//d, err := torCtx.Dialer(ctx, conf1)
	//conn1, err := d.DialContext(ctx, "tcp", hsaddr)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn1, proxyAddress, sshConf)
	if err != nil {
		return nil, err
	}
	//torProxy1.client, err := &ssh.NewClient(c, chans, reqs)
	//client1 := &ssh.NewClient(c, chans, reqs
	torProxy1.Client = ssh.NewClient(c, chans, reqs)
	//client1 := ssh.NewClient(c, chans, reqs)
	fmt.Println("Connected to .onion successfully!")

	//defer client1.Close()

	fmt.Println("connected to ssh server")
	fmt.Println("we trine kreate socks serva at"+socks5Address, err)
	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return torProxy1.Client.Dial(network, addr)
		},
	}

	torProxy1.Socks5s, err = socks5.New(conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//if err := torProxy1.socks5s.ListenAndServe("tcp", socks5Address); err != nil {
	//	fmt.Println("failed to create socks5 server", err)
	//}
	//fmt.Println("kreated u a socks serva at "+socks5Address, err)

	if err != nil {
		return nil, err
	}
	//	torProxy.
	return nil, err
}

func CreateTorProxy2(circuitInterval int, hsaddr string) (*TorProxy1, error) {
	torProxy1 := &TorProxy1{}
	ctx := context.Background()

	port, err := GetFreePort()
	if err != nil {
		return nil, err
	}

	var extraArgs []string
	// Set socks port
	extraArgs = append(extraArgs, "--SocksPort")
	extraArgs = append(extraArgs, strconv.Itoa(port))

	// Set new circuit interval after circuit was used once
	extraArgs = append(extraArgs, "--MaxCircuitDirtiness")
	extraArgs = append(extraArgs, strconv.Itoa(circuitInterval))

	torCtx, err := tor.Start(ctx, &tor.StartConf{
		ExtraArgs: extraArgs,
		//TorrcFile: "torrc",
		//NoAutoSocksPort: true,
		EnableNetwork: true,
	})
	if err != nil {
		return nil, err
	}

	torProxy1.Ctx = torCtx
	torProxy1.ProxyPort = &port
	torProxy1.ControlPort = &torCtx.ControlPort
	torProxy1.CircuitInterval = &circuitInterval
	conf1 := &tor.DialConf{}

	conf1.ProxyAddress = "127.0.0.1:" + strconv.Itoa(*torProxy1.ProxyPort)
	//conf1.ProxyAddress = "127.0.0.1:9050"
	conf1.Forward = proxy.Direct
	doubleproxyport, err := GetFreePort()
	torProxy1.DoubleProxyPort = &doubleproxyport
	return torProxy1, nil
}

func CreateTorProxy1(circuitInterval int, hsaddr string) (*TorProxy1, error) {
	torProxy1 := &TorProxy1{}
	ctx := context.Background()

	port, err := GetFreePort()
	if err != nil {
		return nil, err
	}

	var extraArgs []string
	// Set socks port
	extraArgs = append(extraArgs, "--SocksPort")
	extraArgs = append(extraArgs, strconv.Itoa(port))

	// Set new circuit interval after circuit was used once
	extraArgs = append(extraArgs, "--MaxCircuitDirtiness")
	extraArgs = append(extraArgs, strconv.Itoa(circuitInterval))

	torCtx, err := tor.Start(ctx, &tor.StartConf{
		ExtraArgs: extraArgs,
		//NoAutoSocksPort: true,
		EnableNetwork: true,
	})
	if err != nil {
		return nil, err
	}

	torProxy1.Ctx = torCtx
	torProxy1.ProxyPort = &port
	torProxy1.ControlPort = &torCtx.ControlPort
	torProxy1.CircuitInterval = &circuitInterval
	conf1 := &tor.DialConf{}

	conf1.ProxyAddress = "127.0.0.1:" + strconv.Itoa(*torProxy1.ProxyPort)
	//conf1.ProxyAddress = "127.0.0.1:9050"
	conf1.Forward = proxy.Direct
	doubleproxyport, err := GetFreePort()
	torProxy1.DoubleProxyPort = &doubleproxyport
	socks5Address := "127.0.0.1:" + strconv.Itoa(*torProxy1.DoubleProxyPort)

	sshConf := &ssh.ClientConfig{
		User:            "based",
		Auth:            []ssh.AuthMethod{ssh.Password("lab")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tp, err := NewTorGate(conf1.ProxyAddress)

	if err != nil {
		return nil, err
	}
	conn1, err := tp.DialTor(hsaddr + ":22")
	//d, err := torCtx.Dialer(ctx, conf1)
	//conn1, err := d.DialContext(ctx, "tcp", hsaddr)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn1, conf1.ProxyAddress, sshConf)
	if err != nil {
		return nil, err
	}
	//torProxy1.client, err := &ssh.NewClient(c, chans, reqs)
	//client1 := &ssh.NewClient(c, chans, reqs
	torProxy1.Client = ssh.NewClient(c, chans, reqs)
	//client1 := ssh.NewClient(c, chans, reqs)
	fmt.Println("Connected to .onion successfully!")

	//defer client1.Close()

	fmt.Println("connected to ssh server")
	fmt.Println("we trine kreate socks serva at"+socks5Address, err)
	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return torProxy1.Client.Dial(network, addr)
		},
	}

	torProxy1.Socks5s, err = socks5.New(conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	torProxy1.Socks5s.ListenAndServe("tcp", socks5Address)
	//if err := torProxy1.socks5s.ListenAndServe("tcp", socks5Address); err != nil {
	//	fmt.Println("failed to create socks5 server", err)
	//}
	//fmt.Println("kreated u a socks serva at "+socks5Address, err)

	if err != nil {
		return nil, err
	}
	//	torProxy.
	return torProxy1, nil
}

func (m *TorProxy) Close() {
	if m.Ctx != nil {
		m.Ctx.Close()
	}
}
func (m *TorProxy1) Close1() {
	if m.Ctx != nil {
		m.Ctx.Close()
	}
}

type TorGate string

// TOR_GATE string constant with localhost's Tor port
const TOR_GATE_ string = "127.0.0.1:9050"

func NewTorGate(torgate string) (*TorGate, error) {
	//torgate = TOR_GATE_
	duration, _ := time.ParseDuration("10s")
	connect, err := net.DialTimeout("tcp4", torgate, duration)

	if err != nil {
		return nil, errors.New("Could not test TOR_GATE_: " + err.Error())
	}

	// Tor proxies reply to anything that looks like
	// HTTP GET or POST with known error message.
	connect.Write([]byte("GET /\n"))
	connect.SetReadDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 4096)

	for {
		n, err := connect.Read(buf)

		if err != nil {
			return nil, errors.New("It is not TOR_GATE_")
		}

		if bytes.Contains(buf[:n], []byte("Tor is not an HTTP Proxy")) {
			connect.Close()
			gate := TorGate(torgate)

			return &gate, nil
		}
	}
}

// DialTor dials to the .onion address
func (gate *TorGate) DialTor(address string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5("tcp4", string(*gate), nil, proxy.Direct)

	if err != nil {
		return nil, errors.New("Could not connect to TOR_GATE_: " + err.Error())
	}

	connect, err := dialer.Dial("tcp", address)

	if err != nil {
		return nil, errors.New("Failed to connect: " + err.Error())
	}

	return connect, nil

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}
