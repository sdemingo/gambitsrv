package main

import (
	"fmt"
	"flag"
	"os"
)

func main() {
	portflag := flag.String("p", "22022", "Puerto de escucha")
	helpflag := flag.Bool("h", false, "Mostrar este mensaje de ayuda")
	
	flag.Parse()
	if *helpflag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Start gambitsrv on port "+ (*portflag)+" ...")
	InitServer(*portflag)
}
