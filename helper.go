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
	// TODO: improve performance.. - only make changes necessary -> can query yabai for focused win.. maybe just open an issue and ask for help
	nameStr := C.GoString(name)
	senderStr := C.GoString(sender)
	infoStr := C.GoString(info)
	/*configDirStr := C.GoString(config_dir)
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
			),
		).JQ(`.[] | select(."is-visible" == true) | .index`)

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
		),
	).
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
	// TODO: account for bar width?
	// TODO: somehow force not rendering past right container
	// if i have n <= 4; have them wider; n <=6 shorter; n >6 is icon only?
	// WIP: scale widths with ranges
	// TODO: get bounding rects of spaces_bracket and right_section
	// TODO; get total x coord with yabai -m query --displays
	displayInfo := script.Exec(
		fmt.Sprintf(
			"yabai -m query --displays --display %s",
			currentDisplay,
		),
	)
	//).JQ(".frame.x").String()

	displayType, err := displayInfo.JQ(`if .frame.w > .frame.h then "landscape" else "portrait" end`).String()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	displayType = strings.TrimSpace(displayType)

	/*	displayXPosition, err := displayInfo.JQ(".frame.w").String()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	*/
	/*	rightXPosition, err := script.Exec("sketchybar --query right_section").
		JQ(fmt.Sprintf(`.bounding_rects."display-%s".origin[0]`, currentDisplay)).String()*/
	/*	rightWidth, err := script.Exec("sketchybar --query right_section").
			JQ(fmt.Sprintf(`.bounding_rects."display-%s".size[0]`, currentDisplay)).String()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		spacesWidth, err := script.Exec("sketchybar --query spaces_bracket").
			JQ(fmt.Sprintf(`.bounding_rects."display-%s".size[0]`, currentDisplay)).String()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}*/
	/*displayPos, err := strconv.Atoi(strings.TrimSpace(displayXPosition))
	if err != nil {
		fmt.Printf("%s", err.Error())
	}*/
	/*	displayWidth, err := strconv.Atoi(strings.TrimSpace(displayXPosition))
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		rightWI, err := strconv.Atoi(strings.TrimSpace(rightWidth))
		if err != nil {
			fmt.Printf("%s", err.Error())
		}

		spacesWI, err := strconv.Atoi(strings.TrimSpace(spacesWidth))
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	*/
	//usableWidth := rightWI - displayWidth - (10 * spacesWI)
	//usableWidth := displayWidth - (rightWI + spacesWI)
	numWindows := len(yabaiWindows)
	var titleWidth int
	if numWindows < 3 || numWindows == 0 {
		titleWidth = 200
	} else {
		if strings.Contains(displayType, "portrait") {
			titleWidth = 60
		} else if numWindows > 0 {
			//titleWidth = usableWidth / numWindows
			titleWidth = 100
		}
	}
	/*if strings.Contains(displayType, "portrait") {
		titleWidth = 50
	} else {
		if numWindows < 3 || numWindows == 0 {
			titleWidth = 200
		} else if numWindows > 0 {
			titleWidth = usableWidth / numWindows
		}
	}
	*/
	//titleWidth := 200
	// if numWindows < 4 {
	// 	titleWidth = 150
	// } else if numWindows < 8 {
	// 	titleWidth = 100
	// } else {
	// 	titleWidth = 50
	// }

	// TODO: refine font decl
	const iconFont = "sketchybar-app-font:Regular:16.0"

	// Append strings to the builder
	// builder.WriteString("Hello, ")
	// builder.WriteString("world!")
	// builder.WriteByte('\n') // Append a single byte (newline character)

	// Convert the builder to a string
	//result := builder.String()

	//for _, windowStr := range yabaiWindows {
	numTitleLabels := 8
	for i := 0; i < numTitleLabels; i++ {
		if i >= numWindows {
			// no matching window; set label to empty
			sketchybarArgsBuilder.WriteString(
				fmt.Sprintf(
					//`--set title.%s.%d label=%s label.width=%d background.color=%s `,
					`--set title.%s.%d label="%s" label.width=0 background.border_color=0x0 background.drawing=off icon="" `,
					currentDisplay,
					i,
					"",
				),
			)
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

		appTitleOutput, err := script.Echo(windowStr).JQ(`.app`).String()
		if err != nil {
			panic(fmt.Errorf("error getting app for %s: %s\n", windowStr, err))
		}

		// TODO: parse config_dir variable properly
		cleanedUpAppTitle := strings.Trim(strings.TrimSpace(appTitleOutput), `"`)
		pipe := script.Echo(cleanedUpAppTitle).
			Exec(`xargs -I{} /Users/alex/.config/sketchybar/plugins/icon_map.sh {}`)
		//Exec(fmt.Sprintf(`xargs -I{} %s/plugins/icon_map.sh {}`, configDirStr))
		appIconStr, err := pipe.String()
		if err != nil {
			pipe.SetError(nil)
			panic(
				fmt.Errorf(
					"error mapping app to icon for (%s):\n pipe string result - %s\n error - %s\n",
					cleanedUpAppTitle, appIconStr, err,
				),
			)
		}

		hasFocus, err := strconv.ParseBool(strings.TrimSpace(hasFocusStr))
		if err != nil {
			panic(fmt.Errorf("parsing bool from hasFocusStr - %s: %s\n", hasFocusStr, err))
		}

		// TODO; make this a config file read from launch dir
		// TODO: read from colors.sh maybe too
		borderColor := `0x00000000`
		if hasFocus {
			borderColor = `0xffffffff`
		}

		windowId = strings.TrimSpace(windowId)
		windowTitle = strings.TrimSpace(windowTitle)

		sketchybarArgsBuilder.WriteString(
			fmt.Sprintf(
				`--set title.%s.%d label=%s label.width="%d" background.drawing="on" background.border_color="%s" click_script="${CONFIG_DIR}/plugins/focus.sh %s" icon="%s" icon.font="%s" `,
				//`--set title.%s.%d label=%s label.width=%d background.color=%s click_script="${CONFIG_DIR}/plugins/focus.sh %s" icon=%s icon.font=%s `,
				currentDisplay,
				i,
				windowTitle,
				//fmt.Sprintf(`"%s| %s"`, strings.Trim(appIconStr, `"`), strings.Trim(windowTitle, `"`)),
				titleWidth,
				borderColor,
				windowId,
				strings.TrimSpace(appIconStr),
				iconFont,
			),
		)
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
	/*	rand.Seed(time.Now().UnixNano())
		randomNumber := rand.Intn(100)
		f, err := os.Create(fmt.Sprintf("/Users/alex/cpu-%d.pprof", randomNumber))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	*/
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <bootstrap name> ", os.Args[0])
	}

	bootstrapName := os.Args[1] // first arg is 1st index, 0th is binary name
	fmt.Printf(
		"Starting an event service with bootstrap name: %s\n",
		bootstrapName,
	)
	C.MainCFunction(C.CString(bootstrapName))
}
