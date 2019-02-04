package hclq

import (
	"fmt"
	"github.com/mattolenik/hclq/query"
	"strconv"
)

// Get performs a query and returns a deserialized value
func (doc *HclDocument) Get(q string) (interface{}, error) {
	qry, _ := query.Parse(q)
	resultPairs, err := doc.Query(qry)
	if err != nil {
		return nil, err
	}
	if len(resultPairs) == 1 {
		return resultPairs[0].Value, nil
	}
	results := []interface{}{}
	for _, pair := range resultPairs {
		results = append(results, pair.Value)
	}
	return results, nil
}

// GetAsInt performs Get but converts the result to a string
func (doc *HclDocument) GetAsInt(q string) (int, error) {
	result, err := doc.Get(q)
	if err != nil {
		return 0, err
	}
	num, ok := result.(int)
	if ok {
		return num, nil
	}
	str, ok := result.(string)
	if ok {
		num, err := strconv.Atoi(str)
		if err == nil {
			return num, nil
		}
	}
	return 0, fmt.Errorf("could not find int at '%s' nor a string convertable to an int", q)
}

// GetAsString performs Get but converts the result to a string
func (doc *HclDocument) GetAsString(q string) (string, error) {
	result, err := doc.Get(q)
	if err != nil {
		return "", err
	}
	str, ok := result.(string)
	if ok {
		return str, nil
	}
	num, ok := result.(int)
	if ok {
		return strconv.Itoa(num), nil
	}
	return fmt.Sprintf("%v", result), nil
}

// GetAsList does the same as Get but converts it to a list for you (with type check)
func (doc *HclDocument) GetAsList(q string) ([]interface{}, error) {
	result, err := doc.Get(q)
	if err != nil {
		return nil, err
	}
	arr, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("query does not return a list, cannot be used with GetList")
	}
	return arr, nil
}

// GetAsStringList does the same as GetAsList but converts everything to a string for you.
func (doc *HclDocument) GetAsStringList(q string) ([]string, error) {
	list, err := doc.GetAsList(q)
	if err != nil {
		return nil, err
	}
	results := make([]string, len(list))
	for _, item := range list {
		str, ok := item.(string)
		if ok {
			results = append(results, str)
			continue
		}
		num, ok := item.(int)
		if ok {
			results = append(results, strconv.Itoa(num))
			continue
		}
		// Fall back to general Go print formatting
		results = append(results, fmt.Sprintf("%v", item))
	}
	return results, nil
}

// GetAsIntList does the same as GetAsList but with all values converted to ints.
// Returns an error if a value is found that is not an int and couldn't be parsed into one.
func (doc *HclDocument) GetAsIntList(q string) ([]int, error) {
	list, err := doc.GetAsList(q)
	if err != nil {
		return nil, err
	}
	results := make([]int, len(list))
	for _, item := range list {
		num, ok := item.(int)
		if ok {
			results = append(results, num)
			continue
		}
		str, ok := item.(string)
		if ok {
			num, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("failed to parse '%s' into an integer", str)
			}
			results = append(results, num)
			continue
		}
		return nil, fmt.Errorf("value '%v' is not an integer and could not be parsed into one", item)
	}
	return results, nil
}