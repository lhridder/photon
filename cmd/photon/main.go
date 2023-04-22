package main

import (
	"log"
	"photon/config"
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

	log.Println("Loading proxyconfigs...")

}
