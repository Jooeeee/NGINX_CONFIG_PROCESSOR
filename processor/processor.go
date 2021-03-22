package processor

import (
	"io"
	"fmt"
	"ioutil"
)
type PtType int

const (
	BLOCKSTART PtType
	COMMAND PtType
	BLOCKEND PtType
)
func strParser(line string) (pt PtType, key string,value string,  error) {
	key="abc"
	value="bcd"
	return
}

func strGenerator(key,value string) string {
	return key+" "+value
}

type location struct {
	url string
	cmds map[string]string
}

type server struct{
	cmds map[string]string
	loc []location
}

func serScanner(buf *Reader) server, error{
	var svr server
	for{
		line,err:=buf.ReadString('\n')
		if err!=nil {
			if err==io.EOF {
				return svr,nil
			} else {
				return nil,err
			}
		}

		line=strings.TrimSpace(line)
		pt,key,value,err:=strParser(line)
		if err!=nill{
			fmt.Println(err)
		}else{
			switch pt {
			case BLOCKSTART:
				switch key{
				case "location":
					loc:=locScanner(file)
					append(svr.loc,loc)
				default:
					return nil,errors.New("Error Block")
				}
			case COMMAND:
				svr.cmds[key]=value
			case BLOCKEND:
				return svr,nil			
			}
		}
	}
    return svr,nil
}

func (svr *server) dumper(filename string) {
	return
}

func (svr *server) clean(){
	svr.cmds=make(map[string]string)
	svr.loc=make([]location,0)
}

func scanner(filenam string)server,error{
	file,err:=os.OpenFile(filename,os.O_RDWR,0666)
	if err!=nil {
		return err
	}
	defer file.Close()

	buf:=bufio.NewReader(file)
}