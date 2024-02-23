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
	Id  int
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
		for _, file_id := range *files {
			msg := Msg{}
			msg.Id = file_id
			msg.BugreportId=bugreportId
			msg.PartitionId=partitionId
			(*p).Consumer <- &msg
		}
	}
	fmt.Println("Before closing channel")

	close((*p).Consumer)
	fmt.Println("Before passing true to done")
	(*p).Done <- true
}

