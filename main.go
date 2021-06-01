package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/mitchr/gossip/server"
)

var (
	sPass      bool
	oPass      bool
	configPath = "./config.json"
)

func init() {
	flag.BoolVar(&sPass, "s", false, "sets server password")
	flag.BoolVar(&oPass, "o", false, "add a server operator (username and pass)")
	flag.Parse()
}

func main() {
	c, err := server.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if sPass {
		err := server.SetPass(c)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if oPass {
		err := server.AddOp(c)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	s, err := server.New(c)
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()

	// capture OS interrupt signal so that we can gracefully shutdown server
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		s.Close()
	}()

	s.Serve()
}
