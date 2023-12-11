package protoparser

/*
* this package does some rudimentary parsing of the content of a file
* On occassion, if one "just" has the content of a .proto file one
* needs to work out its package and imports etc in order to store it in a location
* suitable for protoc.
* thus this parser should only be used where protoc cannot (yet)
 */

import (
	"fmt"
	"strings"
)

func Parse(content string) (*ProtoParsed, error) {
	lines := strings.Split(content, "\n")
	res := &ProtoParsed{Content: content}
	comment := false
	for _, line := range lines {
		if comment {
			idx := strings.Index(line, "*/")
			if idx == -1 {
				continue
			}
			comment = false
			line = line[:idx]
		}

		idx := strings.Index(line, "//")
		if idx != -1 {
			line = line[:idx]
		}
		idx = strings.Index(line, "/*")
		if idx != -1 {
			comment = true
			line = line[:idx]
		}

		if strings.Contains(line, "java_package") && strings.Contains(line, "option") && res.JavaPackage == "" {
			res.JavaPackage = lastElement(line)
		} else if strings.Contains(line, "package") {
			res.GoPackage = lastElement(line)
		} else if strings.Contains(line, "import") {
			i := lastElement(line)
			res.Imports = append(res.Imports, i)
		}
	}
	return res, nil
}

// return last token before the ';'
// for example: "package h2gproxy;" returns "h2gproxy"
func lastElement(input string) string {
	input = strings.Trim(input, " ")
	if len(input) <= 1 {
		return ""
	}
	if input[len(input)-1] != ';' {
		fmt.Printf("Invalid input: \"%s\"\n", input)
		return input
	}
	i := len(input) - 2
	s := ""
	for i >= 0 {
		c := input[i]
		if c == ' ' {
			break
		}
		s = fmt.Sprintf("%c%s", c, s)
		i--
	}
	s = strings.Trim(s, "\"")
	return s
}
























































































