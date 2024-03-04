package postgres
import(
	"gorm.io/gorm/clause"
)

func ( k *KnowledgeDef ) GetFullKnowledgeDef( id uint )  {
DbEvents.Preload("Scenarios.States.Events").Preload(clause.Associations).Find( k, id )
return

}


func ( k *KnowledgeDef ) GetKnowledgeDef( ID uint )  {

DbEvents.Find( k, ID )
return 
}

func ( k *KnowledgeDef ) GetEvents() *map[uint]map[uint]map[uint]string {

var events = make( map[uint]map[uint]map[uint]string )
for _,s := range k.Scenarios {
	for _,st := range s.States {
		for _,e := range st.Events {
			if _,ok:=events[s.ID]; !ok {
				events[s.ID]=make( map[uint]map[uint]string )
			}
			if _,ok:=events[s.ID][st.ID]; !ok {
				events[s.ID][st.ID]=make( map[uint]string )
			}
			if _,ok:=events[s.ID][st.ID][e.ID]; !ok {
				events[s.ID][st.ID][e.ID]=e.Log
			}
		}
	}
}
return &events
}

func ( s *State ) GetState( id uint ) {
DbEvents.Find( s, id )
}

func ( s *Scenario ) GetScenario( id uint ) {
DbEvents.Find( s, id )
}


func GetAllTypeScenario() map[uint]string {

var TypeScenarios []TypeScenario
DbEvents.Find( &TypeScenarios )

res := make( map[uint]string )
for _,t := range TypeScenarios {
	res[t.ID]=t.TypeScenarioName
}
return res
}

func GetAllStates() map[uint]map[uint]State {

var states []State
DbEvents.Find( &states )

var res map[uint]map[uint]State
res=make(map[uint]map[uint]State)

for _,s := range states {
	if _,ok:=res[s.ScenarioID]; !ok {
		res[s.ScenarioID]=make(map[uint]State)
	} 
	res[s.ScenarioID][s.ID]=s  //keys scenarioId, stateId
}

return res
}

