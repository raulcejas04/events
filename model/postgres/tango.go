package postgres
import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"database/sql"	 
)

func GetFilesId( bugreportId int ) (*[]int, int) {

	rows,err := DbSql.Query("SELECT partition_id FROM br_partitions WHERE bugreport_id=$1",bugreportId )
	if err == sql.ErrNoRows {
		log.Error("No bugreport Found")
	}
	if err != nil {
		log.Fatal(err)
	}

	var partitionId int
	for rows.Next() {
		if err := rows.Scan(&partitionId); err != nil {
			log.Fatal("Error in scan")
		}
		break
	}

	log.Info( "For bugreport ", bugreportId," Partition found ", partitionId )

	rows, err = DbSql.Query("SELECT id FROM files WHERE bugreport_id=$1 and partition_id=$2", bugreportId, partitionId)
	if err == sql.ErrNoRows {
		log.Fatal("No Results Found")
	}
	if err != nil {
		log.Fatal(err)
	}

	var files []int
	for rows.Next() {
		var file_id int
		if err := rows.Scan(&file_id); err != nil {
			log.Fatal("Error in scan")
		}
		files = append(files, file_id)
	}
	fmt.Println("files ", files)

	if len( files )==0 {
		return nil,0
	} else {
		return &files,partitionId
	}
}

type Message struct{
	Mess string
	Tag string
}

func GetContents(bugreport_id int, partition_id int, file_id int) []Message {
	sql, err := DbSql.Prepare("SELECT tag,message FROM contents WHERE bugreport_id=$1 and partition_id=$2 and file_id=$3")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := sql.Query(bugreport_id, partition_id, file_id)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var messages []Message
	for rows.Next() {
		var message string
		var tag string
		if err := rows.Scan(&tag,&message); err != nil {
			log.Fatal(err)
		}
		messages = append(messages, Message{ Tag: tag, Mess: message})
	}
	fmt.Println(" fileid 2 ", file_id, len(messages))
	return messages
}
