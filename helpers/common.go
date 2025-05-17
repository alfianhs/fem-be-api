package helpers

import "strings"

func StringReplacer(content string, data map[string]string) string {
	for key, value := range data {
		content = strings.ReplaceAll(content, "{{"+key+"}}", value)
	}
	return content
}

func ExtractIds[T any](items []T, extractor func(T) string) []string {
	seen := map[string]bool{}
	var result []string

	for _, item := range items {
		id := extractor(item)
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}

	return result
}

func InArrayString(array []string, val string) bool {
	for _, a := range array {
		if a == val {
			return true
		}
	}
	return false
}

func ArrayStringtoString(array []string) string {
	return strings.Join(array, ", ")
}

func Intersect(a, b []string) []string {
	set := make(map[string]struct{})
	for _, val := range a {
		set[val] = struct{}{}
	}
	var common []string
	for _, val := range b {
		if _, ok := set[val]; ok {
			common = append(common, val)
		}
	}
	return common
}
