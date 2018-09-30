package tagsugar

import (
	"log"
	"strings"
)

type tagOptions map[string]string

var supportTag = map[string]int{
	"-":           1,
	"initial":     1,
	"url":         2,
	"assign_to":   2,
	"assign_type": 2,
}

func parseTag(tag string) tagOptions {
	options := tagOptions{}
	for tag != "" {
		var next string
		i := strings.Index(tag, ";")
		if i > 0 {
			tag, next = tag[:i], tag[i+1:]
		}
		if supportTag[tag] == 1 {
			options[tag] = ""
		} else if i := strings.Index(tag, "("); i > 0 && strings.Index(tag, ")") == len(tag)-1 {
			name := tag[:i]
			if supportTag[name] == 2 {
				options[name] = tag[i+1 : len(tag)-1]
			}
		} else {
			log.Println("unsupported ts tag", tag)
		}

		tag = next
	}
	return options
}
