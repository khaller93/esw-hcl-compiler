package reactor

import (
    "fmt"
    "github.com/fatih/color"
    "os"
    "sort"
)

type ReactorNode struct {
    file               ReactorFile
    extends            []string // list of identifiers from which to inherit
    extendsUnProcessed []string // list of identifiers not processed yet
    precedenceOver     []string // list of identifiers over which this node takes precedence
}

func getArray(data interface{}, key string) ([]interface{}, bool) {
    switch data.(type) {
    case map[string]interface{}:
        arrData, found := (data.(map[string]interface{}))[key]
        if found {
            switch arrData.(type) {
            case []interface{}:
                return arrData.([]interface{}), true
            }
        }
    }
    return nil, false
}

func toStringArray(array []interface{}) ([]string, error) {
    var stringArray = make([]string, len(array))
    for i := range array {
        var elem = array[i]
        switch elem.(type) {
        case string:
            stringArray[i] = elem.(string)
        default:
            return nil, ParsingError{"Array contains not only strings as it would be expected."}
        }
        stringArray[i] = array[i].(string)
    }
    return stringArray, nil
}

func pack(files []ReactorFile) []ReactorNode {
    var nodes = make([]ReactorNode, len(files))
    for f := range files {
        file := files[f]

        /* check extends */
        var extends []string;
        extendNode, found := getArray(file.Data, "extends")
        if found {
            extendsArr, err := toStringArray(extendNode)
            if err == nil {
                extends = extendsArr
            } else {
                _, _ = color.New(color.FgYellow).Fprintf(os.Stdin, "    WARNING for 'extends' variable in \"%s\": %s\n", file.Filepath, err.Error())
            }
        }
        if extends == nil {
            extends = make([]string, 0)
        }
        delete(file.Data, "extends")

        /* check precedence over */
        var precedence []string;
        precedenceNode, found := getArray(file.Data, "precedence_over")
        if found {
            precedenceArr, err := toStringArray(precedenceNode)
            if err == nil {
                precedence = precedenceArr
            } else {
                _, _ = color.New(color.FgYellow).Fprintf(os.Stdin, "    WARNING for 'precedence_over' variable in \"%s\": %s\n", file.Filepath, err.Error())
            }
        }
        if precedence == nil {
            precedence = make([]string, 0)
        }
        delete(file.Data, "precedence_over")

        nodes[f] = ReactorNode{file: file, extends: extends, extendsUnProcessed: append([]string{}, extends...), precedenceOver: precedence}
    }
    return nodes
}

func getNextReactorNodeEntry(stack []ReactorNode) (*ReactorNode, []ReactorNode) {
    switch len(stack) {
    case 0:
        return nil, nil
    case 1:
        return &stack[0], nil
    default:
        sort.SliceStable(stack, func(i, j int) bool {
            return len(stack[i].extendsUnProcessed) < len(stack[j].extendsUnProcessed)
        })
        return &stack[0], stack[1:]
    }
}

func assembleClassNodes(stack []ReactorNode) map[string]interface{} {
    var data = make(map[string]interface{})
    for {
        elem, newStack := getNextReactorNodeEntry(stack)
        if elem == nil {
            if stack == nil {
                break
            } else {
                continue
            }
        }
        // merge
        var newData = elem.file.Data
        for i := range elem.extends {
            extendsId := elem.extends[i]
            extendData, found := data[extendsId]
            if found {
                fmt.Printf("-- %v \n", elem.file.Identifier)
                newData = merge(extendData, newData, "").(map[string]interface{})
                // manage precedence array
                preIds, found := extendData.(map[string]interface{})["precedence_over"]
                if found {
                    preIdsArr := preIds.([]string)
                    for n := range preIdsArr {
                        elem.precedenceOver = append(elem.precedenceOver, preIdsArr[n])
                    }
                }
            } else {
                _, _ = color.New(color.FgYellow).Fprintf(os.Stdin, "    WARNING for %s: Extension '%s' could not be found.\n", elem.file.Identifier, extendsId)
            }
            elem.precedenceOver = append(elem.precedenceOver, extendsId)
        }
        newData["precedence_over"] = elem.precedenceOver
        data[elem.file.Identifier] = newData
        // remove entry from unprocessed extends
        for i := range newStack {
            newStack[i].extendsUnProcessed = remove(newStack[i].extendsUnProcessed, elem.file.Identifier)
        }
        stack = newStack
    }
    return data
}

func assemblePatternNodes(stack []ReactorNode) map[string]interface{} {
    var data = make(map[string]interface{})
    for i := range stack {
        node := stack[i]
        data[node.file.Identifier] = node.file.Data
    }
    return data
}

func filterOfType(t string, files []ReactorFile) []ReactorFile {
    newFilesArr := make([]ReactorFile, 0)
    for i := range files {
        switch files[i].Type {
        case t:
            newFilesArr = append(newFilesArr, files[i])
            break
        default:
            continue
        }
    }
    return newFilesArr
}

func remove(list []string, elem string) []string {
    var newList = make([]string, 0)
    for i := range list {
        if list[i] != elem {
            newList = append(newList, list[i])
        }
    }
    return newList
}

func Assemble(dirPath string, prefixMap map[string]string) (interface{}, error) {
    reactorFiles, err := readAndTransformHCLFiles(dirPath, prefixMap)
    if err == nil {
        data := make(map[string]interface{})
        // instance of nodes
        classNodes := assembleClassNodes(pack(filterOfType("class", reactorFiles)))
        data["instance_of"] = classNodes
        // pattern nodes
        patternNodes := assemblePatternNodes(pack(filterOfType("iri", reactorFiles)))
        data["pattern"] = patternNodes
        return data, nil
    }
    return nil, err
}
