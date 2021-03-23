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
	svr,err:=Scanner(filename)
	if err!=nil{
		fmt.Println(err)
		fmt.Println("Create a new one")
	}
	for i:=0;i<len(opts);i+=5{
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
	err=Dumper(svr,filename)
	if err!=nil{
		return err
	}
	return nil;
}

func ProcessorTerm(){
	fmt.Print("Conf Filepath to create or edit:")
	var filename string
	fmt.Scanln(&filename)
	svr,err:=Scanner(filename)
	if err!=nil{
		fmt.Println("Load file error",err)
		fmt.Println("Create a new one")
	}
	fmt.Println(`Please enter command (format: ACTION KEY VALUE LOCATION_URL LOCATION_ACTION):
		or enter "SAVE" to save the conf and exit.
		- ACTION: ADD,MODIFY,DELETE,LOCATION
		- KEY: command key string
		- VALUE: comand value string. set any no empty string for delete
		- LOCATION_URL: location url string if ACTION is 3. Otherwise set to empty
		- LOCATION_ACTION:ADD,MODIFY,DELETE`)
	input:=bufio.NewReader(os.Stdin)
	for{		
		str,err:=input.ReadString('\n')
		if err!=nil {
			panic(err)
		}
		str=strings.TrimSpace(strings.TrimSuffix( str,"\n"))
		opts:=strings.Fields(str)
		if len(opts)==1 && opts[0]=="SAVE"{
			fmt.Println(svr)
			err=Dumper(svr,filename)
			if err!=nil{
				fmt.Println(err)
			}
			fmt.Println("Finish!")
			break
		}else if len(opts)==3{			
			switch opts[0]{
			case "ADD":
				svr.add_cmd(opts[1],opts[2])
			case "MODIFY":
				svr.modify_cmd(opts[1],opts[2])
			case "DELETE":
				svr.del_cmd(opts[2])
			default:
				fmt.Println("Unknow command")
			}
		}else if len(opts)==5{
			opts_switch:
				switch opts[0]{
				case "ADD":
					svr.add_cmd(opts[1],opts[2])
				case "MODIFY":
					svr.modify_cmd(opts[1],opts[2])
				case "DELETE":
					svr.del_cmd(opts[2])
				case "LOCATION":
					var lopt ACTION
					switch opts[4]{
					case "ADD":
						lopt=ADD
					case "MODIFY":
						lopt=MODIFY
					case "DELETE":
						lopt=DELETE
					default:
						fmt.Println("Unkown command")
						break opts_switch
					}					
					err:=svr.modify_loc(lopt,opts[3],opts[1],opts[2])
					if err!=nil{
						fmt.Println(err)
					}
				default:
					fmt.Println("Unknow command")
					break opts_switch
				}
		}else{
			fmt.Println("Unknow command")
		}	
	}
}

func Dumper(svr server,filename string) error{
	filetmp:=filename+"tmp"
	file,err:=os.OpenFile(filetmp,os.O_CREATE,0666)
	if err!=nil{
		file.Close()
		return err
	}

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
					continue
					// fmt.Println("server block start")
				case "location":
					loc,err:=locScanner(bufScanner)
					if err!=nil{
						return svr,err
					}
					svr.loc[value]=&loc
				default:
					return svr,errors.New("error block")
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
    return svr,errors.New("error block")
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
				return loc,errors.New("error block")
			}
		}
	}
	if err:=bufScanner.Err();err!=nil{
		fmt.Println(err)
	}
	return loc,errors.New("error block")
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
