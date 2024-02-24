package main

import (
	postgres "argus-events/model/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"	
	"os"
	"fmt"
	"encoding/json"
)

func main() {

	initialize()

	postgres.NewConnection()
	
	postgres.Setup()
	postgres.InitLoad()
	know := postgres.KnowledgeDef{}
	know.GetKnowledgeDef(1)
	//fmt.Printf("knw %+v\n",knw )
	
        u, err := json.Marshal(know)
        if err != nil {
            panic(err)
        }
        fmt.Println(string(u))
}

func initialize() {


	// setup logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// load configuration
	viper.SetConfigName("argus-events-config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs") // config file path
	err := viper.ReadInConfig()

	if err != nil {
		log.Error("server: failed to read config file")
		log.Fatal(err)
	}

}
