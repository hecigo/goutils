package goutils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Remove Vietnamese accents
func RemoveAccents(s string) string {
	viChars := "ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰựỲỳỴỵỶỷỸỹ"
	enChars := "AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUuYyYyYyYy"

	var r = ""
	for _, c := range s {
		i := strings.IndexRune(viChars, c)
		if i >= 0 {
			i = utf8.RuneCountInString(viChars[:i])
			r += string(enChars[i])
		} else {
			r += string(c)
		}
	}

	return r
}

// Remove special characters
func RemoveSpecialChar(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(regexp.MustCompile(`[^\p{L}\p{N} ]+`).ReplaceAllString(s, ""), "  ", " "))
}

// Convert a string to URL format
func ToURL(s string) string {
	s = strings.ToLower(RemoveAccents(s))
	s = RemoveSpecialChar(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// Clean and split keywords
func CleanKeyword(kw string) []string {
	kw = strings.ToLower(strings.TrimSpace(kw))

	// remove leading numbers
	kw = regexp.MustCompile(`^(\d+)\s+`).ReplaceAllString(kw, "")

	var keywords []string
	keywords = append(keywords, strings.Split(kw, ",")...)

	// remove special characters and trim spaces
	for i, k := range keywords {
		keywords[i] = RemoveSpecialChar(k)
	}

	keywords = Deduplicate(keywords)

	return keywords
}

/*
Generate Elasticsearch query string
View more: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/query-dsl-query-string-query.html
Example: ((\"hcm\") AND (\"quận1\")) OR ((\"hcm\") AND (\"quan1\"))
*/
func ToDSLQueryString( /*Note: Clean keyword firstly*/ freshKw string) string {
	viWords := strings.Split(freshKw, " ")

	var enWords []string
	for _, word := range viWords {
		enWords = append(enWords, RemoveAccents(word))
	}

	viQuery := "(" + strings.Join(viWords, " AND ") + ")"
	enQuery := "(" + strings.Join(enWords, " AND ") + ")"

	return viQuery + " OR " + enQuery
}
