package main
import(
	"fmt"
	postgres "argus-events/model/psql"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"	
    	"encoding/json"	
)

func main() {
	initialize()
	postgres.NewConnection()
	l:=postgres.GetLearning( 2711, 1, "boot_1" )

        u, err := json.MarshalIndent(l,"", "\t")
        if err != nil {
            panic(err)
        }

	if l!=nil {
		fmt.Printf("%s\n",string(u))
	}
}

func initialize() {
	viper.SetConfigName("argus-events-config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs") 
	err := viper.ReadInConfig()

	if err != nil {
		log.Error("server: failed to read config file")
		log.Fatal(err)
	}
}
