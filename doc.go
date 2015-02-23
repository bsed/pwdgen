// Copyright 2014 MortalSkulD@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Password generator
//
// Usage:
//	pwdgen [options]... [id]...
//
// Algorithm:
//	base58(sha512(md5hex(flagEncryptKey+flagEncryptSalt)+site_id+flagSiteSalt)[0:16]
//
// Example:
//	pwdgen id0
//	pwdgen id0 id1 id2

//	pwdgen -salt=site0 id0
//	pwdgen -salt=site0 id0 id1 id2
//
//	pwdgen -encrypt-key=111 id0
//	pwdgen -encrypt-key=111 id0 id1 id2
//
//	pwdgen -encrypt-key=111 -salt=site0 id0 id1
//	pwdgen -encrypt-key=111 -salt=site0 id0 id1
//	pwdgen -encrypt-key=111 -salt=site0 id0 id1
//
//	# Generate config.ini template
//	pwdgen -gen-config=config.ini
//
//	# KeePass: See config.ini
//	# output: *.ini -> *.keepass1x.csv
//	pwdgen -keepass-config=config.ini
//	pwdgen -keepass-config=config.ini -encrypt-key=111
//
//	pwdgen -version
//	pwdgen -help
//	pwdgen -h
//
//
// Report bugs to <MortalSkulD@gmail.com>.
package main
