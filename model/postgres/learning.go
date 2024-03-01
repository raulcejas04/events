package postgres
import(
	//"gorm.io/gorm"
	//"fmt"
	"gorm.io/gorm/clause"
)

type Learning struct {
	ExtraKnow
	DefName			string
	LearningScenarios	[]LearningScenario
}

type LearningScenario struct {
	ExtraScenario
	ScenarioName		string
	LearningStates		[]LearningState
}

type LearningState struct {
	ExtraState
	ExtraFailCond			//Calculated Fail Condition
	Message			string //message with parameters
	LearningEvents		[]LearningEvent
}

type LearningEvent struct {
	ExtraEvent
	Process			string
	Name			string //name with parameters
	LearningParameters	[]LearningParameter
}

type LearningParameter struct {
	ExtraParameter
}


func GetLearning( bugreportID int, knowledgeDefID int ) *Learning {

k:=GetKnowledgeDef( knowledgeDefID )
e:=GetExtraKnow( bugreportID, knowledgeDefID )

return &Learning{ ExtraKnow: e, DefName: k.DefName }
}

/*func (l *Learning ) makeFailureConds() {
	scenarios:=l.getScenarios()
	for _,scenario :=range scenarios {
		states:=l.getStates( scenario.ID )
		for _,state := range states {
			
		
	
		}
	}
}*/

