package postgres

func GetFullExtraKnow( bugreportID int, knowledgeDefID int ) *ExtraKnow {

var e ExtraKnow
DbEvents.Where("bugreport_id=? and knowledge_def_id=?", bugreportID, knowledgeDefID ).Preload("ExtraScenarios.ExtraStates.ExtraEvents.ExtraParameters").Preload(clause.Associations).Find( &e )

return &e
}


