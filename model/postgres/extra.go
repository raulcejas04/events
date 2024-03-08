package postgres
import(
	 "gorm.io/gorm/clause"
	 "time"
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
	ev:=ExtraEvent{ EventID: eventID, Location: location,  FileID: fileId, FileName: fileName, LineNumber: lineNumber, Pid: pid, Tid: tid, Timestamp: timestamp, Message: message }
	return &ev
}

func ( ev *ExtraEvent ) AddExtraParameter( param ExtraParameter ) {
	ev.ExtraParameters=append( ev.ExtraParameters, param )
} 
