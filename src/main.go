package main

import (
	"fmt"
	"github.com/Jooeeee/NGINX_CONFIG_PROCESSOR/processor"
)

func main(){
	fmt.Println("This is Test!")
	// cmds:=[]interface{processor.ADD,"wdsss","kkkk","","",processor.MODLOC,"loc","loc","/",processor.ADD}
	err:=processor.Processor(
		"server.conf",processor.ADD,"listen","kkk","",0,
	// processor.MODLOC,"loc","loc","/te",processor.ADD,
	// processor.MODLOC,"locasdf","loc","/te",processor.ADD
)
	fmt.Println(err)
}