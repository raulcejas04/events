// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"strings"
)

type UsedParam struct {
	RefIndEvent int
	RefEvent string
	RefIndParam int
	RefParam string
}

func main() {
	//       012345678901234567890
	line := "[%d,%d,%d,%0%2/%s%0%1]"
	//line = "afasdfdas%0%2xsdfs"
	result := UsedParams(line)
	fmt.Printf("%+v\n", result)
}

func UsedParams(line string) []UsedParam {
	var res []UsedParam
	var begin int
	begin = 0
	//i:=0
	for {
		var index int = 0
		//fmt.Println( "begin ",begin," line ",line[begin:], len(line[begin:]) )

		if index = strings.Index(line[begin:], "%"); index != -1 {
			index=begin+index
			//fmt.Println( "1 ",  index, line[index+1:index+2], line[index+3:index+4] )
			if line[index+2:index+3] == "%" && !( line[index+1:index+2] == "s" || line[index+3:index+4] == "d" ) {
				//fmt.Println( "2 ", line[index+1:index+2], line[index+3:index+4] )
				refEvent := line[index+1 : index+2]
				refParam := line[index+3 : index+4]
				res = append(res, UsedParam{ RefIndEvent: index+1, RefEvent: refEvent, RefIndParam: index+3, RefParam: refParam })
			}
		}
		//fmt.Println( "index ",index," len ", len(line) )
		//fmt.Printf("\n")
		if index < len(line) && index!=-1 {
			if line[index+1:index+2]=="s" || line[index+1:index+2]=="d" {
				begin =index+2 
			} else {
				begin = index + 4
			}
		} else {
			break
		}
	}
	return res
}

