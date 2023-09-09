package main

import (
	"github.com/lhridder/photon"
	"github.com/lhridder/photon/config"
	"log"
)

func main() {
	log.Println("Starting Photon: A Minecraft Bedrock reverse proxy")

	cfg, err := config.LoadGlobal()
	if err != nil {
		log.Printf("Failed to load global config: %s", err)
		return
	}

	if cfg.Debug {
		log.Println("Running in debug mode")
	}

	//TODO log.Println("Loading proxyconfigs...")
	gw := photon.Gateway{
		ListenTo: ":25565",
		Proxies:  nil,
		Cfg:      *cfg,
	}

	err = gw.Listen()
	if err != nil {
		log.Printf("failed to open listener: %s", err)
		return
	}

	gw.Serve()

}
