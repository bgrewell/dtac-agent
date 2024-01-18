package utilities

import (
	"errors"
	"strconv"
	"strings"
)

// ExtractString extracts a single string value from the parameters.
func ExtractString(parameters map[string][]string, key string) (string, error) {
	values, ok := parameters[key]
	if !ok || len(values) == 0 {
		return "", errors.New("key not found or no values available")
	}
	return values[0], nil
}

// ExtractStrings extracts multiple string values from the parameters.
func ExtractStrings(parameters map[string][]string, key string) ([]string, error) {
	values, ok := parameters[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return values, nil
}

// ExtractBool extracts a single bool value from the parameters.
func ExtractBool(parameters map[string][]string, key string) (bool, error) {
	value, err := ExtractString(parameters, key)
	if err != nil {
		return false, err
	}
	return evaluateBool(value)
}

// ExtractBools extracts multiple bool values from the parameters.
func ExtractBools(parameters map[string][]string, key string) ([]bool, error) {
	values, err := ExtractStrings(parameters, key)
	if err != nil {
		return nil, err
	}
	bools := make([]bool, len(values))
	for i, value := range values {
		bools[i], err = evaluateBool(value)
		if err != nil {
			return nil, err
		}
	}
	return bools, nil
}

// ExtractInt extracts a single int value from the parameters.
func ExtractInt(parameters map[string][]string, key string) (int, error) {
	value, err := ExtractString(parameters, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// ExtractInts extracts multiple int values from the parameters.
func ExtractInts(parameters map[string][]string, key string) ([]int, error) {
	values, err := ExtractStrings(parameters, key)
	if err != nil {
		return nil, err
	}
	ints := make([]int, len(values))
	for i, value := range values {
		ints[i], err = strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
	}
	return ints, nil
}

func evaluateBool(input string) (value bool, err error) {
	trueValues := map[string]struct{}{
		"true": {},
		"t":    {},
		"tr":   {},
		"1":    {},
		"yes":  {},
	}
	falseValues := map[string]struct{}{
		"false": {},
		"f":     {},
		"0":     {},
		"no":    {},
	}

	lowerInput := strings.ToLower(input)

	if _, exists := trueValues[lowerInput]; exists {
		return true, nil
	}
	if _, exists := falseValues[lowerInput]; exists {
		return false, nil
	}

	return false, errors.New("invalid boolean value")
}
