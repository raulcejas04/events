package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
     	prod "argus-events/pkg/producer"
	"github.com/gorilla/mux"
	"net/http"
	postgres "argus-events/model/psql"
	"strings"
	"github.com/spf13/viper"
	"argus-events/pkg/parser"
	"os"	
	"encoding/json"
	"argus-events/pkg/debug"
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
               log.Infof("Send bugreport ID %s ", key )
               message:= prod.ChanIn{ Id: key, Task: "P" }
		(*p).In <- message
	}
}

func channelHandlerDel( p *prod.ProducerBR ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	       vars := mux.Vars(r)
               key := vars["id"]
               log.Infof("Reprocess bugreport ID %s ", key )
               message:= prod.ChanIn{ Id: key, Task: "D" }
		(*p).In <- message
	}
}


func main () {

	initialize()
	
	postgres.NewConnection()
	postgres.NewConnectionSql()

	var msgWorker = make(chan prod.MsgWorker )

	f:=debug.NewOpenedFile( "./debug.txt" )	
	producer:= &prod.ProducerBR{}
	producer.InitProducerDB(f)
	var knowledgeDefId uint =1
	parser:=parser.NewParser( knowledgeDefId )
	//p.InitProducerDB()
	for _,e := range (*parser).Events {
		fmt.Printf("%+v\n\n",e)
	}
		
	log.Infof("Launch web server " )	
	
    	go handleRequests( producer )

	go producer.ProducerBugRep()
	
	for i:=1;i<=1;i++ {
		go consumer( &(producer.Consumer), &msgWorker, knowledgeDefId )
	}
	
	for i:=1;i<=1;i++ {
		go worker( &msgWorker, producer.Parser, f)
	}
	<-producer.Done
	debug.FileClose(f)
	
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



func consumer( chConsumers *chan *prod.Msg, msgWorker *chan prod.MsgWorker, knowledgeDefId uint ) {
	for msg := range *chConsumers {
		//fmt.Println("Consumer file_id: ", msg.FileId)
		//fmt.Printf( "\nlifecycles %+v\n",(*((*msg).LifeCycles)).Lc )
		if (*msg).FileId==0 {
			//postgres.DbEvents.Create( *msg.ExtraEvent )
			fmt.Println( " Save Events ", msg.BugreportId, msg.PartitionId )
			fmt.Printf( "lifecycle %+v\n",(*((*msg).LifeCycles)).Lc)
			((*msg).LifeCycles).Save( msg.BugreportId, msg.PartitionId, knowledgeDefId )
		} else {
			messages:=postgres.GetContents(msg.BugreportId,msg.PartitionId,msg.FileId,msg.FileName )
			for _,m := range messages {
				//fmt.Printf("mess %+v\n",m )
				msgPar:= prod.MsgWorker{ Message: m, ExtraEvent: msg.ExtraEvent, LifeCycles: msg.LifeCycles }
				*msgWorker <- msgPar
			}
		}	
	}
}

func worker( msgWorker *chan prod.MsgWorker, parser *parser.Parser, f *os.File ) {

	for input :=range *msgWorker {
		//fmt.Printf( "worker mess %+v\n",input.Message.Mess )
		for _,e := range (*parser).Events {

			//TODO get lifecycles context
			//scenarioId:=e.ScenarioId
			//stateId:=e.StateId
			//eventId:=e.EventId
			
			keyDebug:=fmt.Sprintf( "%s-%d",input.Message.FileName,input.Message.LineNumber)
			if e.Approximate( input.Message.Mess ) {
				patternLine:=e.LlRegex

				debug.FileWrite( f,"****************************************************************************************************\n")
				debug.FileWrite( f,fmt.Sprintf( "%s Approx line: %s\n", keyDebug, input.Message.Mess ))
				debug.FileWrite( f,fmt.Sprintf( "%s Approx pattern: %s regex: %s\nreplac: %+v\nstartEnd %s\n patterLine %s\n", keyDebug,e.LogLine, e.LlRegex, e.Replacement, e.StartEnd,patternLine))
				if e.StartEnd=="E" {
					if ! input.LifeCycles.ReplValueParams(&patternLine, input.Message.BootId, e.ScenarioId, e.StateId, e.Replacement) {
						debug.FileWrite( f,fmt.Sprintf(" params NOT found %+v\n", e.Replacement))			
						continue
					}
				}


				debug.FileWrite( f, fmt.Sprintf("%s Candidate: %s\ncalculated pattern %s\n\n",keyDebug,input.Message.Mess, patternLine ))

				
				match,param:= e.ItMatchParam( &patternLine, input.Message.Mess ) 
				if match {

					extraEventIndex := postgres.ExtraEvent{
								EventID: e.EventId,
								Location: input.Message.Location,
								Pid: input.Message.Pid,
								Tid: input.Message.Tid,
								FileID: input.Message.FileId,
								FileName: input.Message.FileName,
								LineNumber: input.Message.LineNumber,
								Timestamp: input.Message.Timestamp,
								Tag: input.Message.Tag,
								Priority: input.Message.Priority,
								LogMessage: input.Message.Mess,
								ExtraParameters: []postgres.ExtraParameter{},
								}

					fmt.Printf( "\n\n**********IT MATCHED %s startend %s eventid %d\n\n", input.Message.Mess, e.StartEnd, e.EventId )
					if e.StartEnd=="E" {
					
						debug.FileWrite(f, fmt.Sprintf("%s Param: %s\n", keyDebug, param ))
						//stateIndex,stateData:=
						input.LifeCycles.AddLineToTrue(input.Message.BootId, e.ScenarioId, e.StateId, e.StartEnd, input.Message.Timestamp, &extraEventIndex )
						//fmt.Println( " index ",stateIndex," candidate state ", stateData, " scenarios ",  (*(input.LifeCycles)).Lc[input.Message.BootId] )
						//fmt.Println( " index ",stateIndex," candidate state ", stateData  )
					} else {

						var params *[]string
						if strings.Contains(e.LogLine,"%s") || strings.Contains(e.LogLine,"%d") {
							params = e.GetParameters( input.Message.Mess )
							if len(*params)>0 {
								//fmt.Printf( "IT MATCHED2 %d %d %d\nlog line %s\nparse %sparam %+v\n ", e.ScenarioId,e.StateId,e.EventId, input.Message.Mess,e.LogLine,*params )
								fmt.Printf( "param %+v\n\n",*params)
								for o,p := range *params {
									if o > 0 {
										extraEventIndex.ExtraParameters=append(extraEventIndex.ExtraParameters, postgres.ExtraParameter{ Value:p, Offset:uint(o), } )
									}
								}
							}
						}
					
						fmt.Printf("\n\n**********IT MATCHED2 scen %d  state %d event %d\n log line %s\n parse %s\n matchpara %+v words %+v\n regex %s\n params %+v\n\n",e.ScenarioId,e.StateId,e.EventId, input.Message.Mess, e.LogLine, *param, e.Words, e.LlRegex, *params )

						input.LifeCycles.AddLine( parser, &input.Message.Mess, input.Message.BootId, e.ScenarioId, e.TypeScenarioId , e.StateId, &extraEventIndex, input.Message.BootName )
				
					}
					
					for _,scen := range (*input.LifeCycles).Lc[input.Message.BootId].ResultScenarios {
						if scen.ScenarioID == e.ScenarioId {
							u, err := json.MarshalIndent(scen,"", " " )
							if err != nil {
								panic(err)
							}		
					
							fmt.Printf( "json %+v\n\n", string(u) )
						
							debug.FileWrite( f, fmt.Sprintf( "%s %s",keyDebug, string(u) ))
							debug.FileWrite( f, "\n\n" )
						}
					}

				}				
			}
		}

	}
	
}


