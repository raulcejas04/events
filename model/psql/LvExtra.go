package psql
import(
	//"gorm.io/gorm"
	//"fmt"
	//"gorm.io/gorm/clause"
	"regexp"
	"strings"
	"strconv"	
)

type LvExtra struct {		//LogViewer Business Object
	DefName			string
	LvScenarios		[]LvScenario
	LvErrors		[]LvFailureCond	
}

type LvScenario struct {
	Type	string
	Start 	LvState
	End	LvState
}

type LvState struct {
	MessageFormat	string
	Message		string
	PassCondName	string
	LvEvents	[]LvEvent
}

type LvFailureCond struct {
	ScenarioID		uint
	ScenarioName		string
	StateID		uint
	FailureMessage		string
	CondName		string
	TypeScenarioName	string
	Events			[]string
}

type LvEvent struct {
	Name		string //name with parameters
	ExtraEvent	ExtraEvent
	LvParameters	[]LvParameter
}

type LvParameter struct {
	Value		string
}


func GetLvExtra( bugreportID int, knowledgeDefID uint, bootName string ) *LvExtra {

var k KnowledgeDef
(&k).GetKnowledgeDef( knowledgeDefID )

e:=GetFullExtraKnow( bugreportID, knowledgeDefID, bootName )

var lvScenarios []LvScenario
for _,sce := range e.GetScenarios() {

	var startState,endState LvState
	var lvEvents []LvEvent	
	var knowState = State{}
	var lastState uint	
	for _,st :=range sce.GetStates() {
		lastState=st.StateID


		var lvParameters []LvParameter
		lvEvents=[]LvEvent{}
		for _,ev :=range st.GetEvents() {
			lvParameters=[]LvParameter{}
			for _,par :=range ev.GetParameters() {
				lvParameters=append(lvParameters, LvParameter{ Value: par.Value } )
			}
			knowEvent := Event{}
			(&knowEvent).GetEvent( ev.EventID )
			ev.Event=knowEvent
			lvEvents=append( lvEvents, LvEvent{  Name: knowEvent.Name, ExtraEvent:ev, LvParameters: lvParameters } )
		}
		//TODO failcond

		knowState= State{}
		(&knowState).GetState( st.StateID )

		if knowState.StartEnd==START {
		
			var passCond = PassCond{}
			(&passCond).GetPassCond( knowState.PassCondID )

			message:=replaceParams( knowState.Message, lvEvents )
			startState=LvState{ MessageFormat:knowState.Message, Message: message, PassCondName: passCond.PassCondName, LvEvents: lvEvents }
		}
	}

	if knowState.StartEnd==END {
		(&knowState).GetState( lastState )
		
		var passCond = PassCond{}
		(&passCond).GetPassCond( knowState.PassCondID )
	
		message:=replaceParams( knowState.Message, startState.LvEvents )
		endState = LvState{ MessageFormat: knowState.Message, Message: message, PassCondName: passCond.PassCondName, LvEvents: lvEvents }	
	}
	scenario:=Scenario{}
	(&scenario).GetScenario( sce.ScenarioID )
	typeScenario:=TypeScenario{}
	typeScenario.GetTypeScenario( scenario.TypeScenarioID )
	lvScenarios = append( lvScenarios, LvScenario{ Type: typeScenario.TypeScenarioName, Start: startState, End: endState } )
}

lvExtra:=LvExtra{ DefName: (k).DefName, LvScenarios: lvScenarios }
lvExtra.addFailureConds(bugreportID, knowledgeDefID, bootName )
return &lvExtra
}


func replaceParams( message string, events []LvEvent ) string {
//TODO
	re := regexp.MustCompile( "%\\d+%\\d+" )
	params:=re.FindStringSubmatch( message ) 
	for _,param := range params {
		re := regexp.MustCompile("%(\\d+)%(\\d+)")
		match := re.FindStringSubmatch(param)
		indexEvent,_:=strconv.Atoi(match[1])
		indexParam,_:=strconv.Atoi(match[2])
		if len(events)>indexEvent && len(events[indexEvent].LvParameters) > indexParam {
			value:= events[indexEvent].LvParameters[indexParam].Value
			message=strings.Replace( message, param, value, 1 )
			//fmt.Println( "replace2 ", message, " param ",param," match ", match, " event ",indexEvent,indexParam, " value ", value )					
		}

	}

return message
}

func (l *LvExtra ) addFailureConds( bugReportID int,  knowledgeDefId uint, bootName string ) {

lvFailureConds:=*GetFatalErrors( bugReportID, knowledgeDefId, bootName )
for k, lvFailureCond := range lvFailureConds {
	state:=State{ID: lvFailureCond.StateID }
	for _,ev :=range *(state.GetEvents()) {
		lvFailureConds[k].Events=append(lvFailureConds[k].Events, ev.Log )
	}
}

(*l).LvErrors=append( (*l).LvErrors, lvFailureConds...)
return
}

