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
	EMPTYLINE
)

type location struct {
	url string
	cmds map[string]string
}

type server struct{
	cmds map[string]string
	loc map[string]location
}

func Scanner(filename string) (server,error){
	var srv server
	file,err:=os.OpenFile(filename,os.O_RDWR,0666)
	if err!=nil {
		return srv,err
	}
	defer file.Close()
	bufScanner:=bufio.NewScanner(file)
	srv,err=serScanner(bufScanner)
	if err!=nil{
		return srv,err
	}
	return srv,nil
}

func Dumper(svr server,filename string) error{
	file,err:=os.OpenFile(filename,os.O_WRONLY|os.O_CREATE,0666)
	if err!=nil{
		return err
	}
	defer file.Close()

	bufWriter:=bufio.NewWriter(file)
	err=serDumper(svr,bufWriter)
	if err!=nil{
		return err
	}

	return nil
}

// func (svr *server) clean(){
// 	svr.cmds=make(map[string]string)
// 	svr.loc=make(map[string]location)
// }

func serScanner(bufScanner *bufio.Scanner)( server, error){
	var svr server
	svr.cmds=make(map[string]string)
	svr.loc=make(map[string]location)
	for bufScanner.Scan() {
		line:=strings.TrimSpace(bufScanner.Text())
		pt,key,value,err:=strParser(line)
		if err!=nil{
			fmt.Println(err)
		}else{
			switch pt {
			case BLOCKSTART:
				switch key{
				case "server":
					fmt.Println("server block start")
				case "location":
					loc,err:=locScanner(bufScanner)
					if err!=nil{
						return svr,err
					}
					svr.loc[value]=loc
				default:
					return svr,errors.New("Error Block")
				}
			case COMMAND:
				svr.cmds[key]=value
			case BLOCKEND:
				return svr,nil			
			// case EMPTYLINE:
			// 	continue
			}
		}
	}
    return svr,errors.New("Error Block")
}

func serDumper(svr server,bufWriter *bufio.Writer) error {
	defer bufWriter.Flush()
	bufWriter.WriteString("server {\n")
	defer bufWriter.WriteString("}\n")
	
	for k,v:=range svr.cmds{
		_,err:=bufWriter.WriteString(strGenerator(k,v))
		if err!=nil{
			return err
		}
	}
	bufWriter.Flush()
	for k,v:=range svr.loc{
		err:=locDumper(k,v,bufWriter)
		if err!=nil{
			return err
		}
		bufWriter.Flush()
	}
	return nil
}

func locScanner(bufScanner *bufio.Scanner) (location,error) {
	var loc location
	loc.cmds=make(map[string]string)
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
	if err:=bufScanner.Err();err!=nil{
		fmt.Println(err)
	}
	return loc,errors.New("Error Block")
}

func locDumper(key string,loc location,bufWriter *bufio.Writer) error{
	defer bufWriter.Flush()
	bufWriter.WriteString("location "+key+" {\n")
	defer bufWriter.WriteString("}\n")

	for k,v:=range loc.cmds{
		_,err:=bufWriter.WriteString(strGenerator(k,v))
		if err!=nil{
			return err
		}
	}
	return nil
}

func strParser(line string) (pt PtType, key string,value string, err error) {
	if len(line)==0{
		return EMPTYLINE,"","",nil
	}
	words:=strings.Fields(line)
	if len(words)>0 {
		switch words[0] {
		case "server":
			return BLOCKSTART, words[0],"",nil
		case "location":
			return BLOCKSTART, words[0],strings.TrimSuffix(words[1],";"),nil
		case "}":
			return BLOCKEND,"","",nil
		default:
			return COMMAND,words[0],strings.TrimSuffix(words[1],";"),nil
		}
	}
	return EMPTYLINE,"","",nil
}

func strGenerator(key,value string) string {
	return key+" "+value+";\n"
}
