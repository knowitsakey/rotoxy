package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gtuk/rotating-tor-proxy/core"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
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

	var txtlines []string
	readtodict(hslist, &txtlines)
	numberTorInstances := len(txtlines)

	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]core.TorProxy1, 0)
	ch := make(chan core.TorProxy1, numberTorInstances)
	//var proxyports []int
	//proxies1 = []int{9002, 9008}

	g, _ := errgroup.WithContext(context.Background())

	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}
	//var wg sync.WaitGroup
	//err := wg.Add(numberTorInstances)
	for i := 0; i < numberTorInstances-1; i++ {
		//i := 1
		//for _, hsaddr := range txtlines {
		//go func(i int) {
		g.Go(func() error {
			fmt.Println("Starting proxy number " + strconv.Itoa(i) + " at " + txtlines[i])
			torProxy, err := core.CreateTorProxy1(circuitInterval, txtlines[i])
			if err != nil {
				if torProxy != nil {
					torProxy.Close1()
				}
				return err
			}
			//port1 := torProxy.DoubleProxyPort
			//proxyports = append(proxyports, *port1)

			ch <- *torProxy
			return nil
			//}(i)
		})
	}
	err := g.Wait()
	close(ch)
	for proxy := range ch {
		proxies = append(proxies, proxy)
	}
	//return nil
	/*	if err := torProxy1.socks5s.ListenAndServe("tcp", socks5Address); err != nil {
			fmt.Println("failed to create socks5 server", err)
		}
		fmt.Println("kreated u a socks serva at "+socks5Address, err)

		if err != nil {
			return nil, err
		}*/
	defer core.CloseProxies1(proxies)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	reverseProxy1 := &core.ReverseProxy1{}

	err = reverseProxy1.Start2(proxies, port)
	if err != nil {
		return err
	}

	return err
}
