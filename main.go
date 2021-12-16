package main

import (
	"bufio"
	"sync"

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

	var txtlines []string
	readtodict(hslist, &txtlines)
	numberTorInstances := len(txtlines)

	log.Println(fmt.Sprintf("Starting tor proxies"))

	proxies := make([]core.TorProxy1, numberTorInstances)
	//ch := make(chan core.TorProxy1, numberTorInstances)
	//var proxyports []int
	//proxies1 = []int{9002, 9008}

	//g, _ := errgroup.WithContext(context.Background())

	for _, eachline := range txtlines {
		fmt.Println(eachline)
	}
	//wg := new(sync.WaitGroup)
	//wg.Add(numberTorInstances - 1)
	var err error
	for i := 0; i < numberTorInstances; i++ {
		proxies[i], err = core.CreateTorProxy1(circuitInterval, txtlines[i])
	}
	/*
		for i := 0; i < numberTorInstances-1; i++ {
			//i := 1
			//for _, hsaddr := range txtlines {
			//go func(i int) {
			//g.Go(func() error {
			fmt.Println("Starting proxy number " + strconv.Itoa(i) + " at " + txtlines[i])
			torProxy, err := core.CreateTorProxy1(circuitInterval, txtlines[i])
			if err != nil {
				if torProxy != nil {
					torProxy.Close1()
				}
				fmt.Println("ERRORRR")
				//return
			}
			//port1 := torProxy.DoubleProxyPort
			//proxyports = append(proxyports, *port1)
			//torProxy.socks5s.ListenAndServe("tcp", socks5Address)
			proxies[i] = *torProxy
			//proxies = append(proxies, *torProxy)

			//ch <- *torProxy
			//fmt.Println("ERRORRR")
			//	return
			//	wg.Done()

			//}(i)
			//})
		}
	*/
	//wg.Wait()
	//err := g.Wait()

	//wg1 := new(sync.WaitGroup)
	//wg1.Add(numberTorInstances)
	/*	close(ch)
		for proxy := range ch {
			proxies = append(proxies, proxy)
			//wg1.Done()
		}*/

	//wg1.Wait()
	fmt.Println(err)
	wg2 := new(sync.WaitGroup)
	wg2.Add(numberTorInstances - 1)
	for i := 0; i < len(proxies)-1; i++ {
		go func() {
			socks5Address := "127.0.0.1:" + strconv.Itoa(*proxies[i].DoubleProxyPort)
			err := proxies[i].Socks5s.ListenAndServe("tcp", socks5Address)
			if err != nil {
				fmt.Println("failed to create socks5 server", err)
				//return err
				return
			}
			fmt.Println("kreated u a socks serva at "+socks5Address, err)

			//wg2.Done()
			return
		}()
	}
	wg2.Wait()

	//return nil
	/*	if err := torProxy1.socks5s.ListenAndServe("tcp", socks5Address); err != nil {
			fmt.Println("failed to create socks5 server", err)
		}
		fmt.Println("kreated u a socks serva at "+socks5Address, err)

		if err != nil {
			return nil, err
		}*/

	//	defer core.CloseProxies1(proxies)

	//if err != nil {
	//	return err
	//}

	log.Println(fmt.Sprintf("Started %d tor proxies", len(proxies)))
	log.Println(fmt.Sprintf("Start reverse proxy on port %d", port))

	//reverseProxy1 := &core.ReverseProxy1{}
	//
	//err := reverseProxy1.Start2(proxies, port)
	//if err != nil {
	//	return err
	//}

	return nil
}
