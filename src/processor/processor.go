package processor

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"errors"
)

type PtType int

const (
	BLOCKSTART PtType=iota
	COMMAND
	BLOCKEND
)

type location struct {
	url string
	cmds map[string]string
}

type server struct{
	cmds map[string]string
	loc []location
}

func Scanner(filename string) (server,error){
	var srv server
	file,err:=os.OpenFile(filename,os.O_RDWR,0666)
	if err!=nil {
		return srv,err
	}
	defer file.Close()
	srv,err=serScanner(file)
	if err!=nil{
		return srv,err
	}
	return srv,nil
}

func (svr *server) dumper(filename string) {
	return
}

func (svr *server) clean(){
	svr.cmds=make(map[string]string)
	svr.loc=make([]location,0)
}

func serScanner(file *os.File)( server, error){
	var svr server
	svr.cmds=make(map[string]string)
	svr.loc=make([]location,0)
	bufScanner:=bufio.NewScanner(file)
	for bufScanner.Scan() {
		line:=strings.TrimSpace(bufScanner.Text())
		pt,key,value,err:=strParser(line)
		if err!=nil{
			fmt.Println(err)
		}else{
			switch pt {
			case BLOCKSTART:
				switch key{
				case "location":
					_,err:=locScanner(file)
					if err!=nil{
						return svr,err
					}
					// append(svr.loc,loc)
				default:
					return svr,errors.New("Error Block")
				}
			case COMMAND:
				svr.cmds[key]=value
			case BLOCKEND:
				return svr,nil			
			}
		}
	}
    return svr,errors.New("Error Block")
}

func locScanner(file *os.File) (location,error) {
	var loc location
	bufScanner:=bufio.NewScanner(file)
	for bufScanner.Scan() {
		line:=strings.TrimSpace(bufScanner.Text())
		pt,key,value,err:=strParser(line)
		if err!=nil{
			fmt.Println(err)
		}else{
			switch pt {			
			case COMMAND:
				loc.cmds[key]=value
			case BLOCKEND:
				return loc,nil		
			default:
				return loc,errors.New("Error Block")
			}
		}
	}
	return loc,errors.New("Error Block")
}

func strParser(line string) (pt PtType, key string,value string, err error) {
	key="abc"
	value="bcd"
	return
}

func strGenerator(key,value string) string {
	return key+" "+value
}
