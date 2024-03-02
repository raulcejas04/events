package producer
import(
	"argus-events/model/postgres"
	"argus-events/pkg/parser"
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"	
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
    LifeCycles *map[uint]LifeCycle //boot_id is the key
}
type Msg struct {
	BugreportId int
	PartitionId int
	FileId  uint
	FileName string
	ExtraEvent *[]postgres.ExtraEvent
	LifeCycles *map[uint]LifeCycle //boot_id is the key
}

type MsgWorker struct {
	Message postgres.Message
	ExtraEvent *[]postgres.ExtraEvent
	LifeCycles *map[uint]LifeCycle //boot_id is the key	  	
}

type LifeCycle struct {
	Scenarios		map[uint]ScenarioProcessed
	BootName		string
}

type ScenarioProcessed struct {
	States			map[uint]StateProcessed
	Result			bool //if all states where processed
}

type StateProcessed struct {
	Events			map[uint]EventProcessed
	Result			bool //successful or not for example all_found requires all events found	
}

type EventProcessed struct {
	Mess	string		//pattern
	Line 	string		//log line
	Result	bool		//found or not

}

func (l LifeCycle) GetEventsToProcess() {


}

//add the scenario if not exist, all states and foreach state all event to process
func (l LifeCycle) AddState( line string, stateId uint, scenarioId uint, eventId uint) {


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

