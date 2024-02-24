package producer
import(
	"argus-events/model/postgres"
	"strconv"
	"fmt"
)

type ChanIn struct {
	Id string
	Task string
}
type ProducerBR struct {
    In chan ChanIn
    Consumer chan *Msg
    Done chan bool
}
type Msg struct {
	BugreportId int
	PartitionId int
	FileId  int
	FileName string
	EventIndex *[]EventIndex
	//Sql *sql.Stmt
}

func (p *ProducerBR) InitProducerDB()  {
	p.In = make(chan ChanIn, 10 )
	p.Consumer = make(chan *Msg)
	p.Done = make(chan bool)
}

func (p *ProducerBR ) ProducerBugRep( ) {
	for in := range (*p).In {
		bugreportId,_:=strconv.Atoi(in.Id)
		files,partitionId:=postgres.GetFilesId( bugreportId )
		var eventIndex []EventIndex  
		for _, inMsg := range *files {
			msg := Msg{}
			msg.FileId = inMsg.FileId
			msg.FileName = inMsg.FileName
			msg.BugreportId=bugreportId
			msg.PartitionId=partitionId
			msg.EventIndex=&eventIndex
			(*p).Consumer <- &msg
		}
		msg:=Msg{}
		msg.Id=0
		msg.BugreportId=bugreportId
		msg.PartitionId=partitionId
		msg.EventIndex=&eventIndex		
		(*p).Consumer <- &msg	
	}
	fmt.Println("Before closing channel")

	close((*p).Consumer)
	fmt.Println("Before passing true to done")
	(*p).Done <- true
}

