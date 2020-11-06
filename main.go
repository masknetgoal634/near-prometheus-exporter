package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	nearapi "github.com/bisontrails/near-exporter/client"
	"github.com/bisontrails/near-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var version = "undefined"

	flag.Usage = func() {
		const (
			usage = "Usage: near_exporter [option] [arg]\n\n" +
				"Prometheus exporter for Near node metrics\n\n" +
				"Options and arguments:\n"
		)

		fmt.Fprint(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()

		os.Exit(2)
	}

	url := flag.String("url", "http://localhost:3030", "Near JSON-RPC URL")
	externalRpc := flag.String("external-rpc", "https://rpc.betanet.near.org", "Near JSON-RPC URL")
	addr := flag.String("addr", ":9333", "listen address")
	accountId := flag.String("accountId", "test", "Validator account id")
	ver := flag.Bool("v", false, "print version number and exit")

	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
	}

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	client := nearapi.NewClient(*url)

	devClient := nearapi.NewClient(*externalRpc)

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		collector.NewNodeRpcMetrics(client, devClient, *accountId),
		collector.NewDevNodeRpcMetrics(devClient),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
		ErrorHandling: promhttp.ContinueOnError,
	})

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
