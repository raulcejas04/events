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

package psql

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	logdebug "log"		
	// "gorm.io/gorm/clause"
	 "time"
	"database/sql"
	_ "github.com/lib/pq"
	//sqldblogger "github.com/simukti/sqldb-logger"
   	//"github.com/rs/zerolog"	
	//zerologadapter "github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"os"		
)

const (
	START = "S"
	END   = "E" 
)

type KnowledgeDef struct {
	ID        	uint `gorm:"primaryKey"`
	gorm.Model
	DefName 	string
	Scenarios  	[]Scenario 	`gorm:"constraint:OnDelete:CASCADE;foreignkey:KnowledgeDefID;references:ID;"`
}

type Scenario struct {
	ID        		uint `gorm:"primaryKey"`
	ScenarioName		string
	KnowledgeDefID		uint
	TypeScenarioID		uint
	TypeScenario 		TypeScenario  	//`gorm:"foreignkey:ScenarioId;references:ID"`
	FailureCond  		FailureCond	//`gorm:"constraint:OnDelete:CASCADE;foreignkey:ScenarioID;references:ID;"`
	States       		[]State	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ScenarioID;references:ID;"`
}

type TypeScenario struct {
	ID               uint     	`json:"id" gorm:"primaryKey;"`
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
	TypeStateID	uint
	TypeState	TypeState 	
	Message	string
	StartEnd	string		//true, false or whatever in the future others	
	Events		[]Event	`gorm:"foreignkey:StateID;references:ID;constraint:OnDelete:CASCADE"`
}

type TypeState struct {
	ID               uint     	`json:"id" gorm:"primaryKey;"`
	TypeStateName	string
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

type ExtraKnow struct {
	ID        		uint 		`json:"id" gorm:"primaryKey"`
	BugreportID		int
	PartitionID		int
	BootID			uint
	BootName		string
	KnowledgeDefID		uint
	KnowledgeDef		KnowledgeDef
	ExtraScenarios  	[]ExtraScenario 	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ExtraKnowID;references:ID;"`	
}

type ExtraScenario struct {
	ID        		uint 		`json:"id" gorm:"primaryKey"`
	ExtraKnowID    		uint
	ScenarioID		uint	
	Scenario		Scenario	 		
	ExtraStates  		[]ExtraState	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ExtraScenarioID;references:ID;"`
}

type ExtraState struct {
	ID        		uint 		`json:"id" gorm:"primaryKey"`
	StateID			uint	
	State			State
	ExtraScenarioID    	uint 		
	ExtraEvents	  	[]ExtraEvent	`gorm:"constraint:OnDelete:CASCADE;foreignkey:ExtraStateID;references:ID;"`
}

//Calculated entity
type ExtraFailCond struct {
	ExtraStateID 	uint
	FailureCondID 	uint	
}

type ExtraEvent struct {
	ID		uint
	ExtraStateID  uint
	EventID	uint
	Event 		Event
	Location	string
	FileID		uint
	FileName	string
	LineNumber 	uint
	Pid		int
	Tid		int
	Tag		string
	Priority	string
	Timestamp	time.Time
	LogMessage	string
	ExtraParameters	[]ExtraParameter  `gorm:"foreignkey:ExtraEventID;references:ID;constraint:OnDelete:CASCADE"`
}

/*func (EventIndex) TableName() string {
  return "event_index"
}*/

type ExtraParameter struct {
	ID		uint
	ExtraEventID	uint
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

	newLogger := logger.New(
	  logdebug.New(os.Stdout, "\r\n", logdebug.LstdFlags), // io writer
	  logger.Config{
	    SlowThreshold:              time.Second,   // Slow SQL threshold
	    LogLevel:                   logger.Info, // Log level
	    IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
	    Colorful:                  false,          // Disable color
	  },
	)

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)

	// open connection
	// Note: if we can't make a connection to the database no err is being returned
	DbEvents, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger,})
	//DbEvents, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

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



