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

type ACTION int

const(
	ADD ACTION=iota
	MODIFY 
	DELETE
	MODLOC
)

type location struct {
	cmds map[string]string
}

func (loc *location) add_cmd(key,value string) {
	loc.cmds[key]=value
}

func (loc *location) modify_cmd(key, value string) {
	loc.add_cmd(key,value)
}

func (loc *location) del_cmd(key string) {
	delete(loc.cmds,key)
}


type server struct{
	cmds map[string]string
	loc map[string]*location
}

func (svr *server) add_cmd(key,value string) {
	svr.cmds[key]=value
}

func (svr *server) modify_cmd(key, value string) {
	svr.add_cmd(key,value)
}

func (svr *server) del_cmd(key string) {
	delete(svr.cmds,key)
}

func(svr *server) modify_loc(act ACTION,loc_key,key,value string) error{
	if _,ok:=svr.loc[loc_key];!ok{
		svr.loc[loc_key]=&location{cmds:make(map[string]string)}
	}
	loc:=svr.loc[loc_key]
	switch act{
	case ADD:
		loc.add_cmd(key,value)
	case MODIFY:
		loc.modify_cmd(key,value)
	case DELETE:
		loc.del_cmd(key)
	default:
		return errors.New("action flag doesn't exit")
	}
	return nil
}

func Scanner(filename string) (server,error){
	var svr server
	file,err:=os.OpenFile(filename,os.O_RDWR,0666)
	if err!=nil {
		svr.cmds=make(map[string]string)
		svr.loc=make(map[string]*location)
		return svr,err
	}
	defer file.Close()
	bufScanner:=bufio.NewScanner(file)
	svr,err=serScanner(bufScanner)
	if err!=nil{
		return svr,err
	}
	return svr,nil
}

// action,key,value,loc_key,loc_action
func Processor(filename string,opts... interface{}) error {
	if len(opts)%5!=0{
		return errors.New("options errors")
	}
	svr,_:=Scanner(filename)
	for i:=0;i<len(opts);i+=5{
		fmt.Println(opts[i:i+5])
		switch opts[i]{
		case ADD:
			svr.add_cmd(opts[i+1].(string),opts[i+2].(string))
		case MODIFY:
			svr.modify_cmd(opts[i+1].(string),opts[i+2].(string))
		case DELETE:
			svr.del_cmd(opts[i+2].(string))
		case MODLOC:
			err:=svr.modify_loc(opts[i+4].(ACTION),opts[i+3].(string),opts[i+1].(string),opts[i+2].(string))
			if err!=nil{
				return err
			}
		}
	}
	fmt.Println(svr.loc["/te"])
	err:=Dumper(svr,"./new.conf")
	if err!=nil{
		return err
	}
	return nil;
}

func Dumper(svr server,filename string) error{
	filetmp:=filename+"tmp"
	file,err:=os.OpenFile(filetmp,os.O_CREATE,0666)
	if err!=nil{
		return err
	}
	file.Close()

	bufWriter:=bufio.NewWriter(file)
	err=serDumper(svr,bufWriter)
	if err!=nil{
		file.Close()
		return err
	}
	file.Close()
	err=os.Rename(filetmp,filename)
	if err!=nil{
		return err
	}
	return nil
}


func serScanner(bufScanner *bufio.Scanner)( server, error){
	var svr server
	svr.cmds=make(map[string]string)
	svr.loc=make(map[string]*location)
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
					svr.loc[value]=&loc
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
	for k,v:=range svr.loc{
		err:=locDumper(k,v,bufWriter)
		if err!=nil{
			return err
		}
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

func locDumper(key string,loc *location,bufWriter *bufio.Writer) error{
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
