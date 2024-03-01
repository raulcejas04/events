package postgres


func GetKnowledgeDef( ID int ) *KnowledgeDef {

k KnowledgeDef
DbEvent.Find( k, ID )
return &k
}
