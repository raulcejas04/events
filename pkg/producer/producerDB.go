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
    LifeCycles *map[uint]Results //boot_id is the key
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
	postgres.ExtraEvent
}


type ResultState struct {
	ID		uint
	ResultEvents	[]ResultEvent
}

type ResultScenario struct {
	ID		uint
	Timestamp 	time.Time
	TypeScenario	string //fatal_error, dynamic_state, etc.
	Start		ResultState
	End		ResultState
}

type Result struct {
	BootName 	string
	ResultScenarios	map[int64]ResultScenario //key is unix time (int64)
}

type Results map[uint]Result //key is boot_id



func (r Results) GetEventsToProcess() {


}

//add the scenario if not exist, all states and foreach state all event to process
func (l *Results) AddLine( parser *parser.Parser, line string, bootId int, scenarioId uint,  stateId uint, eventId uint) {
	scenarios:=(*Results)[bootId].ResultScenarios
	//line belonf to a true or false state
	for _,sce := range scenarios {
		if sce.ID==scenarioID {
			eventsOfState:=parser.GetEvents( scenarioId, stateId )
				
		}
	}
	//if belong to a true state and it doesn't exist any event create a new scenario
	
	//if belong ot a true state and it is a any_found 
	
	
	
}

func (p *ProducerBR) InitProducerDB()  {
	p.In = make(chan ChanIn, 10 )
	p.Consumer = make(chan *Msg)
	p.Done = make(chan bool)
	p.Parser = parser.NewParser( 1 )
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
			(*p).Consumer <- &msg
		}
		msg:=Msg{}
		msg.FileId=0
		msg.BugreportId=bugreportId
		msg.PartitionId=partitionId
		msg.ExtraEvent=&extraEvent		
		(*p).Consumer <- &msg	
	}
	fmt.Println("Before closing channel")

	close((*p).Consumer)
	fmt.Println("Before passing true to done")
	(*p).Done <- true
}

