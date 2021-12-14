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
	ControlPort     *int
	ProxyPort       *int
	CircuitInterval *int
	DoubleProxyPort *int
	//Onionproxy      *socks5.Server
}

func CreateTorProxy(circuitInterval int, hsaddr string) (*TorProxy, error) {
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
		ExtraArgs:       extraArgs,
		NoAutoSocksPort: true,
	})

	if err != nil {
		return nil, err
	}

	torProxy := &TorProxy{}
	torProxy.Ctx = torCtx
	torProxy.ProxyPort = &port
	torProxy.ControlPort = &torCtx.ControlPort
	torProxy.CircuitInterval = &circuitInterval
	conf1 := &tor.DialConf{}
	conf1.ProxyAddress = hsaddr
	conf1.ProxyAddress = "fefix3iwkb5b3b2z2sicik7re2qsv2o5hrch7pyvuvifklou2fnblayd.onion:80"
	conf1.Forward = proxy.Direct
	// Make connection
	_, err = torCtx.Dialer(ctx, conf1)
	if err != nil {
		return nil, err
	}
	//	torProxy.
	return torProxy, nil
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
func CreateTorProxy1(circuitInterval int, hsaddr string) (*TorProxy1, error) {
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

	torProxy1 := &TorProxy1{}
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

	client1 := ssh.NewClient(c, chans, reqs)
	fmt.Println("Connected to .onion successfully!")

	defer client1.Close()

	fmt.Println("connected to ssh server")
	fmt.Println("we trine kreate socks serva at"+socks5Address, err)
	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return client1.Dial(network, addr)
		},
	}

	serverSocks, err := socks5.New(conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := serverSocks.ListenAndServe("tcp", socks5Address); err != nil {
		fmt.Println("failed to create socks5 server", err)
	}
	fmt.Println("kreated u a socks serva at "+socks5Address, err)

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

	connect, err := dialer.Dial("tcp4", address)

	if err != nil {
		return nil, errors.New("Failed to connect: " + err.Error())
	}

	return connect, nil

}
