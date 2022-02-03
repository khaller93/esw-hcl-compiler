package reactor

// merges the two given objects, where the object b is taking precedence
// over the first one.
func mergeObject(a, b map[string]interface{}, path string) map[string]interface{} {
    var object = make(map[string]interface{})

    /* fill in 'a' object values */
    for keyA, valueA := range a {
        object[keyA] = valueA
    }

    /* extract b fields */
    for keyB, valueB := range b {
        valueA, found := object[keyB]
        if found {
            object[keyB] = merge(valueA, valueB, path+"/"+keyB)
        } else {
            object[keyB] = valueB
        }
    }

    return object
}

// merges the given object A with object B. if object
// a and b are not both a map, otherwise always b is
// returned. i.e. two objects are merged, but other
// combinations are always overwritten by b.
func merge(a, b interface{}, path string) interface{} {
    switch a.(type) {
    case map[string]interface{}:
        switch b.(type) {
        case map[string]interface{}:
            return mergeObject(a.(map[string]interface{}), b.(map[string]interface{}), path)
        default:
            return b
        }
    default:
        return b
    }
}
