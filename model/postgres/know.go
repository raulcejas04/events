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
