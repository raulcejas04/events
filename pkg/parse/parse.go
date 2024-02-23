package parse

import (
    //"fmt"
    "strings"

)
type Event struct{
	LogLine string
	Words []string
	TotalLengthWords int
}

/*func main() {
	input := "Start proc 6578:com.android.provision/u0a160 for activity {com.android.provision/com.android.provision.DefaultActivity}"
	e:=Event{ LogLine: "Start proc %d:%s for activity {%s/%s}" }
	//split
	e.getWords()
	if e.approximate( input ) {
		fmt.Println( "it matched" )
	}
}*/

func (event *Event ) Approximate( line string ) bool {
	//fmt.Println( "Words ",event.Words )
	//fmt.Println( "Line ",line," totallengthword ", event.TotalLengthWords, len(line) )
	if event.TotalLengthWords > len(line) {
		return false
	}
	i:=0
	lenLine:=len( line)
	var matched []bool
	for _,word := range event.Words {
		//fmt.Println( " aprox 1 ", word )
		lenWord:=len(word)
		find:=true
		for j:=i;j<lenLine;j++ {
			if len(line[j:])<len(word) {
				find=false
				break
			}
			//fmt.Println( " aprox 2 ",word," j ",j," lenw ",lenWord," word cal -", line[j:j+lenWord],"-" );
			if line[j:j+lenWord]==word {
				//fmt.Println( "Matched")
				i=i+lenWord
				matched=append(matched, true )
				break
			} 
			
		}
		if !find {
			break
		}	
	}
	if len( event.Words ) == len(matched ) {
		for k,_ :=range event.Words {
			if !matched[k] {
				return false
			}
		}
		return true
	} else {
		return false
	}

}

func ( event *Event ) GetWords( ) {

    i:=0
    //j:=0
    
    for  {
    	var word string
    	cond:=event.LogLine[i:]
    	//fmt.Println( "cond ", cond )
    	start := strings.Index(cond, "%")
    	if start!=-1 {
    		word = cond[:start]
    	} else {
    		word = cond
    	}
    	
    	//fmt.Println( "word1 ", word, " i ",i, " start ", start )
    	if len(word)> 0 {
    		event.Words=append(event.Words,word)
    		event.TotalLengthWords+=len(word)
    		i=i+start+2
    	} else {
    		break
    	}

    }
 

}
