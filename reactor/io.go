package reactor

import (
    "fmt"
    "github.com/fatih/color"
    "gitlab.isis.tuwien.ac.at/semsys/ess/esw-hcl-compiler/io"
    "gitlab.isis.tuwien.ac.at/semsys/ess/esw-hcl-compiler/transformer"
    "os"
    "path/filepath"
    "strings"
)

type ReactorFile struct {
    Identifier string
    Type       string // iri or pattern
    Data       map[string]interface{}
    Filepath   string
}

type ParsingError struct {
    message string
}

func (e ParsingError) Error() string {
    return e.message
}

// gets the HCL files in the given path directory
func getHCLFiles(dirPath string) ([]string, error) {
    var files = make([]string, 0)
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".hcl") && !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

func cleanReactorType(data map[string]interface{}) map[string]interface{} {
    delete(data, "class")
    delete(data, "iri")
    return data
}

// parses the file on the given file path.
func parseFile(file string, prefixMap map[string]string) (*ReactorFile, error) {
    astData, err := io.ReadHCLFile(file)
    if err != nil {
        printParsingStatus(file, false)
        return nil, ParsingError{message: fmt.Sprintf("    Failed to read '%s'. %s\n", file, err.Error())}
    }
    node, err := transformer.TransformWithPrefixMap(astData, prefixMap)
    if err != nil {
        return nil, ParsingError{message: fmt.Sprintf("    Failed to parse '%s'. %s\n", file, err.Error())}
    }
    reactorNode := node.(map[string]interface{})

    /* check for IRI */
    result, found := reactorNode["class"]
    if found {
        switch result.(type) {
        case string:
            iri := result.(string)
            return &ReactorFile{Identifier: iri, Type: "class", Data: cleanReactorType(reactorNode), Filepath: file}, nil
        default:
            return nil, ParsingError{message: fmt.Sprintf("'iri' must be a string, but was not in '%s'.", file)}
        }
    }
    /* check for Pattern */
    result, found = reactorNode["iri"]
    if found {
        switch result.(type) {
        case string:
            pattern := result.(string)
            return &ReactorFile{Identifier: pattern, Type: "iri", Data: cleanReactorType(reactorNode), Filepath: file}, nil
        default:
            return nil, ParsingError{message: fmt.Sprintf("'pattern' must be a string, but was not in '%s'.", file)}
        }
    }

    return nil, ParsingError{message: fmt.Sprintf("'iri' or 'pattern' must be specified, but was not in '%s'.", file)}
}

// prints the status for a configuration file (OK or FAILED)
func printParsingStatus(filepath string, ok bool) {
    fmt.Printf("  + \"%s\" parsing: ", filepath)
    if ok {
        _, _ = color.New(color.FgGreen).Print("OK")
    } else {
        _, _ = color.New(color.FgRed).Print("FAILED")
    }
    fmt.Print("\n")
}

// read and transform HCL files
func readAndTransformHCLFiles(dirPath string, prefixMap map[string]string) ([]ReactorFile, error) {
    //
    files, err := getHCLFiles(dirPath)
    if err != nil {
        return nil, err
    }
    //
    var nodes []ReactorFile = make([]ReactorFile, len(files))
    var failed bool = false
    for f := range files {
        var file string = files[f]
        reactorFile, err := parseFile(file, prefixMap)
        if err == nil {
            nodes[f] = *reactorFile
            printParsingStatus(file, true)
        } else {
            failed = true
            printParsingStatus(file, false)
            _, _ = color.New(color.FgRed).Fprintf(os.Stderr, fmt.Sprintf("    %s\n", err.Error()))
        }
    }
    if failed {
        return nil, ParsingError{message: "Parsing the reactor configurations failed for at least one file."}
    }
    return nodes, nil
}
