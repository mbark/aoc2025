package util

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
)

func GetInput(day int) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2025/day/%d/input", day), nil)
	must(err)
	req.Header.Set("cookie", "session="+os.Getenv("SESSION"))

	res, err := http.DefaultClient.Do(req)
	must(err)

	if res.StatusCode != 200 {
		d, err := httputil.DumpResponse(res, true)
		fmt.Println(string(d), err)
		panic(fmt.Sprintf("non-200 status code: %d", res.StatusCode))
	}

	bytes, err := io.ReadAll(res.Body)
	must(err)

	return string(bytes)
}

func RegexCaptureGroups(re *regexp.Regexp, input string) []map[string]string {
	matches := re.FindAllStringSubmatch(input, -1)
	groupNames := re.SubexpNames()

	var results []map[string]string
	for _, match := range matches {
		result := make(map[string]string)
		for i, name := range groupNames {
			if i > 0 && name != "" {
				result[name] = match[i]
			}
		}
		results = append(results, result)
	}

	return results
}

func RegexSubexps(re *regexp.Regexp, input string) []string {
	matches := re.FindAllStringSubmatch(input, -1)

	var results []string
	for _, match := range matches {
		results = append(results, match[1:]...)
	}

	return results
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
