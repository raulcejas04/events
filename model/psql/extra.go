package psql
import(
	"gorm.io/gorm/clause"
	"time"
	log "github.com/sirupsen/logrus"
	//"fmt"		 
	)

func GetFullExtraKnow( bugreportID int, knowledgeDefID uint, bootName string ) *ExtraKnow {

var e ExtraKnow
DbEvents.Where("bugreport_id=? and knowledge_def_id=?", bugreportID, knowledgeDefID ).Preload("ExtraScenarios.ExtraStates.ExtraEvents.ExtraParameters").Preload(clause.Associations).Find( &e )
// add filter by bootName
return &e
}

func NewExtraKnow( bugreportID int, partitionID int, knowledgeDefID uint, bootId uint, bootName string ) *ExtraKnow {

	//k:=GetKnowledgeDef( knowledgeDefID )
	ek:=ExtraKnow{ BugreportID: bugreportID, PartitionID: partitionID, KnowledgeDefID: knowledgeDefID, BootID: bootId, BootName: bootName, }
	return &ek
}

func ( ek *ExtraKnow ) GetScenarios() []ExtraScenario {
	return ek.ExtraScenarios
}

func ( sc *ExtraScenario ) GetStates() []ExtraState {
	return sc.ExtraStates
}

func ( st *ExtraState ) GetEvents() []ExtraEvent {
	return st.ExtraEvents
}

func ( ev *ExtraEvent ) GetParameters() []ExtraParameter {
	return ev.ExtraParameters
}


func ( ek *ExtraKnow ) AddExtraScenario( scenario ExtraScenario ) {
	//scenario.ExtraKnowID=ek.ID
	ek.ExtraScenarios=append( ek.ExtraScenarios, scenario )
} 

func NewExtraScenario(  scenarioID uint ) *ExtraScenario {
	es:=ExtraScenario{ ScenarioID: scenarioID }
	return &es
}

func ( ec *ExtraScenario ) AddExtraState( state ExtraState ) {
	//state.ExtraScenarioID=ec.ID
	ec.ExtraStates=append( ec.ExtraStates, state )
}

func NewExtraState(  stateID uint ) *ExtraState {
	et:=ExtraState{ StateID: stateID }
	return &et
}

func ( et *ExtraState ) AddExtraEvent( event ExtraEvent ) {
	et.ExtraEvents=append( et.ExtraEvents, event )
} 

func NewExtraEvent(  eventID uint, location string,  fileId uint, fileName string, lineNumber uint, pid int, tid int, timestamp time.Time, message string ) *ExtraEvent {
	ev:=ExtraEvent{ EventID: eventID, Location: location,  FileID: fileId, FileName: fileName, LineNumber: lineNumber, Pid: pid, Tid: tid, Timestamp: timestamp, LogMessage: message }
	return &ev
}

func ( ev *ExtraEvent ) AddExtraParameter( param ExtraParameter ) {
	ev.ExtraParameters=append( ev.ExtraParameters, param )
} 

func GetFatalErrors( bugReportId int, knowledgeDefId uint , bootName string ) *[]LvFailureCond {

//TODO use gorm to make this join
sql:= `select sc.id,sc.scenario_name,st.id as state_id,fc.failure_message,tf.condition_name,ts.type_scenario_name
	from type_scenarios ts, scenarios sc, failure_conds fc, type_cond_fails tf, type_states tys, states st
	left join ( extra_states est 
	inner join extra_scenarios esc on est.extra_scenario_id=esc.scenario_id
	inner join extra_knows ek on esc.extra_know_id=ek.id and ek.bugreport_id=$1 and ek.boot_name=$2 and
	ek.knowledge_def_id=$3) on st.id=est.state_id
	where 
	ts.id=sc.type_scenario_id and 
	sc.id = st.scenario_id and
	st.type_state_id=tys.id and
	sc.id=fc.scenario_id and
	fc.type_cond_fail_id=tf.id and
	tys.type_state_name='true' and 
	ts.type_scenario_name='single_shot'`

rows, err := DbEvents.Raw( sql, bugReportId, bootName, knowledgeDefId ).Rows()
defer rows.Close()
if err!=nil {
	log.Errorf("Error getting failure condition %+v\n",err )
	return nil
}

var lvFcs []LvFailureCond
for rows.Next() {
	var lvFc LvFailureCond
	rows.Scan(&lvFc.ScenarioID, &lvFc.ScenarioName, &lvFc.StateID, &lvFc.FailureMessage, &lvFc.CondName, &lvFc.TypeScenarioName)
	lvFcs=append(lvFcs,lvFc)
}

return &lvFcs


}


