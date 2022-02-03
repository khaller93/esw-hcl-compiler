package io

import (
    "fmt"
    "github.com/hashicorp/hcl/hcl/ast"
    "github.com/hashicorp/hcl/hcl/parser"
    "gitlab.isis.tuwien.ac.at/semsys/ess/esw-hcl-compiler/transformer"
    "io/ioutil"
)

type ReadHCLError struct {
    filepath string
    message  string
}

func (e ReadHCLError) Error() string {
    return "Error in file '" + e.filepath + "': " + e.message
}

// reads HCL file
func ReadHCLFile(filepath string) (*ast.File, error) {
    inputContent, err := ioutil.ReadFile(filepath)
    if err == nil {
        return ReadHCLBytes(inputContent)
    }
    return nil, err
}

// reads HCL string
func ReadHCLString(data string) (*ast.File, error) {
    return ReadHCLBytes([]byte(data))
}

// reads HCL bytes
func ReadHCLBytes(data []byte) (*ast.File, error) {
    astData, err := parser.Parse(data)
    if err == nil {
        return astData, nil
    }
    return nil, err
}

// reads HCL file for prefix mapping
func ReadPrefixHCL(filepath string) (map[string]string, error) {
    astData, err := ReadHCLFile(filepath)
    if err == nil {
        node, err := transformer.Transform(astData)
        if err == nil {
            switch t := node.(type) {
            case map[string]interface{}:
                pMapNode := node.(map[string]interface{})
                prefixMap := make(map[string]string)
                for k, v := range pMapNode {
                    switch t := v.(type) {
                    case string:
                        prefixMap[k] = v.(string)
                    default:
                        return nil, ReadHCLError{filepath: filepath, message: fmt.Sprintf("The value of a prefix must be a string, but %T was given for prefix '%s'.", t, k)}
                    }
                }
                return prefixMap, nil
            default:
                return nil, ReadHCLError{filepath: filepath, message: fmt.Sprintf("Expected to find a set of assignment, but was %T.", t)}
            }
        }
        return nil, err
    }
    return nil, err
}

// reads general HCL file
func ReadGeneralHCL(filepath string, prefixMap map[string]string) (interface{}, error) {
    astData, err := ReadHCLFile(filepath)
    if err == nil {
        return transformer.TransformWithPrefixMap(astData, prefixMap)
    }
    return nil, err
}
