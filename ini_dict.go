// Copyright 2012 <MortalSkulD@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Dict map[string]map[string]string

func newDict() Dict {
	dict := make(map[string]map[string]string)
	return Dict(dict)
}

func (dict Dict) GetSections() []string {
	sections := make([]string, 0, len(dict))
	for sec, _ := range dict {
		sections = append(sections, sec)
	}
	sort.StringSlice(sections).Sort()
	return sections
}
func (dict Dict) GetKeys(section string) []string {
	sec, ok := dict[section]
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(sec))
	for key, _ := range sec {
		keys = append(keys, key)
	}
	sort.StringSlice(keys).Sort()
	return keys
}

func (dict Dict) GetString(section, key string) (string, error) {
	sec, ok := dict[section]
	if !ok {
		return "", errors.New("Section Not Found")
	}
	value, ok := sec[key]
	if !ok {
		return "", errors.New("Key Not Found")
	}
	return value, nil
}

func (dict Dict) GetBool(section, key string) (bool, error) {
	sec, ok := dict[section]
	if !ok {
		return false, errors.New("Section Not Found")
	}
	value, ok := sec[key]
	if !ok {
		return false, errors.New("Key Not Found")
	}
	v := value[0]
	if v == 'y' || v == 'Y' || v == '1' || v == 't' || v == 'T' {
		return true, nil
	}
	if v == 'n' || v == 'N' || v == '0' || v == 'f' || v == 'F' {
		return false, nil
	}
	return false, errors.New("Parse Bool Failed")
}

func (dict Dict) GetInt(section, key string) (int, error) {
	sec, ok := dict[section]
	if !ok {
		return 0, errors.New("Section Not Found")
	}
	value, ok := sec[key]
	if !ok {
		return 0, errors.New("Key Not Found")
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("Parse Int Failed")
	}
	return i, nil
}

func (dict Dict) GetFloat(section, key string) (float64, error) {
	sec, ok := dict[section]
	if !ok {
		return 0, errors.New("Section Not Found")
	}
	value, ok := sec[key]
	if !ok {
		return 0, errors.New("Key Not Found")
	}
	d, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.New("Parse Float Failed")
	}
	return d, nil
}

func (dict Dict) String() string {
	var b bytes.Buffer
	sections := dict.GetSections()
	for _, sec := range sections {
		if sec != "" {
			fmt.Fprintf(&b, "[%s]\n", sec)
		}
		keys := dict.GetKeys(sec)
		for _, key := range keys {
			if key != "" {
				fmt.Fprintf(&b, "%s=%s\n", key, dict[sec][key])
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String()
}

func (dict *Dict) setSection(section string) {
	section = strings.TrimFunc(section, unicode.IsSpace)
	if section == "" {
		return
	}

	if _, ok := (*dict)[section]; !ok {
		(*dict)[section] = make(map[string]string)
	}
}
func (dict *Dict) setString(section, key string, value string) {
	section = strings.TrimFunc(section, unicode.IsSpace)
	if section == "" {
		return
	}

	if _, ok := (*dict)[section]; !ok {
		(*dict)[section] = make(map[string]string)
	}
	(*dict)[section][key] = value
}
