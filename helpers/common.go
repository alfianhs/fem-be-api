package helpers

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

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

func SetToStartOfDayUTC(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func SetToEndOfDayUTC(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 23, 59, 59, 0, time.UTC)
}

func SetToStartOfDayWIB(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t = t.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func SetToEndOfDayWIB(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t = t.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
}

func IsSameDateWIB(a, b time.Time) bool {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	a = a.In(loc)
	b = b.In(loc)
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func ConvertUTCToWIB(t time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return t
	}
	return t.In(loc)
}

func FormatDateWIB(t time.Time, layout string) string {
	return ConvertUTCToWIB(t).Format(layout)
}

func GenerateInvoiceExternalId(count int64) string {
	prefix := "TRX"
	randomChar, _ := GenerateSecureRandomChar(3)
	timestamp := time.Now().Format("20060102")
	counter := fmt.Sprintf("%03d", count+1)

	return fmt.Sprintf("%s-%s%s-%s", prefix, timestamp, randomChar, counter)
}

func SanitizeString(str string) string {
	re := regexp.MustCompile(`[^\w\d_-]`)
	return re.ReplaceAllString(str, "_")
}
