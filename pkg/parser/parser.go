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
	Words []string
	TotalLengthWords int
	UsedParams []UsedParam
}

type Parser struct {
	knowledgeDef postgres.KnowledgeDef
	Events []Event

}

/*func main() {
	input := "Start proc 6578:com.android.provision/u0a160 for activity {com.android.provision/com.android.provision.DefaultActivity}"
	e:=Event{ LogLine: "Start proc %d:%s for activity {%s/%s}" }
	//split
	e.getWords()
	if e.approximate( input ) {
		fmt.Println( "it matched" )
	}
}*/


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

func NewParser() *Parser {
	parser := Parser{}
	know := postgres.KnowledgeDef{}
	know.GetKnowledgeDef(1)
	dbEvents:=*(know.GetEvents())
	var events []Event
	for scenarioId,scen := range dbEvents {
		for stateId,state := range scen {
			for eventId,e := range state {	
				event:= Event{ScenarioId: scenarioId, StateId: stateId, EventId: eventId, LogLine:e}
				event.GetWords()
				event.RegularExpression()
				events=append(events, event)
			}
		}
	}
	parser.knowledgeDef=know
	parser.Events=append(parser.Events, events... )
	return &parser
}

func (parser *Parser ) Approximate( line string ) bool {
	//fmt.Println( "Words ",event.Words )
	//fmt.Println( "Line ",line," totallengthword ", event.TotalLengthWords, len(line) )
	for _,event := range  parser.Events {
		if event.TotalLengthWords > len(line) {
			return false
		}
 
		i:=0
		lenLine:=len( line)
		var matched []bool
		for _,word := range event.Words {
			//fmt.Println( " aprox 1 ", word )
			lenWord:=len(word)
			find:=true
			for j:=i;j<lenLine;j++ {
				if len(line[j:])<len(word) {
					find=false
					break
				}
				//fmt.Println( " aprox 2 ",word," j ",j," lenw ",lenWord," word cal -", line[j:j+lenWord],"-" );
				if line[j:j+lenWord]==word {
					//fmt.Println( "Matched")
					i=i+lenWord
					matched=append(matched, true )
					break
				} 
			
			}
			if !find {
				break
			}	
		}
		if len( event.Words ) == len(matched ) {
			for k,_ :=range event.Words {
				if !matched[k] {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}
	return false
}

func ( event *Event ) GetParameters( input string ) *[]string {
	if strings.Contains(event.LogLine, "%s") || strings.Contains(event.LogLine, "%d") {
		pattern:= regexp.MustCompile( event.LlRegex )
		parameters := pattern.FindStringSubmatch( input )
		if len(parameters)>0 {
			fmt.Println( "PARAMETERS ", parameters )
			return &parameters
		}
	} else {
		return nil
	}
	return nil
}

func ( event *Event ) RegularExpression() {
	event.LlRegex=strings.Replace( event.LogLine, "%s", "(.+)", -1)
	event.LlRegex=strings.Replace( event.LlRegex, "%d", "(\\d+)", -1)
	return
}

func ( event *Event ) GetWords( ) {

	re := regexp.MustCompile("%(d|s|\\d)")
	newStr := re.ReplaceAllString(event.LogLine, "")
	event.Words=strings.Fields(newStr)
	for _,w :=range event.Words {
		event.TotalLengthWords+=len(w)
	}

}
