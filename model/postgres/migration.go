package postgres

import (
	//"fmt"
	log "github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
	//"gorm.io/gorm/clause"
	//"time"
)

// Setup automatically creates or updates the tables in Postgres
func Setup() {
	DbEvents.Migrator().DropTable( &Event{},
		&PassCond{},
		&State{},
		&KnowledgeDef{},
		&Scenario{},
		&TypeScenario{},
		&FailureCond{},
		&TypeCondFail{},
		&TypeErrorCond{},
		&EventIndex{},
		&Parameter{},
 	)
	err = DbEvents.AutoMigrate(
		&Event{},
		&PassCond{},
		&State{},
		&KnowledgeDef{},
		&Scenario{},
		&TypeScenario{},
		&FailureCond{},
		&TypeCondFail{},
		&TypeErrorCond{},
		&EventIndex{},
		&Parameter{},
		
	)

	log.Info("postgres: migrated tables")
}



func InitLoad() {

var TypeScenarios = []TypeScenario{{TypeScenarioName: "fatal_error"},{TypeScenarioName: "single_shot"},{TypeScenarioName: "dynamic_state"},}
DbEvents.Create(&TypeScenarios)	

var TypeCondFail = []TypeCondFail{{ConditionName: "error"},}
DbEvents.Create(&TypeCondFail)

var TypeErrorConds = []TypeErrorCond{{TypeErrorCondName: "error"}}
DbEvents.Create(&TypeErrorConds)

var PassConds = []PassCond{{PassCondName: "any_found"},{PassCondName: "all_found"}}
DbEvents.Create(&PassConds)

var KnowledgeDefs = []KnowledgeDef{KnowledgeDef{DefName: "Log Database", 
					Scenarios: []Scenario{
						Scenario{
							TypeScenarioID: 3, 
							States: []State{
								State{ PassCondID: 1, 
									Message: "Application %0%2 starting",
									Events: []Event{
										Event{
											Log: "Start proc %d:%s for activity {%s/%s}",
									 		Name: "Process started %3", 
									 		Process: "ActivityManager",
									 		},
									 	},
									},
								State{ PassCondID: 1, 
									Message: "Application %0%2 stopped",
									Events: []Event{
										Event{
											Log: "[%d,%d,%d,%0%2/%s]",
									 		Name: "Process %0%3 stopped", 
									 		Process: "am_destroy_activity",
									 		},
										Event{
											Log: "[%d,%d,%0%2/%s]",
									 		Name: "Process %0%3 stopped", 
									 		Process: "am_kill",
									 		},
										Event{
											Log: "[%d,%d,%0%2/%s]",
									 		Name: "Process %0%3 stopped", 
									 		Process: "am_stop_activity",
									 		},
									 	},

									},
								},

							},
						Scenario{
							TypeScenarioID: 3, 
							States: []State{
								State{ PassCondID: 1, 
									Message: "Displays are switched on",
									Events: []Event{
										Event{
											Log: "INF: Setting display power state DISPLAY_POWER_STATE_ON",
									 		Name: "Display switched to ON state", 
									 		Process: "MagicFlinger",
									 		},
										Event{
											Log: "INF: Setting display power state DISPLAY_POWER_STATE_ON_VBLANK",
									 		Name: "Display switched to ON state", 
									 		Process: "MagicFlinger",
									 		},
									 		
									 	},
									},
								State{ PassCondID: 1, 
									Message: "Displays are switched off",
									Events: []Event{
										Event{
											Log: "INF: Setting display power state DISPLAY_POWER_STATE_LED_OFF",
									 		Name: "Display switched OFF", 
									 		Process: "MagicFlinger",
									 		},
									 	},

									},
								},

							},
						Scenario{
							TypeScenarioID: 2, 
							FailureCond: FailureCond{ 
										TypeCondFailID: 1,
										TypeErrorCondID: 1,
										FailureMessage: "MagicFlinger didn't initialize correctly",
										},
							States: []State{
								State{ PassCondID: 2, 
									Message: "MagicFlinger successfully booted",
									Events: []Event{
										Event{
											Log: "Deferred initalization complete",
									 		Name: "MagicFlinger deferred intialization complete", 
									 		Process: "MagicFlinger",
									 		}, 		
									 	},
									},
								},

							},														
						},
					},
				}
DbEvents.Create(&KnowledgeDefs)

}

