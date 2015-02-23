// Copyright 2014 MortalSkulD@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen_helper.go -goos=windows -output=pwdgen_windows.go
//go:generate go run gen_helper.go -goos=linux -output=pwdgen_linux.go
//go:generate go run gen_helper.go -goos=darwin -output=pwdgen_darwin.go

package main
