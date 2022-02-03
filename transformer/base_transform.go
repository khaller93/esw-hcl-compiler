package transformer

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/hashicorp/hcl/hcl/ast"
    "gitlab.isis.tuwien.ac.at/kg-utils/prefix.cc/api"
    "os"
    "regexp"
)

// transforms the HCL AST to better fit the ESW use case.
func Transform(node *ast.File) (interface{}, error) {
    return TransformWithPrefixMap(node, make(map[string]string))
}

// transforms the HCL AST to better fit the ESW use case.
func TransformWithPrefixMap(node *ast.File, prefixMap map[string]string) (interface{}, error) {
    return traverse(node, prefixMap)
}

func traverse(node interface{}, prefixMap map[string]string) (interface{}, error) {
    if node != nil {
        switch t := node.(type) {
        case *ast.File:
            return traverse(node.(*ast.File).Node, prefixMap)
        case *ast.ObjectList:
            return traverseObjectList(node.(*ast.ObjectList), prefixMap)
        case *ast.ObjectType:
            return traverseObjectList(node.(*ast.ObjectType).List, prefixMap)
        case *ast.ListType:
            return traverseList(node.(*ast.ListType), prefixMap)
        case *ast.LiteralType:
            return traverseLiteral((node).(*ast.LiteralType).Token.Value(), prefixMap)
        default:
            return nil, traverseError{node: node.(*ast.Node), message: fmt.Sprintf("Unknown type %T encountered.", t)}
        }
    }
    return nil, traverseError{node: nil}
}

func traverseList(list *ast.ListType, prefixMap map[string]string) (interface{}, error) {
    if list != nil {
        nodes := list.List
        transNodes := make([]interface{}, len(nodes))
        for i := 0; i < len(nodes); i++ {
            entry, err := traverse(nodes[i], prefixMap)
            if err != nil {
                return nil, err
            }
            transNodes[i] = entry
        }
        return transNodes, nil
    } else {
        return nil, traverseError{node: nil}
    }
}

var prefixRegex = regexp.MustCompile("^([0-9\\p{L}_]+[^:]):[0-9\\p{L}_]+$")
var prefixRepRegex = regexp.MustCompile("^([0-9\\p{L}_]+[^:]):")

func substitutePrefix(text string, prefixMap map[string]string) string {
    subMatches := prefixRegex.FindAllStringSubmatch(text, 1)
    if len(subMatches) > 0 && len(subMatches[0]) > 0 {
        var prefix = subMatches[0][1]
        namespace, found := prefixMap[prefix]
        if found {
            return prefixRepRegex.ReplaceAllString(text, namespace)
        } else {
            namespaces, err := api.GetOnTheFlyPrefixCC().GetNamespace(prefix)
            if err == nil && len(namespaces) > 0 {
                _, _ = color.New(color.FgYellow).Fprintf(os.Stdin, "    Uses namespace '%s' for prefix '%s'. Found using prefix.cc.\n", namespaces[0], prefix)
                return prefixRepRegex.ReplaceAllString(text, namespaces[0])
            } else {
                var errorMessage = namespaces[0]
                if err != nil {
                    errorMessage = "--" + err.Error() + " > " + namespaces[0]
                }
                _, _ = color.New(color.FgYellow).Fprintf(os.Stdin, "    No namespace could be found for prefix '%s'. %s\n", prefix, errorMessage)
            }
        }
    }
    return text
}

var keyRegex = regexp.MustCompile(`^"|"$`)

func prepareStringKey(key string, prefixMap map[string]string) string {
    return substitutePrefix(keyRegex.ReplaceAllString(key, ""), prefixMap)
}

func traverseObjectList(objectList *ast.ObjectList, prefixMap map[string]string) (interface{}, error) {
    if objectList != nil {
        var err error = nil
        items := objectList.Items
        data := make(map[string]interface{})
        for i := 0; i < len(items); i++ {
            keys := items[i].Keys
            var currentData = &data
            if len(keys) > 1 {
                for k := 0; k < len(keys)-1; k++ {
                    var key = prepareStringKey((*(keys[k])).Token.Text, prefixMap)
                    keyVal, found := (*currentData)[key]
                    if found {
                        switch keyVal.(type) {
                        case map[string]interface{}:
                            valueMap := keyVal.(map[string]interface{})
                            currentData = &valueMap
                            continue
                        default:
                            break
                        }
                    }
                    newMap := make(map[string]interface{})
                    (*currentData)[key] = newMap
                    currentData = &newMap
                }
            }
            var key = prepareStringKey((*(keys[len(keys)-1])).Token.Text, prefixMap)
            val, err := traverse(items[i].Val, prefixMap)
            if err != nil {
                return nil, err
            }
            (*currentData)[key] = val
        }
        return data, err
    } else {
        return nil, traverseError{node: nil}
    }
}

func traverseLiteral(literal interface{}, prefixMap map[string]string) (interface{}, error) {
    switch literal.(type) {
    case string:
        return substitutePrefix(literal.(string), prefixMap), nil
    default:
        return literal, nil
    }
}

// error when fetching namespaces for a prefix
type traverseError struct {
    node    *ast.Node
    message string
}

func (e traverseError) Error() string {
    node := e.node
    if node != nil {
        filename := (*node).Pos().Filename
        line := (*node).Pos().Line
        return "Failed to traverse at file '" + filename + "', line=" + string(line) + "." + e.message
    } else {
        return "Failed to traverse the node." + e.message
    }
}
