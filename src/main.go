package main

import (
	"fmt"
	"github.com/Jooeeee/NGINX_CONFIG_PROCESSOR/processor"
)

func main(){
	fmt.Println("This is Test!")
	srv,err:=processor.Scanner("./server.conf")
	fmt.Println(srv,err)
	err=processor.Dumper(srv,"./new.conf")
}