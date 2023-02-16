package main

import (
	"ChatServer/models"
	"ChatServer/pkg/handler"
	"encoding/json"
	"log"
	"os"
)

const configPath = "configs/config.local.json"

func main() {
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := serverRun(cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func serverRun(addr string) error {
	router := handler.InitRoutes()
	if err := router.Run(addr); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func loadConfig(path string) (*models.Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var config models.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, err
}
