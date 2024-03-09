package producer
import(
	postgres "argus-events/model/psql"
	"argus-events/pkg/parser"
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"	
	"time"
	"regexp"
	"strings"
	"os"
	//"sort"
	"argus-events/pkg/debug"
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
	Message 	postgres.Message
	//Debug		*os.File
	ExtraEvent 	*[]postgres.ExtraEvent
	LifeCycles 	*Results //boot_id is the key	  	
}

type ResultEvent struct {
	EventID		uint
	TimestampUnix		int64 //unix time
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
	ResultScenarios	[]ResultScenario
}

type Results struct {
	TypeScenariosName	map[uint]string
	States			map[uint]map[uint]postgres.State
	Debug			*os.File
	Lc map[uint]Result //key is boot_id
}

type Candidate struct {
	ScenarioID 	uint
	StateID	uint
	EventID	uint
}


func (l *Results) AddLineToTrue( bootId uint, scenarioId uint, stateId uint, startEnd string, timestamp time.Time, newExtraEvent *postgres.ExtraEvent )  {

//var index int=-1


if (*l).Lc==nil {
	return
	//return -1,nil
	//(*l).Lc=make(map[uint]Result)
}

result:=(*l).Lc[bootId]

//timeUnixEvent := time.Unix()
var candidates map[int64]Candidate
candidates = make(map[int64]Candidate)
var arrCandidates []int64
Loop:
for ksc,resultScenario := range result.ResultScenarios {
	//fmt.Println( "entro scencario", resultScenario.ScenarioID, len( resultScenario.ResultStates ), startEnd )
	if resultScenario.ScenarioID==scenarioId {
		if len( resultScenario.ResultStates ) == 2 { //jump completed scenarios
			continue
		} else if startEnd == "E" {
			for _,resultState := range resultScenario.ResultStates {
				//fmt.Println( "entro stado ", resultState.StateID )
				if resultState.StartEnd=="S" {
					for _,resultEvent := range resultState.ResultEvents{
						//fmt.Println( "entro evento ", resultEvent.EventID, resultEvent.TimestampUnix, timestamp.Unix()  )
						if resultEvent.TimestampUnix < timestamp.Unix() {
							candidates[resultEvent.TimestampUnix]=Candidate{ ScenarioID: scenarioId, StateID: stateId, EventID: resultEvent.EventID }
							arrCandidates=append(arrCandidates,resultEvent.TimestampUnix)
							newResultEvent := ResultEvent{
										EventID: (*newExtraEvent).EventID,
										//StartEnd: (*l).States[scenarioId][stateId].StartEnd,
										TimestampUnix: (*newExtraEvent).Timestamp.Unix(),
										ExtraEvent: *newExtraEvent,
										}

							newResultEvent.ExtraEvent.ExtraParameters=append( newResultEvent.ExtraEvent.ExtraParameters, (*newExtraEvent).ExtraParameters... )
							
							resState := ResultState{ StateID: stateId, StartEnd: (*l).States[scenarioId][stateId].StartEnd }
							resState.ResultEvents=append( resState.ResultEvents, newResultEvent ) 
		
							debug.FileWrite( (*l).Debug, fmt.Sprintf("Add state 1 %+v\n", (*l).Lc[bootId].ResultScenarios[ksc].ResultStates ))
							(*l).Lc[bootId].ResultScenarios[ksc].ResultStates = append( (*l).Lc[bootId].ResultScenarios[ksc].ResultStates, resState )
							debug.FileWrite( (*l).Debug, fmt.Sprintf("Add state 2 %+v\n", (*l).Lc[bootId].ResultScenarios[ksc].ResultStates ))
							break Loop 						
						}
					}
				}
			}
		} else {
		
		}	
	}
}

/*found:=-1
if len(arrCandidates)>0  {
	sort.Slice(arrCandidates, func(i, j int) bool { return arrCandidates[i] < arrCandidates[j] })


	for k,timeUnixFound := range arrCandidates {
		if timeUnixFound < timestamp.Unix() {
			found=k
			break
		}
	}
	//fmt.Printf( "Candidates %+v found %d timeoftheline %d\n", arrCandidates, found,timestamp.Unix()  )	

}
if found >= 0 {
	stateFound:=candidates[arrCandidates[found]]
	return found,&stateFound
} else {
	return found,nil
}*/
return
}

func InitResults( debug *os.File ) *Results {
	Lc := make(map[uint]Result)
	return &Results{ TypeScenariosName: postgres.GetAllTypeScenario(), States: postgres.GetAllStates(), Lc: Lc, Debug: debug }
}

