package main

import (
    "flag"
    "fmt"
    "github.com/fatih/color"
    "gitlab.isis.tuwien.ac.at/semsys/ess/esw-hcl-compiler/io"
    "gitlab.isis.tuwien.ac.at/semsys/ess/esw-hcl-compiler/reactor"
    "os"
    "path/filepath"
    "strings"
)

const ApplicationVersion string = "1.0.0"

// usage message for this application.
func usage() {
    fmt.Printf("Usage of %s (version %s):\n    %s <input-dir> <output-dir>\n", os.Args[0], ApplicationVersion,
        os.Args[0])
    flag.PrintDefaults()
}

// prints the status for a configuration file (OK or FAILED)
func printStatus(filepath string, ok bool) {
    fmt.Printf("\"%s\": ", filepath)
    if ok {
        _, _ = color.New(color.FgGreen).Print("OK")
    } else {
        _, _ = color.New(color.FgRed).Print("FAILED")
    }
    fmt.Print("\n")
}

// entry point
func main() {
    flag.Parse()
    args := flag.Args()

    if len(args) == 2 {
        var inputDirPath string = args[0]
        var outputDirPath string = args[1]

        // create output directory
        err := os.MkdirAll(outputDirPath, os.ModePerm)
        if err != nil {
            _, _ = color.New(color.FgRed).Fprint(os.Stderr, "Could not create output directory '%s'. %s", outputDirPath, err.Error())
            os.Exit(1)
        }

        // Prefix Map
        var prefixFilepath = filepath.Join(inputDirPath, "prefix.hcl")
        prefixMap, err := io.ReadPrefixHCL(filepath.Join(inputDirPath, "prefix.hcl"))
        if err == nil {
            err = io.WriteJSON(prefixMap, filepath.Join(outputDirPath, "prefix.json"))
            if err == nil {
                printStatus(prefixFilepath, true)
            } else {
                printStatus(prefixFilepath, false)
                _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
            }
        } else {
            printStatus(prefixFilepath, false)
            _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
        }
        if prefixMap == nil {
            prefixMap = make(map[string]string)
        }

        // General Configurations
        var failed bool = false
        generalFiles := []string{"general.hcl", "explorer.hcl"}
        for i := range generalFiles {
            var generalFilepath string = filepath.Join(inputDirPath, generalFiles[i])
            data, err := io.ReadGeneralHCL(generalFilepath, prefixMap)
            if err != nil {
                failed = true
                printStatus(generalFilepath, false)
                _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
                continue
            }
            err = io.WriteJSON(data, filepath.Join(outputDirPath, strings.Replace(generalFiles[i], ".hcl", ".json", 1)))
            if err != nil {
                failed = true
                printStatus(generalFilepath, false)
                _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
                continue
            }
            printStatus(generalFilepath, true)
        }

        // Reactor configurations
        var reactorDir = filepath.Join(inputDirPath, "reactor")
        fmt.Printf("\"%s\" ...\n", reactorDir)
        data, err := reactor.Assemble(reactorDir, prefixMap)
        if err == nil {
            err = io.WriteJSON(data, filepath.Join(outputDirPath, "reactor.json"))
            if err == nil {
                printStatus(reactorDir, true)
            } else {
                failed = true
                printStatus(reactorDir, false)
                _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
            }
        } else {
            failed = true
            printStatus(reactorDir, false)
            _, _ = color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
        }

        // exit
        if failed {
            os.Exit(1)
        } else {
            os.Exit(0)
        }
    } else {
        usage()
    }
}
