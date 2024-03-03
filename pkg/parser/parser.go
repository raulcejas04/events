package parser

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"argus-events/model/postgres"    
)

type UsedParam struct {
	RefIndEvent int
	RefEvent int
	RefIndParam int
	RefParam int
}

type Event struct{
	ScenarioId uint
	StateId uint
	EventId uint
	LogLine string //Log line original without parameters
	LlRegex string //with regex
	Replacement []string
	Words []string
	TotalLengthWords int
	UsedParams []UsedParam
}

type Parser struct {
	KnowledgeDef postgres.KnowledgeDef
	Events []Event
}


func UsedParams(line string) []UsedParam {
	var res []UsedParam
	var begin int
	begin = 0
	//i:=0
	for {
		var index int = 0
		//fmt.Println( "begin ",begin," line ",line[begin:], len(line[begin:]) )

		if index = strings.Index(line[begin:], "%"); index != -1 {
			index=begin+index
			//fmt.Println( "1 ",  index, line[index+1:index+2], line[index+3:index+4] )
			if line[index+2:index+3] == "%" && !( line[index+1:index+2] == "s" || line[index+3:index+4] == "d" ) {
				//fmt.Println( "2 ", line[index+1:index+2], line[index+3:index+4] )
				refEvent,_ := strconv.Atoi(line[index+1 : index+2])
				refParam,_ := strconv.Atoi(line[index+3 : index+4])
				res = append(res, UsedParam{ RefIndEvent: index+1, RefEvent: refEvent, RefIndParam: index+3, RefParam: refParam })
			}
		}
		//fmt.Println( "index ",index," len ", len(line) )
		//fmt.Printf("\n")
		if index < len(line) && index!=-1 {
			if line[index+1:index+2]=="s" || line[index+1:index+2]=="d" {
				begin =index+2 
			} else {
				begin = index + 4
			}
		} else {
			break
		}
	}
	return res
}

func NewParser( knowledgeDef uint ) *Parser {
	parser := Parser{}
	know := postgres.KnowledgeDef{}
	know.GetFullKnowledgeDef(knowledgeDef)
	fmt.Printf("know %+v\n",know)
	dbEvents:=*(know.GetEvents())
	var events []Event
	for scenarioId,scen := range dbEvents {
		for stateId,state := range scen {
			for eventId,e := range state {	
				event:= Event{ScenarioId: scenarioId, StateId: stateId, EventId: eventId, LogLine:e}
				event.GetWords()
				event.RegularExpression()
				event.InitValueParams()
				events=append(events, event)
			}
		}
	}
	parser.KnowledgeDef=know
	parser.Events=append(parser.Events, events... )
	return &parser
}

func (parse *Parse ) GetEvents( scenarioId int, stateId int ) {


}

func (event *Event ) Approximate( line string ) bool {

	if event.TotalLengthWords > len(line) {
		//fmt.Println( "Line ",line," totallengthword ", event.TotalLengthWords, len(line) )			
		return false
	}
 
	x:=0
	res:=true
	for _,word := range event.Words {
		y := strings.Index(line[x:], word)
		//fmt.Printf("----x line %s\n word %s\n x %+v\n words %+v\n\n", line[x:], word, x, event.Words)
		if y > -1 {
			x=x+y+len(word)
		} else {
			res=false
			break
		}
	}

	return res
}

func ( event *Event ) ItMatchParam( input string ) (bool,*[]string) {
	pattern:= regexp.MustCompile( event.LlRegex )
	if pattern.MatchString(input) {
		if strings.Contains(event.LogLine, "%s") || strings.Contains(event.LogLine, "%d") {

			parameters := pattern.FindStringSubmatch( input )
			if len(parameters)>0 {
				return true,&parameters
			}
		}
		return true,nil
	} else {
		return false,nil
	}
	return false,nil
}

//in the constructor
func (event *Event ) InitValueParams( )  {
	re := regexp.MustCompile( "%\\d+%\\d+" )
	event.Replacement=re.FindStringSubmatch( event.LogLine ) 
	return
}

//in the matching
func (event *Event ) ReplValueParams( )  {
	for _,rep := range event.Replacement {
		fmt.Printf( "rep %+v\n",rep )
	}
	return
}


func ( event *Event ) RegularExpression() {
	event.LlRegex=strings.Replace( event.LogLine, "%s", "(.+)", -1)
	event.LlRegex=strings.Replace( event.LlRegex, "%d", "(\\d+)", -1)
	return
}

func ( event *Event ) GetWords( ) {

	re := regexp.MustCompile("%(d|s|\\d)")
	newStr := re.ReplaceAllString(event.LogLine, " ")
	event.Words=strings.Fields(newStr)
	for _,w :=range event.Words {
		event.TotalLengthWords+=len(w)
	}

}
