package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
     	prod "argus-events/pkg/producer"
	"github.com/gorilla/mux"
	"net/http"
	"argus-events/model/postgres"
	"strings"
	"github.com/spf13/viper"
	"argus-events/pkg/parser"		
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
	postgres.NewConnectionSql()

	var msgParser = make(chan prod.MsgWorker )
	
	producer:= &prod.ProducerBR{}
	producer.InitProducerDB()
	parser:=parser.NewParser( 1 )
	//p.InitProducerDB()
	for _,e := range (*parser).Events {
		fmt.Printf("%+v\n\n",e)
	}
		
	log.Infof("Launch web server " )	
	
    	go handleRequests( producer )

	go producer.ProducerBugRep()
	
	for i:=1;i<4;i++ {
		go consumer( &(producer.Consumer), &msgParser )
	}
	for i:=1;i<10;i++ {
		go worker( &msgParser, producer.Parser )
	}
	<-producer.Done
	
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



func consumer( chConsumers *chan *prod.Msg, msgWorker *chan prod.MsgWorker ) {
	for msg := range *chConsumers {
		//fmt.Println("Consumer file_id: ", msg.FileId)
		//fmt.Printf( "\nlifecycles %+v\n",(*((*msg).LifeCycles)).Lc )
		if msg.FileId==0 {
			//postgres.DbEvents.Create( *msg.ExtraEvent )
			fmt.Println( " Save Events ", msg.BugreportId, msg.PartitionId )
			fmt.Printf( "%+v\n",(*((*msg).LifeCycles)).Lc)		
		} else {
			messages:=postgres.GetContents(msg.BugreportId,msg.PartitionId,msg.FileId,msg.FileName )
			for _,m := range messages {
				//fmt.Println( "tag ",m.Tag)
				//if strings.Contains(m.Tag,"ActivityManager" ) {
					
					msgPar:= prod.MsgWorker{ Message: m, ExtraEvent: msg.ExtraEvent, LifeCycles: msg.LifeCycles }
					*msgWorker <- msgPar
				//}
			}
		}	
	}
}

func worker( msgWorker *chan prod.MsgWorker, parser *parser.Parser ) {

	for input :=range *msgWorker {
		//fmt.Printf( "worker mess %+v\n",input.Message.Mess )
		for _,e := range (*parser).Events {
		//fmt.Printf("\nevent %+v\n",e)

		//TODO get lifecycles context
		/*lc:=input.LifeCycles[input.Message.BootId]
		evp:=lc.GetEventsToProcess()*/
		//fmt.Println( " Message ",input.Message.Mess )


			//scenarioId:=e.ScenarioId
			//stateId:=e.StateId
			//eventId:=e.EventId
			
					//e:=parse.Event{ LogLine: even }

			if e.Approximate( input.Message.Mess ) {
				//fmt.Printf( "\n\n**********IT MATCHED %s\n\n", input.Message.Mess )
				e.ReplValueParams()
				match,param:= e.ItMatchParam( input.Message.Mess ) 
				if match {
					fmt.Printf("\n\n**********IT MATCHED2 scen %d  state %d event %d\n log line %s\n parse %s\n matchpara %+v\n\n",e.ScenarioId,e.StateId,e.EventId, input.Message.Mess, e.LogLine, *param )

					extraEventIndex := postgres.ExtraEvent{
								EventID: e.EventId,
								Location: input.Message.Location,
								Pid: input.Message.Pid,
								Tid: input.Message.Tid,
								FileID: input.Message.FileId,
								FileName: input.Message.FileName,
								LineNumber: input.Message.LineNumber,
								Timestamp: input.Message.Timestamp,
								Message: input.Message.Mess,
								ExtraParameters: []postgres.ExtraParameter{},
								}
					if strings.Contains(e.LogLine,"%s") || strings.Contains(e.LogLine,"%d") {
						params := e.GetParameters( input.Message.Mess )
						if len(*params)>0 {
							//fmt.Printf( "IT MATCHED2 %d %d %d\nlog line %s\nparse %sparam %+v\n ", e.ScenarioId,e.StateId,e.EventId, input.Message.Mess,e.LogLine,*params )
							fmt.Printf( "param %+v\n\n",*params)
							for o,p := range *params {
								extraEventIndex.ExtraParameters=append(extraEventIndex.ExtraParameters, postgres.ExtraParameter{ Value:p, Offset:uint(o), } )
							}
						}
					}
					input.LifeCycles.AddLine( parser, &input.Message.Mess, input.Message.BootId, e.ScenarioId, e.TypeScenarioId , e.StateId, &extraEventIndex )
				
				}			
				
				
				/*


				*(input.EventIndex) = append( *(input.EventIndex),  eventIndex )

				length:=len( *(input.EventIndex) )
				if length > 0 {

					for _,p := range  eventIndex.Parameters {
						(*input.EventIndex)[length-1].Parameters = append((*input.EventIndex)[length-1].Parameters,  p )
					}
				}*/
			}
		}

	}
	
}

