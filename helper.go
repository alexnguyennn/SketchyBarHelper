package main

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
	"strconv"
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
	// TODO: somehow make right side be in front even if labels spill over - happens by default
	// TODO: improve performance.. - only make changes necessary -> can query yabai for focused win.. maybe just open an issue and ask for help
	// TODO: fixup icons
	// TODO: add icons in windows
	// TODO: dial in titles - icon  App / Title?
	// TODO; dial in title theming
	// TODO: scale widths with ranges
	nameStr := C.GoString(name)
	senderStr := C.GoString(sender)
	infoStr := C.GoString(info)
	/*	configDirStr := C.GoString(config_dir)
		buttonStr := C.GoString(button)
		modifierStr := C.GoString(modifier)
	*/
	/*	fmt.Printf("*** inside handler, called with values: "+
		"name=%s, sender=%s, info=%s, configdir=%s, button=%s, modifier=%s ***\n",
		nameStr, senderStr, infoStr, configDirStr, buttonStr, modifierStr,
	)*/

	currentDisplay := ""
	spacesJsonStr := ""
	switch senderStr {
	case "display_change":
		//fmt.Printf("display change event activated; got: %s", infoStr)
		currentDisplay = infoStr
	case "space_change":
		//fmt.Printf("space change event activated; got: %s", infoStr)
		spacesJsonStr = infoStr
	default:
	}

	if currentDisplay == "" {
		currentDisplay = nameStr[strings.LastIndex(nameStr, ".")+1:]
	}
	//fmt.Printf("picked up currentDisplay value: %s\n", currentDisplay)

	// Run the 'yabai' command and capture its output
	activeSpaceOnCurrentDisplay := ""
	if spacesJsonStr == "" {
		yabaiPipeline := script.Exec(
			fmt.Sprintf(
				`yabai -m query --spaces --display "%s"`,
				currentDisplay,
			)).JQ(`.[] | select(."is-visible" == true) | .index`)

		yabaiOutput, err := yabaiPipeline.String()
		if err != nil {
			fmt.Printf("ran yabai command 1 and ran into this error: %s\n", err.Error())
		}
		activeSpaceOnCurrentDisplay = yabaiOutput
	} else {
		spaceLookupOutput, err := script.Echo(spacesJsonStr).
			JQ(fmt.Sprintf(`."display-%s"`, currentDisplay)).
			String()
		if err != nil {
			fmt.Printf("error while looking up space for %s in %s", currentDisplay, spacesJsonStr)
		}

		activeSpaceOnCurrentDisplay = spaceLookupOutput
		//fmt.Printf("looked up spacejsonstr instead and got %s\n", activeSpaceOnCurrentDisplay)
	}

	//fmt.Printf("got current space on display from yabai: %s\n", yabaiOutput)
	if len(activeSpaceOnCurrentDisplay) == 0 {
		//panic("never got space for display")
		fmt.Println("never got space for display")
		// TODO; happens with break; improve by focusing last focused?
	}

	yabaiWindows, err := script.Exec(
		fmt.Sprintf(
			`yabai -m query --windows --space "%s"`,
			strings.TrimSpace(activeSpaceOnCurrentDisplay),
		)).
		JQ(`sort_by(.frame.x, .frame.y, ."stack-index") | .[]`).
		Slice()

	if err != nil {
		fmt.Printf("ran yabai command 2 and ran into this error: %s\n", err.Error())
	}

	//fmt.Print("got windows from yabai: %s\n", yabaiWindows)

	// Create a new strings.Builder
	//var sketchybarRemoveArgsBuilder strings.Builder
	var sketchybarArgsBuilder strings.Builder

	// TODO: lookup how to substitute new method without having to manually move the closeing bracket
	/*	sketchybarArgsBuilder.WriteString(
		fmt.Sprintf(
			`--remove /title\.%s\./ `,
			currentDisplay,
		))*/
	//sketchybarRemoveArgsBuilder.WriteByte('\n') // Append a single byte (newline character)
	//C.sketchybar(C.CString(sketchybarRemoveArgsBuilder.String()))

	// TODO: be smarter about this calculation
	// if i have n <= 4; have them wider; n <=6 shorter; n >6 is icon only?
	titleWidth := 200

	// Append strings to the builder
	// builder.WriteString("Hello, ")
	// builder.WriteString("world!")
	// builder.WriteByte('\n') // Append a single byte (newline character)

	// Convert the builder to a string
	//result := builder.String()

	numWindows := len(yabaiWindows)
	//for _, windowStr := range yabaiWindows {
	numTitleLabels := 8
	for i := 0; i < numTitleLabels; i++ {
		if i >= numWindows {
			// no matching window; set label to empty
			sketchybarArgsBuilder.WriteString(
				fmt.Sprintf(
					//`--set title.%s.%d label=%s label.width=%d background.color=%s `,
					`--set title.%s.%d label="%s" label.width=0 background.color=0x0 `,
					currentDisplay,
					i,
					"",
				))
			continue
		}

		windowStr := yabaiWindows[i]
		//fmt.Print("got window from yabai: %s\n", windowStr)
		windowId, err := script.Echo(windowStr).JQ(`.id`).String()
		if err != nil {
			panic(fmt.Errorf("error getting window id for %s: %s\n", windowStr, err))
		}

		hasFocusStr, err := script.Echo(windowStr).JQ(`."has-focus"`).String()
		if err != nil {
			panic(fmt.Errorf("error getting has focus for %s: %s\n", windowStr, err))
		}

		windowTitle, err := script.Echo(windowStr).JQ(`.title`).String()
		if err != nil {
			panic(fmt.Errorf("error getting window title for %s: %s\n", windowStr, err))
		}

		hasFocus, err := strconv.ParseBool(strings.TrimSpace(hasFocusStr))
		if err != nil {
			panic(fmt.Errorf("parsing bool from hasFocusStr - %s: %s\n", hasFocusStr, err))
		}

		backgroundColour := `0xff06decd`
		if hasFocus {
			backgroundColour = `0xfff0a104`
		}

		windowId = strings.TrimSpace(windowId)
		windowTitle = strings.TrimSpace(windowTitle)

		sketchybarArgsBuilder.WriteString(
			fmt.Sprintf(
				`--set title.%s.%d label=%s label.width=%d background.color=%s click_script="${CONFIG_DIR}/plugins/focus.sh %s" `,
				currentDisplay,
				i,
				windowTitle,
				titleWidth,
				backgroundColour,
				windowId,
			))
		//sketchybarArgsBuilder.WriteByte('\n') // Append a single byte (newline character)

	}

	// set string
	sketchybarCommand := sketchybarArgsBuilder.String()
	//fmt.Printf("\n\nabout to run this sketchybarCommand: %s\n\n", sketchybarCommand)
	C.sketchybar(C.CString(sketchybarCommand))
	// TODO: free command string after done?

	/*if senderStr == "front_app_switched" {
		// front_app item update
		command := fmt.Sprintf("--set %s label=\"%s\"", nameStr, infoStr)
		C.sketchybar(C.CString(command))
		// TODO: free command string after done
	}*/
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
