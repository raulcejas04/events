package postgres
import (
	//"fmt"
	log "github.com/sirupsen/logrus"
	"database/sql"
	"time" 
)

type FileMsg struct {
	FileId uint
	FileName string
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
		var file_id uint
		var file_name string
		if err := rows.Scan(&file_id,&file_name); err != nil {
			log.Fatal("Error in scan")
		}
		files = append(files, FileMsg{ FileId: file_id, FileName: file_name} )
	}
	//fmt.Println("files ", files)

	if len( files )==0 {
		return nil,0
	} else {
		return &files,partitionId
	}
}

type Message struct{
	BugreportId int
	PartitionId int
	Mess string
	Pid int
	Tid int
	Tag string
	Location string
	Timestamp time.Time
	BootId uint
	BootName string
	FileId uint
	FileName string
	LineNumber uint
}

type MsgBoot struct {
	BootId uint
	BootName string
}

func GetContents(bugreportId int, partitionId int, fileId uint, fileName string) []Message {

	rows,err := DbSql.Query("SELECT l.id,b.id as boot_folder_id,b.boot_folder_name FROM labels l, boot_folders b WHERE l.boot_folder_id=b.id AND bugreport_id=$1 and partition_id=$2", bugreportId, partitionId)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var labels =make( map[int]MsgBoot )
	for rows.Next() {
		var id int
		var boot_folder_id uint
		var boot_folder_name string
		if err := rows.Scan(&id,&boot_folder_id,&boot_folder_name); err != nil {
			log.Fatal(err)
		}
		labels[id]=MsgBoot{ BootId: boot_folder_id, BootName: boot_folder_name}
	}
	
	log.Info( "bugreport ", bugreportId, " partition_id ", partitionId, " file_id ", fileId )
	rows, err = DbSql.Query("SELECT pid,tid,line_number,tag,message,label_id,timestamp,location FROM contents WHERE bugreport_id=$1 and partition_id=$2 and file_id=$3 and message like 'Start proc%'", bugreportId, partitionId, fileId)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var messages []Message
	for rows.Next() {
		var pid int
		var tid int
		var line_number uint
		var message string
		var tag string
		var label_id int
		var timestamp time.Time
		var location string
		if err := rows.Scan(&pid,&tid,&line_number, &tag,&message,&label_id,&timestamp,&location); err != nil {
			log.Fatal(err)
		}
		boot:=labels[label_id]
		messages = append(messages, Message{ BugreportId: bugreportId, PartitionId: partitionId, Location: location, Pid:pid, Tid:tid, Tag: tag, Mess: message, Timestamp: timestamp, BootId: boot.BootId, BootName: boot.BootName, FileId: fileId, FileName: fileName, LineNumber: line_number })
	}
	//fmt.Println(" fileid 2 ", fileId, len(messages))
	return messages
}
