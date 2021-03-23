package main

import (
	"fmt"
	"github.com/Jooeeee/NGINX_CONFIG_PROCESSOR/processor"
)

func main(){
	fmt.Println("Welcome to nginx controller")
	// cmds:=[]interface{processor.ADD,"wdsss","kkkk","","",processor.MODLOC,"loc","loc","/",processor.ADD}
	err:=processor.Processor(
		"serveri.conf",processor.ADD,"listen","kkk","",0,
	processor.MODLOC,"loc","loc","/",processor.DELETE,
	processor.MODLOC,"loc","locddd","/",processor.MODIFY)
	if err!=nil{
		fmt.Println(err)
	}
}