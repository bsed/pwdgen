## Password Generator

Password generator, support [KeePass 1.x](http://keepass.info/) format file.

## Install from Standalone File

- Windows: download [pwdgen_windows.go](https://github.com/chai2010/pwdgen/blob/master/pwdgen_windows.go)
- Darwin: download [pwdgen_darwin.go](https://github.com/chai2010/pwdgen/blob/master/pwdgen_darwin.go)
- Linux: download [pwdgen_linux.go](https://github.com/chai2010/pwdgen/blob/master/pwdgen_linux.go)

Then `go run pwdgen_$GOOS.go` or `go build -o pwdgen.exe pwdgen_$GOOS.go`.

## Install With `go get`

	go get github.com/webfd/pwdgen

## Usage

	$ pwdgen id0
	$ pwdgen id0 id1 id2

	$ pwdgen -salt=site0 id0
	$ pwdgen -salt=site0 id0 id1 id2

	$ pwdgen -encrypt-key=111 id0
	$ pwdgen -encrypt-key=111 id0 id1 id2

	$ pwdgen -encrypt-key=111 -salt=site0 id0 id1
	$ pwdgen -encrypt-key=111 -salt=site0 id0 id1
	$ pwdgen -encrypt-key=111 -salt=site0 id0 id1

	$ # Generate config.ini template
	$ pwdgen -gen-config=config.ini

	$ # *.ini -> *.keepass1x.csv
	$ pwdgen -keepass-config=config.ini
	$ pwdgen -keepass-config=config.ini -encrypt-key=111

	$ pwdgen -version
	$ pwdgen -help
	$ pwdgen -h


## Use in Go

	package main

	import (
		"fmt"
		pwdgen "github.com/webfd/pwdgen"
	)
	func main() {
		fmt.Println(pwdgen.PwdGen("id0", "site0", "111", "fuckcsdn"))
		// Output: 2jNXfMGoXTSK9pFS
	}

See ([Issue4210](https://code.google.com/p/go/issues/detail?id=4210)).

## Algorithm

	base58(sha512(md5hex(encryptKey+encryptSalt)+siteId+siteSalt)[0:16]


## Reference

* http://godoc.org/github.com/webfd/pwdgen
* http://keepass.info/
