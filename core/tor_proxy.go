package core

import (
	"context"
	"fmt"
	"github.com/cretz/bine/tor"
	"github.com/eahydra/socks"

	"net/http"
	"strconv"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"
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
	Onionproxy      *socks.Socks5Server
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
		ExtraArgs:       extraArgs,
		NoAutoSocksPort: true,
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
	conf1.ProxyAddress = hsaddr
	conf1.ProxyAddress = "fefix3iwkb5b3b2z2sicik7re2qsv2o5hrch7pyvuvifklou2fnblayd.onion:80"
	conf1.Forward = proxy.Direct
	// Make connection
	dialer1, err := torCtx.Dialer(ctx, conf1)
	myport := 10000
	//here you need to create a socks 5 proxy listening on 127.0.0.1:myport, by connecting to a hidden service
	//R socks5 proxy object, listening at myport
	//M the torproxy1 struct
	//E what should socksproxy object be: a dialler or a proxyobhect?

	dialer2, err := proxy.SOCKS5("tcp", "fefix3iwkb5b3b2z2sicik7re2qsv2o5hrch7pyvuvifklou2fnblayd.onion:80", nil, dialer1)
	//myconn := dialer2.Dial("tcp", "https://startpage.com")

	httpClient := &http.Client{Transport: &http.Transport{Dial: dialer2.Dial}}
	// Get /
	resp, err := httpClient.Get("https://check.torproject.org")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Grab the <title>
	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Title: %v\n", getTitle(parsed))
	//fmt.Printf(parsed)
	fmt.Print(parsed)
	torProxy1.DoubleProxyPort = &myport

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
