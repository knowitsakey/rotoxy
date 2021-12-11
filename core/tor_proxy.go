package core

import (
	"context"
	"github.com/cretz/bine/tor"
	"strconv"
)

type TorProxy struct {
	Ctx             *tor.Tor
	ControlPort     *int
	ProxyPort       *int
	CircuitInterval *int
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
	//	conf1.ProxyAddress = "fefix3iwkb5b3b2z2sicik7re2qsv2o5hrch7pyvuvifklou2fnblayd.onion"

	// Make connection
	_, err = torCtx.Dialer(ctx, conf1)
	if err != nil {
		return nil, err
	}

	return torProxy, nil
}

func (m *TorProxy) Close() {
	if m.Ctx != nil {
		m.Ctx.Close()
	}
}
