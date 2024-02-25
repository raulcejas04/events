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
	
	postgres.NewConnection()
	know := postgres.KnowledgeDef{}
	know.GetKnowledgeDef(1)
	events:=know.GetEvents()
	postgres.NewConnectionSql()
	var msgParser = make(chan prod.MsgWorker )
	
	p:=&prod.ProducerBR{}
	p.InitProducerDB()
	
	log.Infof("Launch web server " )	
	
    	go handleRequests( p )

	go (*p).ProducerBugRep()
	
	for i:=1;i<4;i++ {
		go consumer( &((*p).Consumer), &msgParser )
	}
	for i:=1;i<10;i++ {
		go worker( &msgParser, events )
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



func consumer( chConsumers *chan *prod.Msg, msgParser *chan prod.MsgWorker ) {
	for msg := range *chConsumers {
		fmt.Println("Consumer file_id: ", msg.FileId)
		if msg.FileId==0 {
			postgres.DbEvents.Create( *msg.EventIndex )
			fmt.Println( " Save Events ", msg.BugreportId, msg.PartitionId )
			fmt.Printf( "%+v\n",*(msg.EventIndex))		
		} else {
			messages:=postgres.GetContents(msg.BugreportId,msg.PartitionId,msg.FileId,msg.FileName )
			for _,m := range messages {
				//fmt.Println( "tag ",m.Tag)
				if strings.Contains(m.Tag,"ActivityManager" ) {
					//fmt.Println( "tag ",m.tag,m.mess )
					msgPar:= prod.MsgWorker{ Message: m, EventIndex: msg.EventIndex }
					*msgParser <- msgPar
				}
			}
		}	
	}
}

func worker( msgParser *chan prod.MsgWorker, events *map[uint]map[uint]map[uint]string ) {
	for input :=range *msgParser {
		for scenarioId,scen := range *events {
			for stateId,state := range scen {
				for eventId,even := range state {
					e:=parse.Event{ LogLine: even }
					//split
					e.GetWords()
					//fmt.Println( "input ",input)
					if e.Approximate( input.Message.Mess ) {
						fmt.Println( "IT MATCHED ", input.Message.Mess )
						eventIndex := postgres.EventIndex{
													BugreportID: input.Message.BugreportId,
													PartitionID: input.Message.PartitionId,
													EventID: eventId,
													Location: input.Message.Location,
													BootID: input.Message.BootId,
													BootName: input.Message.BootName,
													FileID: input.Message.FileId,
													FileName: input.Message.FileName,
													LineNumber: input.Message.LineNumber,
													Timestamp: input.Message.Timestamp,
													Message: input.Message.Mess,
													Parameters: []postgres.Parameter{},
													}
																			
						if strings.Contains(even,"%s") || strings.Contains(even,"%d") {
							params := e.GetParameters( input.Message.Mess )
							if len(*params)>0 {
								fmt.Println( "IT MATCHED2 ", scenarioId,stateId,eventId, input.Message.Mess," param ",*params )
								for o,p := range *params {
									//fmt.Println( " value ",p," offset ",o )
									eventIndex.Parameters=append(eventIndex.Parameters, postgres.Parameter{ Value:p, Offset:uint(o), } )
								}
							}
						}

						//fmt.Println( "Event index param ", eventIndex.Parameters )
						*(input.EventIndex) = append( *(input.EventIndex),  eventIndex )
						//fmt.Printf( "\n\n eventindex %+v\n ", (*input.EventIndex)[0].Parameters )
						length:=len( *(input.EventIndex) )
						if length > 0 {
							//fmt.Printf( "\n\n eventindex %+v\n ", *(input.EventIndex) )
							for _,p := range  eventIndex.Parameters {
								(*input.EventIndex)[length-1].Parameters = append((*input.EventIndex)[length-1].Parameters,  p )
							}
						}
					}
				}

			}
		}			
		
	}
}

