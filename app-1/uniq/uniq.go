package uniq

import (
	"errors"
	"strings"
)

func trim(str string, options Options) string {
	words := strings.Split(str, " ")

	if options.fields >= len(words) {
		return ""
	}

	words = words[options.fields:]

	result := strings.Join(words, " ")

	if options.chars >= len(words) {
		return ""
	}

	return result[options.chars:]
}

func Unique(lines []string, options Options) ([]string, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	if (options.count && options.repeat) || (options.count && options.unique) || (options.unique && options.repeat) {
		err := errors.New("Not Compitable arguments choose 1")
		return nil, err
	}

	result := make([]string, 0)

	var err error

	previous := lines[0]

	count := 1

	for _, data := range lines[1:] {
		current := trim(data, options)
		previous = trim(previous, options)

		if current == previous {
			count++
		} else if options.ignore {
			if !strings.EqualFold(current, previous) {
				// result ++
				count = 1
				previous = data
			}
		} else {

		}
	}

	return result, err
}
