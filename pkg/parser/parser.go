package parser

import (
	"fmt"
	"strings"
	"regexp"
	//"strconv"
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
	TypeScenarioId uint
	StateId uint
	EventId uint
	StartEnd string
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

/*
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
}*/

func NewParser( knowledgeDef uint ) *Parser {
	parser := Parser{}
	know := postgres.KnowledgeDef{}
	know.GetFullKnowledgeDef(knowledgeDef)
	fmt.Printf("know %+v\n",know)
	dbEvents:=*(know.GetEvents())
	var events []Event
	for scenarioId,scen := range dbEvents {
		scenarioRow:=postgres.Scenario{}
		scenarioRow.GetScenario(scenarioId)
		fmt.Printf( " scenario %d %d\n %+v\n", scenarioId, scenarioRow.TypeScenarioID, scenarioRow )
		for stateId,state := range scen {
			stateRow:=postgres.State{}
			stateRow.GetState( stateId )
			for eventId,e := range state {	
				event:= Event{ScenarioId: scenarioId, TypeScenarioId: scenarioRow.TypeScenarioID ,StateId: stateId, EventId: eventId, LogLine:e, StartEnd: stateRow.StartEnd}
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


func (parse *Parser ) GetStateStartEnd( scenarioId uint, stateId uint ) string {
	for _,e := range parse.Events {
		if e.ScenarioId==scenarioId && e.StateId==stateId {
			return e.StartEnd
		}
	}
	return ""
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
	
	fmt.Println("ItMatch 1 ", event.LlRegex, input )
	if pattern.MatchString(input) {
		fmt.Println("ItMatch 2 ", input )
		if strings.Contains(event.LogLine, "%s") || strings.Contains(event.LogLine, "%d") {

			fmt.Println("ItMatch 3 ", event.LogLine )
			parameters := pattern.FindStringSubmatch( input )
			if len(parameters)>0 {
				fmt.Println("ItMatch 4 ", parameters )
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



func ( event *Event ) RegularExpression() {
	event.LlRegex=strings.Replace( event.LogLine, `/`, `\/`, -1)
	event.LlRegex=strings.Replace( event.LlRegex, `[`, `\[`, -1)
	event.LlRegex=strings.Replace( event.LlRegex, `]`, `\]`, -1)	
	event.LlRegex=strings.Replace( event.LlRegex, "%s", "(.+)", -1)
	event.LlRegex=strings.Replace( event.LlRegex, "%d", "(\\d+)", -1)
	event.LlRegex=fmt.Sprintf("^%s$", event.LlRegex )
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

func ( event *Event ) GetParameters( input string ) *[]string {
	if strings.Contains(event.LogLine, "%s") || strings.Contains(event.LogLine, "%d") {
		pattern:= regexp.MustCompile( event.LlRegex )
		parameters := pattern.FindStringSubmatch( input )
		if len(parameters)>0 {
			fmt.Println( "PARAMETERS ", input, event.LlRegex, parameters )
			return &parameters
		}
	} else {
		return nil
	}
	return nil
}
