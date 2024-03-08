package debug
import(
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
)

func NewOpenedFile( name string ) *os.File {

	f, err := os.OpenFile( name, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return f
}

func FileClose(f *os.File) {

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func FileWrite( f *os.File, sentence string ) {

	//fmt.Println(" sentence ", sentence )
	_, err := f.Write( []byte(sentence) )
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	return
}

