package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration structure for prometheus_ps
type Configuration struct {
	WatchList []string
	Wl        map[string]Void
	Port      int
}

// Void : Empty struct
type Void struct{}

// ReadConfig read the configuration file at loc
func ReadConfig(loc string) Configuration {
	fmt.Println("Reading conf.json ...")
	file, _ := os.Open(loc)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error reading config", err)
		os.Exit(1)
	}
	fmt.Printf("WatchList content : %s\n", config.WatchList)
	config.Wl = make(map[string]Void)
	for _, x := range config.WatchList {
		config.Wl[x] = Void{}
	}
	return config
}
