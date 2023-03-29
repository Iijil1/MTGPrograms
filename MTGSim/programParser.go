package main

import (
	"strconv"
	"strings"
)

type parsedLine struct {
	multiplicity     int
	objectType       string
	typeReplacements map[int]string
	keywords         []string
	displayGroup1    string
	displayGroup2    string
}

//returns nil if the string is not a valid line of the program
func parseProgramLine(programLine string) *parsedLine {
	words := strings.Split(programLine, " ")
	if len(words) < 2 || len(words[0]) < 2 || words[0][len(words[0])-1] != 'x' {
		//This line is to short or doesn't start with Nx, so we ignore it for the program
		return nil
	}

	result := parsedLine{
		multiplicity:     0,
		objectType:       words[1],
		typeReplacements: map[int]string{},
		keywords:         []string{},
		displayGroup1:    "",
		displayGroup2:    "",
	}

	var err error
	result.multiplicity, err = strconv.Atoi(words[0][0 : len(words[0])-1])
	if err != nil {
		//The line didn't start with Nx after all, since we couldn't parse N. We ignore it for the program
		return nil
	}

	for _, word := range words[2:] {
		if len(word) == 0 {
			continue
		}
		splitWord := strings.Split(word, ":")
		switch len(splitWord) {
		case 1:
			switch word[0] {
			case '<':
				result.displayGroup1 = word[1:]
			case '>':
				result.displayGroup2 = word[1:]
			default:
				result.keywords = append(result.keywords, word)
			}
		case 2:
			changedType, err := strconv.Atoi(splitWord[0])
			if err != nil {
				return nil //some invalid word got added
			}
			result.typeReplacements[changedType] = splitWord[1]
		default:
			return nil //some invalid word got added
		}
	}

	return &result
}
