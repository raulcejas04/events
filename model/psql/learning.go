package psql
import(
	//"gorm.io/gorm"
	"fmt"
	//"gorm.io/gorm/clause"
)

type LogViewer struct {
	DefName			string
	LVScenarios			[]LVScenario
}

type LVScenario struct {
	Start 	LVState
	End	LVState
}

type LVState struct {
	Message	string
	PassCondName	string
	LVEvents	[]LVEvent
	LVFailCond	LVFailCond
}

type LVFailCond struct {
	FailureMessage	string
}

type LVEvent struct {
	Name		string //name with parameters
	ExtraEvent	ExtraEvent
	LVParameters	[]LVParameter
}

type LVParameter struct {
	Value		string
}


func GetLearning( bugreportID int, knowledgeDefID uint, bootName string ) *LogViewer {

var k KnowledgeDef
(&k).GetKnowledgeDef( knowledgeDefID )

e:=GetFullExtraKnow( bugreportID, knowledgeDefID, bootName )

var lvScenarios []LVScenario
for _,sce := range e.GetScenarios() {
	fmt.Printf("scen %+v\n",sce )
	var startState,endState LVState
	for _,st :=range sce.GetStates() {
		fmt.Printf("stat %+v\n",st )
		
		var lvEvents []LVEvent
		for _,ev :=range st.GetEvents() {
			fmt.Printf("ev %+v\n",ev )
			var lvParameters []LVParameter
			for _,par :=range ev.GetParameters() {
				fmt.Printf("par %+v\n",par )
				lvParameters=append(lvParameters, LVParameter{ Value: par.Value } )
			}
			knowEvent := Event{}
			(&knowEvent).GetEvent( ev.ID )
			lvEvents=append( lvEvents, LVEvent{ Name: knowEvent.Name, ExtraEvent:ev, LVParameters: lvParameters } )
		}
		//TODO failcond
		var knowState = State{}
		(&knowState).GetState( st.ID )
		
		var passCond = PassCond{}
		(&passCond).GetPassCond( knowState.PassCondID )

		message:=replaceParams( knowState.Message )
		lvState := LVState{ Message: message, PassCondName: passCond.PassCondName, LVEvents: lvEvents }
		if knowState.StartEnd==START {
			startState=lvState
		} else {
			endState=lvState		
		}
	}
	lvScenarios = append( lvScenarios, LVScenario{ Start: startState, End: endState } )
}

return &LogViewer{ DefName: (k).DefName, LVScenarios: lvScenarios }
}


func replaceParams( message string ) string {
//TODO
return message
}

/*func (l *Learning ) makeFailureConds() {
	scenarios:=l.getScenarios()
	for _,scenario :=range scenarios {
		states:=l.getStates( scenario.ID )
		for _,state := range states {
			
		
	
		}
	}
}*/

