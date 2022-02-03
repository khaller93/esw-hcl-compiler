package reactor

var functions map[string]interface{} = map[string]interface{}{
    "append": appendToList,
}

// append a list of values to an existing list. if no list has been encountered,
// the passed value stays unchanged.
func appendToList(values interface{}, argument interface{}) interface{} {
    switch values.(type) {
    case []interface{}:
        switch argument.(type) {
        case []interface{}:
            currentValues := values.([]interface{})
            newValues := make([]interface{}, 0)
            for i := range currentValues {
                newValues = append(newValues, currentValues[i])
            }
            currentArgument := argument.([]interface{})
            for i := range currentArgument {
                newValues = append(newValues, currentArgument[i])
            }
            return newValues
        default:
            return values
        }
    default:
        return values
    }
}
