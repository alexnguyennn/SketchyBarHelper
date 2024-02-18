package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/samber/lo"
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

type window struct {
	id    string
	app   string
	title string
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
	/*	configDirStr := C.GoString(config_dir)
		buttonStr := C.GoString(button)
		modifierStr := C.GoString(modifier)
	*/
	/*	fmt.Printf(
			"*** inside handler, called with values: "+
				"name=%s, sender=%s, info=%s, configdir=%s, button=%s, modifier=%s ***\n",
			nameStr, senderStr, infoStr, configDirStr, buttonStr, modifierStr,
		)
	*/
	currentDisplay := ""
	focusedWorkspace := ""
	switch senderStr {
	case "display_change":
		//fmt.Printf("display change event activated; got: %s", infoStr)
		currentDisplay = infoStr
	case "aerospace_workspace_change":
		focusedWorkspace = os.Getenv("FOCUSED_WORKSPACE")
	default:
	}

	if currentDisplay == "" {
		currentDisplay = nameStr[strings.LastIndex(nameStr, ".")+1:]
	}
	//fmt.Printf("picked up currentDisplay value: %s\n", currentDisplay)

	// Run the 'yabai' command and capture its output
	if focusedWorkspace == "" {
		yabaiPipeline := script.Exec(
			fmt.Sprintf(
				`aerospace list-workspaces --monitor %s --visible`,
				currentDisplay,
			),
		)

		yabaiOutput, err := yabaiPipeline.String()
		if err != nil {
			fmt.Printf("ran yabai command 1 and ran into this error: %s\n", err.Error())
		}
		focusedWorkspace = strings.TrimSpace(yabaiOutput)
	}

	//fmt.Printf("got workspace: %s\n", focusedWorkspace)
	if len(focusedWorkspace) == 0 {
		//panic("never got space for display")
		fmt.Println("never got space for display")
		// TODO; happens with break; improve by focusing last focused?
	}

	windowListCmd := fmt.Sprintf(
		`aerospace list-windows --workspace "%s"`,
		focusedWorkspace,
	)
	windowStrs, err := script.Exec(
		windowListCmd,
	).Slice()
	if err != nil {
		fmt.Printf("ran %s: %s\n", windowListCmd, err.Error())
		return
	}

	windows := lo.Map(
		windowStrs, func(winStr string, i int) window {
			parts := strings.SplitN(winStr, "|", 3)
			return window{
				id:    strings.TrimSpace(parts[0]),
				app:   strings.TrimSpace(parts[1]),
				title: strings.TrimSpace(parts[2]),
			}
		},
	)

	//windows := strings.Split(windowTitles, "\n")

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
	/*	displayInfo := script.Exec(
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
	*/
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

	visibleWindowLimit := lo.Ternary(
		currentDisplay == "4", 4, 6,
	)

	numWindows := len(windows)
	var titleWidth int
	if numWindows < 3 || numWindows == 0 {
		titleWidth = 200
	} else {
		/*if strings.Contains(displayType, "portrait") {
			titleWidth = 60
		} else if numWindows > 0 {*/
		//titleWidth = usableWidth / numWindows
		titleWidth = 100
		//}
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

	focusedID, err := script.
		Exec(`aerospace list-windows --focused`).
		Exec(`cut -d'|' -f1`).
		String()
	if err != nil {
		fmt.Printf("ran %s: %s\n", `aerospace list-windows --focused`, err.Error())
		return
	}

	focusedID = strings.TrimSpace(focusedID)

	numTitleLabels := 8
	for i := 0; i < numTitleLabels; i++ {
		if i >= numWindows || i >= visibleWindowLimit {
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

		// TODO: hasfocused by separate call for focused window, match

		// TODO; make this a config file read from launch dir
		// TODO: read from colors.sh maybe too
		borderColor := `0x00000000`
		hasFocus := windows[i].id == focusedID
		if hasFocus {
			//fmt.Printf("got focus match: %s | %s\n", focusedID, windows[i].title)
			borderColor = `0xffffffff`
		} else {
			//fmt.Printf("failed focus match: %s vs %s | win: %s\n", focusedID, windows[i].id, windows[i].title)
		}

		windowTitle := windows[i].title

		sketchybarArgsBuilder.WriteString(
			fmt.Sprintf(
				`--set title.%s.%d label="%s" label.width="%d" background.drawing="on" background.border_color="%s" `, // trailing space per command is important
				//`--set title.%s.%d label=%s label.width=%d background.color=%s click_script="${CONFIG_DIR}/plugins/focus.sh %s" icon=%s icon.font=%s `,
				currentDisplay,
				i,
				windowTitle,
				//fmt.Sprintf(`"%s| %s"`, strings.Trim(appIconStr, `"`), strings.Trim(windowTitle, `"`)),
				titleWidth,
				borderColor,
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
