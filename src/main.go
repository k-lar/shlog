package main

import (
    "fmt"
    "os"
    "io"
    "log"
    "bufio"
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

func writeToFile(content string) {
    f, err := os.OpenFile("shows.txt",
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
    f, _ := os.Open("shows.txt")

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

    f, _ = os.Create("shows.txt")
    tmpfile, _ := os.Open(f_tmp.Name())

    scanner := bufio.NewScanner(tmpfile)
    for scanner.Scan() {
        line := scanner.Text()
        if line != show {
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
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")

    entry := showName + ";S" + showSeason + "E" + showEpisode
    if showTime != "" {
        entry = entry + "T" + showTime
    }

    writeToFile(entry)

    pterm.Println(pterm.LightGreen("Entry added:"))
    pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
		{"Name", "Season", "Episode", "Time"},
		{showName, showSeason, showEpisode, showTime},
	}).Render()
}

func fuzzySearchShows() string {
    fmt.Print("\033[1A\033[K")
    fmt.Print("\033[1A\033[K")
    area, _ := pterm.DefaultArea.Start() // Start the Area printer.

    shows, err := readFile("shows.txt")
    if err != nil {
        log.Fatalf("readFile: %s", err)
    }

    selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(shows).Show()
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
        fuzzySearchShows()
	case "[3] - Edit progress":
        fuzzySearchShows()
	case "[4] - Delete show":
        showToRemove := fuzzySearchShows()
        if (confirmRemoval(showToRemove)) {
            removeShow(showToRemove)
            fmt.Print("\033[1A\033[K")
        }
        fmt.Print("\033[1A\033[K")
        menu()

	case "[5] - Quit":
        os.Exit(0)
	}
    area.Stop()
}

func main() {
    menu()
}

