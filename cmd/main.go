package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
     	prod "argus-events/pkg/producer"
	"github.com/gorilla/mux"
	"net/http"
	"argus-events/model/postgres"
	"argus-events/pkg/parse"
	"strings"
	"github.com/spf13/viper"	
)



func handleRequests(  p *prod.ProducerBR ) {
    r := mux.NewRouter()
    
    chPost := channelHandlerPost(p)
    chDel := channelHandlerDel(p)    
    //chGet := channelHandlerGet(p)
        
    r.HandleFunc("/{id}", chPost ).Methods("POST")
    r.HandleFunc("/{id}", chDel ).Methods("DELETE")
    //r.HandleFunc("/{id}", chGet ).Methods("GET")
    
    srv := &http.Server{
	  Addr:    ":7042",
	  Handler: r,
    }
    srv.ListenAndServe()

}


func channelHandlerPost( p *prod.ProducerBR ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	       vars := mux.Vars(r)
               key := vars["id"]
               log.Infof("Save bugreport ID %s ", key )
               message:= prod.ChanIn{ Id: key, Task: "P" }
		(*p).In <- message
	}
}

func channelHandlerDel( p *prod.ProducerBR ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	       vars := mux.Vars(r)
               key := vars["id"]
               log.Infof("Delete bugreport ID %s ", key )
               message:= prod.ChanIn{ Id: key, Task: "D" }
		(*p).In <- message
	}
}


func main () {

	initialize()
	
	postgres.NewConnectionSql()
	var msgParser = make(chan string)
	
	p:=&prod.ProducerBR{}
	p.InitProducerDB()
	
	log.Infof("Launch web server " )	
	
    	go handleRequests( p )

	go (*p).ProducerBugRep()
	
	for i:=1;i<4;i++ {
		go consume( &((*p).Consumer), &msgParser )
	}
	for i:=1;i<4;i++ {
		go worker( &msgParser )
	}
	<-(*p).Done
	
}

func initialize() {
	viper.SetConfigName("argus-events-config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../configs") 
	err := viper.ReadInConfig()

	if err != nil {
		log.Error("server: failed to read config file")
		log.Fatal(err)
	}
}

func consume( chConsumers *chan *prod.Msg, msgParser *chan string ) {
	for msg := range *chConsumers {
		fmt.Println("Consumer file_id: ", msg.Id)
		messages:=postgres.GetContents(msg.BugreportId,msg.PartitionId,msg.Id )
		for _,m := range messages {
				//fmt.Println( "tag ",m.Tag)
				if strings.Contains(m.Tag,"ActivityManager" ) {
					//fmt.Println( "tag ",m.tag,m.mess )
					*msgParser <- m.Mess
				}
		}

	}
}

func worker( msgParser *chan string ) {
	for input :=range *msgParser {
		e:=parse.Event{ LogLine: "Start proc %d:%s for activity {%s/%s}" }
		//split
		e.GetWords()
		/*if len(input)>=10 && input[0:10]=="Start proc" {
			fmt.Println( "input ",input)
		}*/
		//fmt.Println( "input ",input)
		if e.Approximate( input ) {
			fmt.Println( "IT MATCHED ", input )
			e.GetParameters( input )
		}
	}
}

