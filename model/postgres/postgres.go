// %BANNER_BEGIN%
//---------------------------------------------------------------------
// %COPYRIGHT_BEGIN%
//
// Copyright (c) 2022 Magic Leap, Inc. (COMPANY) All Rights Reserved.
// Magic Leap, Inc. Confidential and Proprietary
//
// %COPYRIGHT_END%
//---------------------------------------------------------------------
// %BANNER_END%
//
// What is this?
// This package is responsible for saving bugreports to a PostgresSQL database
// This package uses gorm.io/gorm as an external dependency. Gorm is a full-featured
// ORM library for Golang that makes operations with SQL databases reliable and friendly.

package postgres

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	 "gorm.io/gorm/clause"
	 "time"
	"database/sql"
	_ "github.com/lib/pq"
	//sqldblogger "github.com/simukti/sqldb-logger"
   	//"github.com/rs/zerolog"	
	//zerologadapter "github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	//"os"		
)

type KnowledgeDef struct {
	gorm.Model
	DefName 	string
	Scenarios  	[]Scenario 	`gorm:"constraint:OnDelete:CASCADE;foreignkey:KnowledgeDefID;references:ID;"`
}

type Scenario struct {
	ID        		uint `gorm:"primaryKey"`
	KnowledgeDefID		uint
	TypeScenarioID		uint
	TypeScenario 		TypeScenario  	//`gorm:"foreignkey:ScenarioId;references:ID"`
	FailureCond  		FailureCond	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ScenarioID;references:ID;"`
	States       		[]State	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ScenarioID;references:ID;"`
}

type TypeScenario struct {
	ID               int     	`json:"id" gorm:"primaryKey;"`
	TypeScenarioName string
}

type FailureCond struct {
	ID        		uint `gorm:"primaryKey"`
	ScenarioID		uint
	TypeCondFailID		uint 
	TypeCondFail 	 	TypeCondFail 
	TypeErrorCondID	uint
	TypeErrorCond 		TypeErrorCond 
	FailureMessage  	string
}

type TypeCondFail struct {
	ID             uint    `json:"id" gorm:"primaryKey;"`
	ConditionName 	string
}

type TypeErrorCond struct {
	ID               	uint     	`json:"id" gorm:"primaryKey;"`
	TypeErrorCondName	string
}

type State struct {
	ID       	uint 		`gorm:"primaryKey"`
	ScenarioID	uint
	PassCondID	uint	
	PassCond	PassCond 	
	Message	string
	Events		[]Event	`gorm:"foreignkey:StateID;references:ID;constraint:OnDelete:CASCADE"`
}

type PassCond struct {
	ID               uint     	`json:"id" gorm:"primaryKey;"`
	PassCondName	string
}

type Event struct {
	ID        	uint 		`json:"id" gorm:"primaryKey"`
	StateID	uint	
	Name		string
	Process 	string
	Log		string
}

type EventIndex struct {
	ID		uint
	BugreportID	int
	PartitionID	int
	EventID	uint
	Event 		Event
	Location	string
	BootID		uint
	BootName	string
	FileID		uint
	FileName	string
	LineNumber 	uint
	Timestamp	time.Time
	Message	string
	Parameters	[]Parameter  `gorm:"foreignkey:EventIndexID;references:ID;constraint:OnDelete:CASCADE"`
}

func (EventIndex) TableName() string {
  return "event_index"
}

type Parameter struct {
	ID		uint
	EventIndexID	uint
	Offset		uint
	Value		string
}


// gorm instance
var DbEvents *gorm.DB
var err error

// newConnection makes a connection to the Postgres database. The connection parameters
// are set in the configuration file.
func NewConnection() {
	name := viper.GetString("postgres.name")
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")

	/*newLogger := logger.New(
	  logdebug.New(os.Stdout, "\r\n", logdebug.LstdFlags), // io writer
	  logger.Config{
	    SlowThreshold:              time.Second,   // Slow SQL threshold
	    LogLevel:                   logger.Info, // Log level
	    IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
	    Colorful:                  false,          // Disable color
	  },
	)*/

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)

	// open connection
	// Note: if we can't make a connection to the database no err is being returned
	//Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger,})
	DbEvents, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		log.Infof("postgres: successfully connected to database on %s %s %s", host, port, name)
	}

	//Setup()
}

var DbSql *sql.DB
/*type Logger interface {
	Log(ctx context.Context, level Level, msg string, data map[string]interface{})
}*/
func NewConnectionSql() {
	name := viper.GetString("postgres_tango.name")
	host := viper.GetString("postgres_tango.host")
	port := viper.GetString("postgres_tango.port")
	user := viper.GetString("postgres_tango.user")
	password := viper.GetString("postgres_tango.password")
	
	//dsn := "host=aggro.magicleap.ds user=aggro dbname=argus_tango sslmode=disable password=orgga port=5432"
	//dsn := "host=localhost user=postgres dbname=argus_prod_logcat sslmode=disable password=postgres port=6432"
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)
	fmt.Println( "dsn ", dsn )
	var err error
	DbSql, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Error: The data source arguments are not valid", err)
	} else {
		log.Infof("postgres: successfully connected to database")
	}
	
	//loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	//DbSql = sqldblogger.OpenDriver(dsn, DbSql.Driver(), loggerAdapter /*, ...options */) 	

}


func ( k *KnowledgeDef ) GetKnowledgeDef( id int )  {


DbEvents.Preload("Scenarios.States.Events").Preload(clause.Associations).Find( k, id )

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
