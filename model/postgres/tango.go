package postgres
import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"database/sql"
	 
)

type FileMsg struct {
	FileId int
	FileName int
}

func GetFilesId( bugreportId int ) (*[]FileMsg, int) {

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

	rows, err = DbSql.Query("SELECT id,file_name FROM files WHERE bugreport_id=$1 and partition_id=$2", bugreportId, partitionId)
	if err == sql.ErrNoRows {
		log.Fatal("No Results Found")
	}
	if err != nil {
		log.Fatal(err)
	}

	var files []FileMsg
	for rows.Next() {
		var file_id int
		var file_name string
		if err := rows.Scan(&file_id,&file_name); err != nil {
			log.Fatal("Error in scan")
		}
		files = append(files, FileMsg{ FileId: file_id, FileName: file_name )
	}
	//fmt.Println("files ", files)

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

type MsgBoot struct {
	BootId int
	BootName string
}

func GetContents(bugreport_id int, partition_id int, file_id int) []Message {

	rows,err := DbSql.Query("SELECT l.id,b.boot_folder_id,b.boot_folder_name FROM labels l, boot_folders b WHERE l.boot_folder_id=b.id AND 
				bugreport_id=$1 and partition_id=$2 and file_id=$3", bugreport_id, partition_id, file_id)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var labels =make( map[int]string )
	for rows.Next() {
		var id,boot_folder_id int
		var boot_folder_name string
		if err := rows.Scan(&id,&boot_id,&boot_folder_name); err != nil {
			log.Fatal(err)
		}
		labels[id]=MsgBoot{ BootId: boot_id, BootName: boot_folder_name}
	}
	
	log.Info( "bugreport ", bugreport_id, " partition_id ", partition_id, " file_id ", file_id )
	rows, err := DbSql.Query("SELECT tag,message,label_id,timestamp FROM contents WHERE bugreport_id=$1 and partition_id=$2 and file_id=$3", bugreport_id, partition_id, file_id)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var messages []Message
	for rows.Next() {
		var message string
		var tag string
		var label_id int
		var timestamp time.Time
		var location string
		if err := rows.Scan(&tag,&message,&label_id,&timestamp,&location); err != nil {
			log.Fatal(err)
		}
		boot:=labels[label_id]
		messages = append(messages, Message{ Location: location, Tag: tag, Mess: message, timestamp: timestamp, BootId: boot.BootId, BootName: boot.BootName, Timestamp: timestamp})
	}
	fmt.Println(" fileid 2 ", file_id, len(messages))
	return messages
}
