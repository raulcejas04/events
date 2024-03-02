package postgres
import(
	//"gorm.io/gorm"
	//"fmt"
	//"gorm.io/gorm/clause"
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


func NewLearning( bugreportID int, partitionID int, knowledgeDefID uint ) *Learning {
	e:=ExtraKnow{ BugreportID: bugreportID, PartitionID: partitionID, KnowledgeDefID: knowledgeDefID }
	var k *KnowledgeDef
	k.GetKnowledgeDef( knowledgeDefID )
	l := Learning{ ExtraKnow: e, DefName: (*k).DefName }
	return &l
}


func GetLearning( bugreportID int, knowledgeDefID uint, bootName string ) *Learning {

var k *KnowledgeDef
k.GetKnowledgeDef( knowledgeDefID )
e:=GetFullExtraKnow( bugreportID, knowledgeDefID, bootName )

return &Learning{ ExtraKnow: *e, DefName: (*k).DefName }
}

/*func (l *Learning ) makeFailureConds() {
	scenarios:=l.getScenarios()
	for _,scenario :=range scenarios {
		states:=l.getStates( scenario.ID )
		for _,state := range states {
			
		
	
		}
	}
}*/

