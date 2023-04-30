package main

import (
    "fmt"
    "os"
    "io"
    "log"
    "strings"
    "bufio"
    "runtime"
    "errors"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

func readFile(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func getShowsPath() string {
    var showsPath string;
    userdir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }

    if runtime.GOOS == "windows" {
        showsPath = userdir + "/shlog"
    } else {
        showsPath = userdir + "/.config/shlog"
    }

	error := os.MkdirAll(showsPath, os.ModePerm)
	if error != nil {
		log.Println(error)
	}

    showsPath = showsPath + "/shows.txt"
    return showsPath
}

func editEntry(showInfo []string) {
    fmt.Print("\033[1A\033[K")
    area, _ := pterm.DefaultArea.Start() // Start the Area printer.

    pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
        {"Name", "Season", "Episode", "Time"},
        {showInfo[0], trimLeftChar(showInfo[1]), trimLeftChar(showInfo[2]), trimLeftChar(showInfo[3])},
    }).Render()

    options := []string {
        "Show name",
        "Season",
        "Episode",
        "Time",
    }

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	printer.KeySelect = keys.Space
	printer.Checkmark = &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedOptions, _ := printer.Show()

    if (selectedOptions == nil) {
        fmt.Println("No edit option selected.")
        result, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show("Exit to menu?")
        if (result == true) {
            menu()
        } else {
            editEntry(showInfo)
            return
        }
    }

    entryName := showInfo[0]
    updatedInfo := showInfo
    updatedArr := [4]int{0, 0, 0, 0}

    // fmt.Println(showInfo)
    for i := 0; i < len(selectedOptions); i++ {

        switch selectedOptions[i] {
            case "Show name":
                result, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Update show name")
                updatedInfo[0] = result
                updatedArr[0] = 1

            case "Season":
                result, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Update season")
                updatedInfo[1] = result
                updatedArr[1] = 1

            case "Episode":
                result, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Update episode")
                updatedInfo[2] = result
                updatedArr[2] = 1

            case "Time":
                result, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Update time")
                updatedInfo[3] = result
                updatedArr[3] = 1
        }
    }

    // fmt.Println(updatedInfo)
    // fmt.Println(len(updatedInfo))

    var updatedEntry string
    if (updatedArr[0] == 1) {
        updatedEntry = updatedInfo[0]
    } else {
        updatedEntry = showInfo[0]
    }

    if (updatedArr[1] == 1) {
        updatedEntry = updatedEntry + ";S" + updatedInfo[1]
    } else {
        updatedEntry = updatedEntry + ";" + showInfo[1]
    }

    if (updatedArr[2] == 1) {
        updatedEntry = updatedEntry + ";E" + updatedInfo[2]
    } else {
        updatedEntry = updatedEntry + ";" + showInfo[2]
    }

    if (updatedArr[3] == 1) {
        updatedEntry = updatedEntry + ";T" + updatedInfo[3]
    } else {
        updatedEntry = updatedEntry + ";" + showInfo[3]
    }

    // fmt.Println(updatedEntry)
    // fmt.Println("Removing: ", entryName)
    removeShow(entryName)
    writeToFile(updatedEntry)
    // area.Clear()

    tmpEntry := strings.Split(updatedEntry, ";")
    fmt.Println()
    pterm.Println(pterm.LightGreen("Updated entry:"))
    printEntry(getShowInfo(getShowsPath(), tmpEntry[0]))
    area.Stop()
}

func writeToFile(content string) {
    f, err := os.OpenFile(getShowsPath(),
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
    }
    defer f.Close()
    if _, err := f.WriteString(content + "\n"); err != nil {
        log.Println(err)
    }
}



func confirmRemoval(show string) bool {
    fmt.Print("\033[1A\033[K")
    area, _ := pterm.DefaultArea.Start() // Start the Area printer.
    pterm.Warning.Printfln("Are sure you want to delete %s?", show)
	result, _ := pterm.DefaultInteractiveConfirm.Show()
    area.Clear()
    area.Stop()
    fmt.Print("\033[1A\033[K")
    return result
}

func removeShow(show string) {
    f, _ := os.Open(getShowsPath())

    // create and open a temporary file
    f_tmp, err := os.CreateTemp("", "tmpfile-*.txt")
    if err != nil {
        log.Fatal(err)
    }

    // Copy content from original to tmp
    _, err = io.Copy(f_tmp, f)
    if err != nil {
        log.Fatal(err)
    }

    f, _ = os.Create(getShowsPath())
    tmpfile, _ := os.Open(f_tmp.Name())

    scanner := bufio.NewScanner(tmpfile)
    for scanner.Scan() {
        line := scanner.Text()
        split := strings.Split(line, ";")
        if split[0] != show {
            if _, err := f.Write([]byte(line + "\n")); err != nil {
                fmt.Println(err)
            }
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    defer os.Remove(f_tmp.Name())
    defer f.Close()
}

func printEntry(showInfo []string) {
    if (len(showInfo) > 3) {
        pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
            {"Name", "Season", "Episode", "Time"},
            {showInfo[0], trimLeftChar(showInfo[1]), trimLeftChar(showInfo[2]), trimLeftChar(showInfo[3])},
        }).Render()
    } else {
        pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
            {"Name", "Season", "Episode"},
            {showInfo[0], trimLeftChar(showInfo[1]), trimLeftChar(showInfo[2])},
        }).Render()
    }
    fmt.Println()
}

