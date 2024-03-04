package producer
import(
	"argus-events/model/postgres"
	"argus-events/pkg/parser"
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"	
	"time"
)

type ChanIn struct {
	Id string
	Task string
}
type ProducerBR struct {
    In chan ChanIn
    Consumer chan *Msg
    Done chan bool
    Parser *parser.Parser
    LifeCycles *Results //boot_id is the key
}
type Msg struct {
	BugreportId int
	PartitionId int
	FileId  uint
	FileName string
	ExtraEvent *[]postgres.ExtraEvent
	LifeCycles *Results //boot_id is the key
}

type MsgWorker struct {
	Message postgres.Message
	ExtraEvent *[]postgres.ExtraEvent
	LifeCycles *Results //boot_id is the key	  	
}

type ResultEvent struct {
	EventID			uint
	ExtraEvent		postgres.ExtraEvent
}


type ResultState struct {
	StateID			uint
	StartEnd		string
	ResultEvents		[]ResultEvent
}

type ResultScenario struct {
	ScenarioID		uint
	Timestamp 		time.Time
	TypeScenarioName	string //fatal_error, dynamic_state, etc.
	ResultStates		[]ResultState
}

type Result struct {
	BootName 		string
	ResultScenarios		map[int64]ResultScenario //key is unix time (int64)
}

type Results struct {
	TypeScenariosName	map[uint]string
	States			map[uint]map[uint]postgres.State
	Lc map[uint]Result //key is boot_id
}


func (r Results) GetEventsToProcess() {


}

func InitResults() *Results {
	Lc := make(map[uint]Result)
	return &Results{ TypeScenariosName: postgres.GetAllTypeScenario(), States: postgres.GetAllStates(), Lc: Lc }
}

//add the scenario if not exist, all states and foreach state all event to process
func (l *Results) AddLine( parser *parser.Parser, line *string, bootId uint, scenarioId uint,  stateId uint, newExtraEvent *postgres.ExtraEvent ) {
	if _,ok:=(*l).Lc[bootId]; !ok {
		(*l).Lc[bootId]=Result{}
	}
	scenarios:=(*l).Lc[bootId].ResultScenarios	
	
	var add=false
	//var startEnd string
	//was the line processed ?
	for _,sce := range scenarios {
		if sce.ScenarioID==scenarioId {
			if sce.TypeScenarioName == "dynamic_state" {
				//it should be added
				add=true
			} else {
				add=true //if it is found add should be false
				if len(sce.ResultStates)==0 {
					//startEnd="T"
					add=true
				} else {
					for _,st := range sce.ResultStates {
						if st.StateID==stateId {
							
							for _,eve := range st.ResultEvents {
								if (*newExtraEvent).ID == eve.EventID {
									add=false
									return
								}
							}
						}
					}
				}
			}
		}	
	}
	
	if add {

		if _,ok:=(*l).Lc[bootId]; !ok {
			(*l).Lc[bootId]=Result{}
		}
		if _,ok:=(*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]; !ok {
			(*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]=ResultScenario{}
		}
		
		newResultEvent := ResultEvent{
					EventID: (*newExtraEvent).ID,
					//StartEnd: (*l).States[scenarioId][stateId].StartEnd,
					ExtraEvent: *newExtraEvent,
					}
							
		resState := ResultState{ StartEnd: (*l).States[scenarioId][stateId].StartEnd }
		resState.ResultEvents=append( resState.ResultEvents, newResultEvent ) 
		
		resScenario:= (*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]
		resScenario.Timestamp=(*newExtraEvent).Timestamp
		resScenario.TypeScenarioName=(*l).TypeScenariosName[scenarioId]
		resScenario.ResultStates=append( resScenario.ResultStates, resState  )

		(*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]=resScenario
	}
	//if belong to a true state and it doesn't exist any event create a new scenario
	
	//if belong ot a true state and it is a any_found 
	
	
	return	
}

func (p *ProducerBR) InitProducerDB()  {
	p.In = make(chan ChanIn, 10 )
	p.Consumer = make(chan *Msg)
	p.Done = make(chan bool)
	p.Parser = parser.NewParser( 1 )
	p.LifeCycles = InitResults()
	//fmt.Printf(" lifecycles %+v\n\n",p.LifeCycles )
}

func (p *ProducerBR ) ProducerBugRep( ) {
	for in := range (*p).In {
		bugreportId,_:=strconv.Atoi(in.Id)
		files,partitionId:=postgres.GetFilesId( bugreportId )
		var extraEvent []postgres.ExtraEvent 
		if len(*files)==0 {
			log.Errorf("Bugreport %d does not contain files",bugreportId)
			return
		}
		for _, inMsg := range *files {
			msg := Msg{}
			msg.FileId = inMsg.FileId
			msg.FileName = inMsg.FileName
			msg.BugreportId=bugreportId
			msg.PartitionId=partitionId
			msg.ExtraEvent=&extraEvent
			msg.LifeCycles=p.LifeCycles
			(*p).Consumer <- &msg
		}
		msg:=Msg{}
		msg.FileId=0
		msg.BugreportId=bugreportId
		msg.PartitionId=partitionId
		msg.ExtraEvent=&extraEvent
		msg.LifeCycles=p.LifeCycles				
		(*p).Consumer <- &msg	
	}
	fmt.Println("Before closing channel")

	close((*p).Consumer)
	fmt.Println("Before passing true to done")
	(*p).Done <- true
}

