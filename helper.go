package main

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
	"strings"
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

	currentDisplay := ""
	switch senderStr {
	case "display_change":
		currentDisplay = infoStr
	default:
	}

	if currentDisplay == "" {
		currentDisplay = nameStr[strings.LastIndex(nameStr, ".")+1:]
	}
	fmt.Printf("picked up currentDisplay value: %s\n", currentDisplay)

	// Run the 'yabai' command and capture its output
	yabaiPipeline := script.Exec(
		fmt.Sprintf(
			`yabai -m query --spaces --display "%s"`,
			currentDisplay,
		)).JQ(`.[] | select(."is-visible" == true) | .index`)

	yabaiOutput, err := yabaiPipeline.String()
	if err != nil {
		fmt.Printf("ran yabai command 1 and ran into this error: %s\n", err.Error())
	}

	fmt.Printf("got current space on display from yabai: %s\n", yabaiOutput)

	yabaiWindows, err := script.Exec(
		fmt.Sprintf(
			`yabai -m query --windows --space "%s"`,
			strings.TrimSpace(yabaiOutput),
		)).
		JQ(`sort_by(.frame.x, .frame.y, ."stack-index") | .[]`).
		Slice()

	if err != nil {
		fmt.Printf("ran yabai command 2 and ran into this error: %s\n", err.Error())
	}

	fmt.Printf("got windows from yabai: %s\n", yabaiWindows)

	// convert windows to a slice of strings
	// TODO: create structure/find some way to iterate (can't use go structs here) (lean on go scripting?)
	// iterate over each component
	// generate args
	// create sketchybar command
	// call sketchybar

	if senderStr == "front_app_switched" {
		// front_app item update
		command := fmt.Sprintf("--set %s label=\"%s\"", nameStr, infoStr)
		C.sketchybar(C.CString(command))
		// TODO: free command string after done
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
