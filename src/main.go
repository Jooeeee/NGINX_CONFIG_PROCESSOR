package main

import (
	"fmt"
	"github.com/Jooeeee/NGINX_CONFIG_PROCESSOR/processor"
)

func main(){
	fmt.Println("This is Test!")
	cmds:=[]interface{processor.ADD,}
	processor.Processor("server.conf",)
}