//add the scenario if not exist, all states and foreach state all event to process
func (l *Results) AddLine( parser *parser.Parser, line *string, bootId uint, scenarioId uint, typeScenarioId uint, stateId uint, newExtraEvent *postgres.ExtraEvent, bootName string ) {

	fmt.Println( "TypeScenariosNames ",typeScenarioId,(*l).TypeScenariosName[typeScenarioId] )
	var add=true
	if (*l).TypeScenariosName[typeScenarioId]!= "dynamic_state" {
		add=false
		// if it is exists in the Result
		if _,ok:=(*l).Lc[bootId]; !ok {
			(*l).Lc[bootId]=Result{}
		}
		scenarios:=(*l).Lc[bootId].ResultScenarios	
	

		//var startEnd string
		//was the line processed ?
		for _,sce := range scenarios {
			if sce.ScenarioID==scenarioId {
				if len(sce.ResultStates)==0 {
					//startEnd="T"
					add=true
				} else {
					for _,st := range sce.ResultStates {
						if st.StateID==stateId {
							
							for _,eve := range st.ResultEvents {
								if (*newExtraEvent).ID == eve.EventID {
									add=false
									break
								}
							}
							if !add {
								break
							}
						}
					}
					if !add {
						break
					}
				}
			}
			if !add {
				break
			}
		}
	}
	
	if add {
	
		//TODO check if it is possible to add or not



		//(*l).Lc[bootId].ResultScenarios=append((*l).Lc[bootId].ResultScenarios, ResultScenario{} )
				
		newResultEvent := ResultEvent{
					EventID: (*newExtraEvent).EventID,
					//StartEnd: (*l).States[scenarioId][stateId].StartEnd,
					TimestampUnix: (*newExtraEvent).Timestamp.Unix(),
					ExtraEvent: *newExtraEvent,
					}
							
		resState := ResultState{ StateID: stateId, StartEnd: (*l).States[scenarioId][stateId].StartEnd }
		resState.ResultEvents=append( resState.ResultEvents, newResultEvent ) 
		
		//resScenario:= (*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]
		var resScenario ResultScenario
		resScenario.ScenarioID = scenarioId
		resScenario.Timestamp=(*newExtraEvent).Timestamp
		resScenario.TypeScenarioName=(*l).TypeScenariosName[scenarioId]
		resScenario.ResultStates=append( resScenario.ResultStates, resState  )
		
		//fmt.Printf( "\nresScenario %+v\n", resScenario )

		//(*l).Lc[bootId].ResultScenarios[(*newExtraEvent).Timestamp.Unix()]=resScenario
		//(*l).Lc[bootId].ResultScenarios=append((*l).Lc[bootId].ResultScenarios, resScenario )
		
		result,ok:=(*l).Lc[bootId]; 
		if !ok {
			result = Result{ BootName: bootName }
		}
		result.ResultScenarios=append( result.ResultScenarios, resScenario )
		(*l).Lc[bootId]=result

	}
	//if belong to a true state and it doesn't exist any event create a new scenario
	
	//if belong ot a true state and it is a any_found 
	
	
	return	
}

//in the matching
func ( l *Results ) ReplValueParams( patterLine *string, bootId uint, scenarioId uint, stateId uint, replacement []string ) bool {

result:=(*l).Lc[bootId]

var param []postgres.ExtraParameter
for _,repla :=range replacement {
	found:=false
	re := regexp.MustCompile("%(\\d+)%(\\d+)")
	match := re.FindStringSubmatch(repla)
	indexEvent,_:=strconv.Atoi(match[1])
	indexParam,_:=strconv.Atoi(match[2])
	fmt.Printf( "++++++++rep %+v match %d %d\n\n",repla, indexEvent, indexParam )
	Out:
	for _,resultScenario := range result.ResultScenarios {
		//fmt.Println( "entro scencario", resultScenario.ScenarioID, len( resultScenario.ResultStates ), startEnd )
		if resultScenario.ScenarioID==scenarioId {
			if len( resultScenario.ResultStates ) == 2 { //jump completed scenarios
				continue
			} else  {
				for _,resultState := range resultScenario.ResultStates {
					//fmt.Println( "entro stado ", resultState.StateID )
					if resultState.StartEnd=="S" {
						i:=0
						for _,resultEvent := range resultState.ResultEvents{
							if i==indexEvent {
								//fmt.Println( "resultevent ", resultEvent)
								param=resultEvent.ExtraEvent.ExtraParameters
								*patterLine=strings.Replace( *patterLine, repla, param[indexParam].Value, 1 )
								found=true
								break Out
							}
							i++
						}
					}
				}
			}
		}
	}
	if !found {
		return false
	}
}
return true
}

func ( l *Results ) Save( bugreportId int, partitionId int, knowledgeDefId uint ) {

var extraKnow *postgres.ExtraKnow
for bootId,lc :=range (*l).Lc {
	extraKnow = postgres.NewExtraKnow( bugreportId, partitionId, knowledgeDefId, bootId, lc.BootName )
	fmt.Println( "bootId %d %s\n",bootId, lc.BootName )

	for _,sce := range lc.ResultScenarios {
		extraScenario:=postgres.ExtraScenario{ ScenarioID: sce.ScenarioID }
		
		for _,sta := range sce.ResultStates {
			extraState := postgres.ExtraState{ StateID: sta.StateID }

			for _,eve := range sta.ResultEvents {		
				extraState.AddExtraEvent( eve.ExtraEvent )
			}
			extraScenario.AddExtraState( extraState )
		}
		extraKnow.AddExtraScenario( extraScenario )
	}
	postgres.DbEvents.Create( extraKnow )
}
}

func (p *ProducerBR) InitProducerDB( debug *os.File)  {
	p.In = make(chan ChanIn, 10 )
	p.Consumer = make(chan *Msg)
	p.Done = make(chan bool)
	p.Parser = parser.NewParser( 1 )
	p.LifeCycles = InitResults( debug )
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

