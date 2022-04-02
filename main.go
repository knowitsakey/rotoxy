package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"

	//"context"
	"fmt"
	"github.com/gtuk/rotating-tor-proxy/core"
	"github.com/urfave/cli/v2"
	//"golang.org/x/sync/errgroup"
	"log"
	"os"
	"strconv"
)

func main() {
	var numberTorInstances int
	var port int
	var circuitInterval int
	var hslist string

	app := &cli.App{
		Name:  "rotoxy",
		Usage: "run a rotating Tor proxy server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "tors",
				Value:       1,
				Usage:       "number of Tor proxies that should run",
				Destination: &numberTorInstances,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       8080,
				Usage:       "port where the reverse proxy should listen on",
				Destination: &port,
			},
			&cli.IntFlag{
				Name:        "circuitInterval",
				Value:       30,
				Usage:       "number in seconds after a new circuit should be requested",
				Destination: &circuitInterval,
			},
			&cli.StringFlag{
				Name:        "hslist",
				Value:       "onionservices.txt",
				Usage:       "number in seconds after a new circuit should be requested",
				Destination: &hslist,
			},
		},
		Action: func(c *cli.Context) error {
			//return run(port, circuitInterval, hslist)
			return dash(port, circuitInterval, hslist)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func readtodict(filename string, txtlines *[]string) {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	//var txtlines []string

	for scanner.Scan() {
		*txtlines = append(*txtlines, scanner.Text())
	}

	file.Close()

}

func dash(port int, circuitInterval int, hslist string) error {

	var textlines []string
	readtodict(hslist, &textlines)
	numberTorInstances := len(textlines)

	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]*core.TorProxy1, numberTorInstances)

	for _, eachline := range textlines {
		fmt.Println(eachline)
	}
	//wg := new(sync.WaitGroup)
	//wg.Add(numberTorInstances - 1)
	var err error
	for i := 0; i < numberTorInstances; i++ {
		//proxies[i], err = core.CreateSimpleSshProxy(circuitInterval, textlines[i])
		proxies[i], err = core.CreateTorProxy2(circuitInterval, textlines[i])
	}
	DoubleProxies := make([]*core.SshProxy, numberTorInstances)

	for i := 0; i < numberTorInstances; i++ {
		DoubleProxies[i], err = core.CreateSshProxy(textlines[i], proxies[i])

	}

	fmt.Println(err)
	//startServers(proxies)
	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	reverseProxy1 := &core.ReverseProxy1{}

	//http:
	
	err = reverseProxy1.Start1(proxies, port)

	if err != nil {
		return err
	}

	//defer core.CloseProxies1(proxies)
	return nil
}

//i user, hidden service, ssh private key, torsocks5port, doubleproxyport,
//r nothing
//m
//e test 2 see if it works n give us a socks5 proxy.
//execute command launching darkssh in a new process

func start2ndProxy() {

}

func startServers(proxies []*core.TorProxy1) {
	var err error
	for i := 0; i < len(proxies); i++ {
		//go func() {
		if proxies[i] != nil {
			//go func() {

			socks5Address := "127.0.0.1:" + strconv.Itoa(*proxies[i].DoubleProxyPort)
			fmt.Println("we trine kreate socks serva at "+socks5Address, err)
			go func() {
				err := proxies[i].Socks5s.ListenAndServe("tcp", socks5Address)
				if err != nil {
					fmt.Println("failed to create socks5 server", err)
				}
			}()
			if err != nil {
				fmt.Println("failed to create socks5 server", err)
				//return err
				//return
			}
			fmt.Println("kreated u a socks serva at "+socks5Address, err)
			//			wg2.Done()
			//		}()
		}

	}
}

func testSocks5(port int) {

	const (
		NoAuth          = uint8(0)
		noAcceptable    = uint8(255)
		UserPassAuth    = uint8(2)
		userAuthVersion = uint8(1)
		authSuccess     = uint8(0)
		authFailure     = uint8(1)
		socks5Version   = uint8(5)
	)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {

	}
	lAddr := l.Addr().(*net.TCPAddr)

	socks5Address := "127.0.0.1:" + strconv.Itoa(port)
	// Get a local conn
	conn, err := net.Dial("tcp", socks5Address)
	if err != nil {
		log.Println(conn)
	}
	req := bytes.NewBuffer(nil)
	req.Write([]byte{5})
	req.Write([]byte{2, NoAuth, UserPassAuth})
	req.Write([]byte{1, 3, 'f', 'o', 'o', 3, 'b', 'a', 'r'})
	req.Write([]byte{5, 1, 0, 1, 127, 0, 0, 1})

	port1 := []byte{0, 0}
	binary.BigEndian.PutUint16(port1, uint16(lAddr.Port))
	req.Write(port1)

	// Send a ping
	req.Write([]byte("ping"))

	// Send all the bytes
	conn.Write(req.Bytes())

	// Verify response
	expected := []byte{
		socks5Version, UserPassAuth,
		1, authSuccess,
		5,
		0,
		0,
		1,
		127, 0, 0, 1,
		0, 0,
		'p', 'o', 'n', 'g',
	}
	out := make([]byte, len(expected))

	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadAtLeast(conn, out, len(out)); err != nil {
		fmt.Println(err)
	}

	// Ignore the port
	out[12] = 0
	out[13] = 0

	if !bytes.Equal(out, expected) {
		fmt.Println(err)
	}

}
