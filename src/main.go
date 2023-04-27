package main

import (
    // "fmt"
    "os"
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

func fuzzySearchShows() string {
    shows, err := readFile("shows.txt")
    if err != nil {
        log.Fatalf("readLines: %s", err)
    }

    selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(shows).Show()
    pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))
    return selectedOption
}


func menu() {
    menu_options := []string {
        "[1] - Add show",
        "[2] - View progress",
        "[3] - Edit progress",
        "[4] - Delete show",
        "[5] - Quit",
    }

    selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(menu_options).Show()

    switch selectedOption {
	case "[1] - Add show":
        fuzzySearchShows()
	case "[2] - Edit progress":
        fuzzySearchShows()
	case "[3] - Delete show":
        fuzzySearchShows()
	case "[4] - Quit":
        os.Exit(0)
	}
}

func main() {
    menu()
}

