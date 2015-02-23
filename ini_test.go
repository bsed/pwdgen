// Copyright 2012 <MortalSkulD@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

var iniStr = `
# Copyright 2011-2013 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

[gmail]
LoginName = abc
WebSite = http://mail.google.com
Comments = 

[taobao]
LoginName = 123
WebSite = http://www.taobao.com
Comments = 
`

var iniGolden Dict = map[string]map[string]string{
	"gmail": map[string]string{
		"LoginName": "abc",
		"WebSite":   "http://mail.google.com",
		"Comments":  "",
		"Bad":       "",
	},
	"taobao": map[string]string{
		"LoginName": "123",
		"WebSite":   "http://www.taobao.com",
		"Comments":  "",
	},
	"Bad": map[string]string{
		"LoginName": "",
		"WebSite":   "",
		"Comments":  "",
	},
}

func TestDict(t *testing.T) {
	dict, err := LoadIniString(iniStr)
	if err != nil {
		t.Errorf("LoadString(iniStr) fail: %v", err)
	}

	section_list := dict.GetSections()
	for i := 0; i < len(section_list); i++ {
		if _, ok := iniGolden[section_list[i]]; !ok {
			t.Errorf("section[%s] not found", section_list[i])
		}
	}

	for sec, sec_values := range iniGolden {
		for key, value := range sec_values {
			if dict[sec][key] != value {
				t.Errorf("dict([%s]:%s): Need=%v, Got=%v", sec, key, value, dict[sec][key])
			}
		}
	}
}

func ExampleDict() {
	dict, _ := LoadIniString(`
# Copyright 2011-2013 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

[gmail]
LoginName = abc
WebSite = http://mail.google.com
Comments = 

[taobao]
LoginName = 123
WebSite = http://www.taobao.com
Comments = 
`)

	fmt.Println(dict["gmail"]["LoginName"])
	fmt.Println(dict["gmail"]["WebSite"])
	fmt.Println(dict["gmail"]["Comments"])
	fmt.Println(dict["taobao"]["LoginName"])
	fmt.Println(dict["taobao"]["WebSite"])
	fmt.Println(dict["taobao"]["Comments"])
	// Output:
	// abc
	// http://mail.google.com
	//
	// 123
	// http://www.taobao.com
	//
}
