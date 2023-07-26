package main

import (
	"fmt"
	"os"
)

/*
#include <sketchybar.h>
#include <stdint.h>

extern void MainCFunction(char* bootstrapName);
extern void handleEvent(char* env);
*/
import "C"

//export AGoFunction
func AGoFunction() {
	fmt.Println("AGoFunction()")
}

//export GoHandler
func GoHandler(
	name *C.char,
	sender *C.char,
	info *C.char,
	config_dir *C.char,
	button *C.char,
	modifier *C.char,
) {
	// Your Go event handling logic here
	nameStr := C.GoString(name)
	senderStr := C.GoString(sender)
	infoStr := C.GoString(info)
	configDirStr := C.GoString(config_dir)
	buttonStr := C.GoString(button)
	modifierStr := C.GoString(modifier)

	fmt.Printf("*** inside handler, called with values: "+
		"name=%s, sender=%s, info=%s, configdir=%s, button=%s, modifier=%s ***\n",
		nameStr, senderStr, infoStr, configDirStr, buttonStr, modifierStr,
	)

	fmt.Printf("got vars name: %s, sender: %s, info: %s\n ", nameStr, senderStr, infoStr)

	if senderStr == "front_app_switched" {
		// front_app item update
		command := fmt.Sprintf("--set %s label=\"%s\"", nameStr, infoStr)
		C.sketchybar(C.CString(command))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <bootstrap name> ", os.Args[0])
	}

	bootstrapName := os.Args[1] // first arg is 1st index, 0th is binary name
	fmt.Printf("Starting an event service with bootstrap name: %s",
		bootstrapName)
	C.MainCFunction(C.CString(bootstrapName))
}
