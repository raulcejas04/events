package psql
import(
	"gorm.io/gorm/clause"
	log "github.com/sirupsen/logrus"
)

func ( k *KnowledgeDef ) GetFullKnowledgeDef( id uint )  {
DbEvents.Preload("Scenarios.States.Events").Preload(clause.Associations).Find( k, id )
return

}


func ( k *KnowledgeDef ) GetKnowledgeDef( ID uint )  {

result:=DbEvents.Find( k, ID )
if result.Error != nil {
	log.Errorf( "Error getting knowledgeDef %+v", result.Error )
}
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
result:=DbEvents.Find( s, id )
if result.Error != nil {
	log.Errorf( "Error getting state %+v", result.Error )
}
}

func ( s *Scenario ) GetScenario( id uint ) {
result:=DbEvents.Find( s, id )
if result.Error != nil {
	log.Errorf( "Error getting scenario %+v", result.Error )
}
}

func ( e *Event ) GetEvent( id uint ) {
result:=DbEvents.Find( e, id )
if result.Error != nil {
	log.Errorf( "Error getting event %+v", result.Error )
}
}

func (p *PassCond ) GetPassCond( id uint ) {
result:=DbEvents.Find( p, id )
if result.Error != nil {
	log.Errorf( "Error getting PassCond %+v", result.Error )
}
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

