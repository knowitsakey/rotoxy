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
func run(port int, circuitInterval int, hslist string) error {

	var txtlines []string
	readtodict(hslist, &txtlines)
	numberTorInstances := len(txtlines)

	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]core.TorProxy, 0)
	ch := make(chan core.TorProxy, numberTorInstances)

	g, _ := errgroup.WithContext(context.Background())

	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}

	for i := 0; i < numberTorInstances; i++ {
		//i := 1
		//for _, hsaddr := range txtlines {
		//		g.Go(func() error {
		fmt.Println("Starting proxy number " + strconv.Itoa(i) + " at " + txtlines[i])
		torProxy, err := core.CreateTorProxy(circuitInterval, txtlines[i])
		if err != nil {
			if torProxy != nil {
				torProxy.Close()
			}

			return err
		}

		ch <- *torProxy

		//		})
	}
	//return nil

	err := g.Wait()
	close(ch)

	for proxy := range ch {
		proxies = append(proxies, proxy)
	}
	defer core.CloseProxies(proxies)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	proxies1 := make([]core.TorProxy, 0)

	reverseProxy1 := &core.ReverseProxy{}
	err = reverseProxy1.Start(proxies1, port)
	if err != nil {
		return err
	}

	return nil
}

func dash(port int, circuitInterval int, hslist string) error {

	var txtlines []string
	readtodict(hslist, &txtlines)
	numberTorInstances := len(txtlines)

	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]core.TorProxy1, 0)
	ch := make(chan core.TorProxy1, numberTorInstances)

	g, _ := errgroup.WithContext(context.Background())

	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}

	for i := 0; i < numberTorInstances; i++ {
		//i := 1
		//for _, hsaddr := range txtlines {
		go func(i int) {
			fmt.Println("Starting proxy number " + strconv.Itoa(i) + " at " + txtlines[i])
			torProxy, err := core.CreateTorProxy1(circuitInterval, txtlines[i])
			if err != nil {
				if torProxy != nil {
					torProxy.Close1()
				}

			}

			ch <- *torProxy

		}(i)
	}
	//return nil

	err := g.Wait()
	close(ch)

	for proxy := range ch {
		proxies = append(proxies, proxy)
	}
	defer core.CloseProxies1(proxies)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	proxies1 := make([]core.TorProxy1, 0)

	reverseProxy1 := &core.ReverseProxy1{}

	err = reverseProxy1.Start1(proxies1, port)
	if err != nil {
		return err
	}
	return err
}

/*
func dash1(port int, circuitInterval int, hslist string) error {
	reverseProxy1 := &core.ReverseProxy1{}
	var proxies1 []int
	proxies1 = []int{9002, 9008}
	err := reverseProxy1.Start1(proxies1, port)

	return err
}
*/
