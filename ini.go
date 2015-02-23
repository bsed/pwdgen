// Copyright 2012 <MortalSkulD@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	regDoubleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$")
	regSingleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*'([^']*)'$")
	regNoQuote     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([^#;]+)")
	regNoValue     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([#;].*)?")
)

func LoadIniFile(filename string) (dict Dict, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return loadIniData(reader, filename)
}
func LoadIniString(str string) (dict Dict, err error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(str)))
	return loadIniData(reader, "[string]")
}

func loadIniData(reader *bufio.Reader, filename string) (dict Dict, err error) {

	dict = newDict()
	lineno := 0
	section := ""

	for err == nil {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineno++
		if len(l) == 0 {
			continue
		}
		line := strings.TrimFunc(string(l), unicode.IsSpace)

		for line[len(line)-1] == '\\' {
			line = line[:len(line)-1]
			l, _, err := reader.ReadLine()
			if err != nil {
				return nil, err
			}
			line += strings.TrimFunc(string(l), unicode.IsSpace)
		}

		section, err = parseIniLine(&dict, section, line)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("%v: '%s:%d'.", err, filename, lineno),
			)
		}
	}

	return
}

func parseIniLine(dict *Dict, section, line string) (string, error) {
	// commets
	if line[0] == '#' || line[0] == ';' {
		return section, nil
	}

	// section name
	if line[0] == '[' && line[len(line)-1] == ']' {
		section = strings.TrimFunc(line[1:len(line)-1], unicode.IsSpace)
		dict.setSection(section)
		return section, nil
	}

	// key = value
	if section != "" {
		if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
			dict.setString(section, m[0][1], m[0][2])
			return section, nil
		} else if m = regSingleQuote.FindAllStringSubmatch(line, 1); m != nil {
			dict.setString(section, m[0][1], m[0][2])
			return section, nil
		} else if m = regNoQuote.FindAllStringSubmatch(line, 1); m != nil {
			dict.setString(section, m[0][1], strings.TrimFunc(m[0][2], unicode.IsSpace))
			return section, nil
		} else if m = regNoValue.FindAllStringSubmatch(line, 1); m != nil {
			dict.setString(section, m[0][1], "")
			return section, nil
		}
	}

	return section, errors.New("syntax error")
}