func getShowInfo(path string, show string) []string {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    var split []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        split = strings.Split(line, ";")
        if (split[0] == show) {
            return split
        }
    }
    return split
}

func addShow() {
    showName, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter show name")
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    showSeason, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter season")
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    showEpisode, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter episode")
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    showTime, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter time (optional)")
    showTime = strings.TrimSpace(showTime)
    if (showTime == "") {
        showTime = "00:00"
    }
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    entry := showName + ";S" + showSeason + ";E" + showEpisode
    if showTime != "" {
        entry = entry + ";T" + showTime
    }

    writeToFile(entry)

    if (showTime != "") {
        pterm.Println(pterm.LightGreen("Entry added:"))
        pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
            {"Name", "Season", "Episode", "Time"},
            {showName, showSeason, showEpisode, showTime},
        }).Render()
    } else {
        pterm.Println(pterm.LightGreen("Entry added:"))
        pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
            {"Name", "Season", "Episode"},
            {showName, showSeason, showEpisode},
        }).Render()
    }
    fmt.Println()
}

func prettyReadFile(path string) ([]string, []string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, nil, err
    }
    defer file.Close()

    var lines []string
    var split []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        split = strings.Split(line, ";")
        lines = append(lines, split[0])
    }
    return lines, split, scanner.Err()
}

func fuzzySearchShows() string {
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    area, _ := pterm.DefaultArea.Start() // Start the Area printer.

    shows, _, err := prettyReadFile(getShowsPath())
    if err != nil {
        log.Fatalf("readFile: %s", err)
    }

    selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(shows).Show("Select a show:")
    pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))

    area.Update(selectedOption)
    area.Clear()
    area.Stop()
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")
    return selectedOption
}

func menu() {
    fmt.Print("\033[1A\033[K")
    area, _ := pterm.DefaultArea.Start() // Start the Area printer.
    menu_options := []string {
        "[1] - Add show",
        "[2] - View progress",
        "[3] - Edit progress",
        "[4] - Delete show",
        "[5] - Quit",
    }

    selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(menu_options).Show()

    area.Update(selectedOption)
    area.Clear()

    switch selectedOption {
	case "[1] - Add show":
        fmt.Print("\033[1A\033[K")
        addShow()

	case "[2] - View progress":
        if _, err := os.Stat(getShowsPath()); err == nil {
            // path/to/whatever exists
            printEntry(getShowInfo(getShowsPath(), fuzzySearchShows()))
        } else if errors.Is(err, os.ErrNotExist) {
            // path/to/whatever does *not* exist
            fmt.Println("No shows can be found.")
            os.Exit(0)
        } else {
            // Schrodinger: file may or may not exist. See err for details.
            // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
            fmt.Println("ERROR: Something is wrong with shows.txt file.")
            os.Exit(1)
        }

	case "[3] - Edit progress":
        if _, err := os.Stat(getShowsPath()); err == nil {
            // path/to/whatever exists
            editEntry(getShowInfo(getShowsPath(), fuzzySearchShows()))
        } else if errors.Is(err, os.ErrNotExist) {
            // path/to/whatever does *not* exist
            fmt.Println("No shows to edit.")
            os.Exit(0)
        } else {
            // Schrodinger: file may or may not exist. See err for details.
            // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
            fmt.Println("ERROR: Something is wrong with shows.txt file.")
            os.Exit(1)
        }

	case "[4] - Delete show":
        if _, err := os.Stat(getShowsPath()); err == nil {
            // path/to/whatever exists
            showToRemove := fuzzySearchShows()
            fmt.Println()
            if (confirmRemoval(showToRemove)) {
                removeShow(showToRemove)
                fmt.Print("\033[1A\033[K")
            }
            fmt.Print("\033[1A\033[K")
            menu()
        } else if errors.Is(err, os.ErrNotExist) {
            // path/to/whatever does *not* exist
            fmt.Println("No shows to delete.")
            os.Exit(0)
        } else {
            // Schrodinger: file may or may not exist. See err for details.
            // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
            fmt.Println("ERROR: Something is wrong with shows.txt file.")
            os.Exit(1)
        }

	case "[5] - Quit":
        os.Exit(0)
	}
    area.Stop()
}

func main() {
    menu()
}

