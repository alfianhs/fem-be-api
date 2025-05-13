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
