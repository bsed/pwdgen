// Copyright 2012 <MortalSkulD@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	majorVersion = 2
	minorVersion = 4
)

var (
	flagSalt          = flag.String("salt", "", "Set site salt.")
	flagSiteSalt      = flag.String("site-salt", "", "Set site salt.")
	flagEncryptKey    = flag.String("encrypt-key", "", "Set encrypt key.")
	flagEncryptSalt   = flag.String("encrypt-salt", "", "Set encrypt salt.")
	flagGenConfig     = flag.String("gen-config", "", "Generate pwdgen config file(*.ini).")
	flagKeepassConfig = flag.String("keepass-config", "", "Generate KeePass 1.x CSV.")

	version = flag.Bool("version", false, "Show version and exit.")
	help    = flag.Bool("help", false, "Show usage and exit.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pwdgen [options]... [id]...\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "%s\n", `
Algorithm:
  base58(sha512(md5hex(encryptKey+encryptSalt)+siteId+siteSalt)[0:16]

Example:
  pwdgen id0
  pwdgen id0 id1 id2

  pwdgen -salt=site0 id0
  pwdgen -salt=site0 id0 id1 id2

  pwdgen -encrypt-key=111 id0
  pwdgen -encrypt-key=111 id0 id1 id2

  pwdgen -encrypt-key=111 -salt=site0 id0 id1
  pwdgen -encrypt-key=111 -salt=site0 id0 id1
  pwdgen -encrypt-key=111 -salt=site0 id0 id1

  # Generate config.ini template
  pwdgen -gen-config=config.ini

  # KeePass: See config.ini
  # output: *.ini -> *.keepass1x.csv
  pwdgen -keepass-config=config.ini
  pwdgen -keepass-config=config.ini -encrypt-key=111
  pwdgen -keepass-config=config.ini -encrypt-key=111 -salt=fuckcsdn

  pwdgen -version
  pwdgen -help
  pwdgen -h


Report bugs to <MortalSkulD@gmail.com.`)
	}
}

func parseCmdLine() {
	flag.Parse()

	if *version {
		fmt.Printf("pwdgen-%d.%d\n", majorVersion, minorVersion)
		os.Exit(0)
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *flagSalt != "" {
		*flagSiteSalt = *flagSalt
	}
}

// base58(sha512(md5hex(flagEncryptKey+flagEncryptSalt)+site_id+flagSiteSalt)[0:16]
func PwdGen(site_id, flagSiteSalt, flagEncryptKey, flagEncryptSalt string) string {
	md5 := md5.New()
	md5.Write([]byte(flagEncryptKey + flagEncryptSalt))
	md5Hex := fmt.Sprintf("%x", md5.Sum(nil))

	sha := sha512.New()
	sha.Write([]byte(md5Hex + site_id + flagSiteSalt))
	shaSum := sha.Sum(nil)

	pwd := EncodeBase58(shaSum)[0:16]
	return string(pwd)
}

func main() {
	parseCmdLine()

	// generate config file
	if genConfig := *flagGenConfig; genConfig != "" {
		if !strings.HasSuffix(strings.ToLower(genConfig), ".ini") {
			genConfig += ".ini"
		}
		err := ioutil.WriteFile(genConfig, []byte(defaultConfigFile), 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Write config failed: %v\n", err)
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if flag.NArg() <= 0 && len(*flagKeepassConfig) <= 0 {
		flag.Usage()
		os.Exit(0)
	}

	if len(*flagEncryptKey) <= 0 {
		fmt.Printf("Encryption key: ")
		(*flagEncryptKey) = string(GetPasswdMasked())
		if len(*flagEncryptKey) <= 0 {
			fmt.Fprintf(os.Stderr, "ERROR: Key must be at least 1 characters.\n")
			os.Exit(-1)
		}
	}

	if *flagKeepassConfig != "" {
		dict, err := LoadIniFile(*flagKeepassConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Load <%s> failed.\n", *flagKeepassConfig)
			os.Exit(-1)
		}

		csvName := strings.Replace(*flagKeepassConfig, ".ini", ".keepass1x.csv", -1)
		file, err := os.Create(csvName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Create <%s> failed.\n", csvName)
			os.Exit(-1)
		}
		defer file.Close()

		// KeePass 1.x csv head line
		fmt.Fprintf(file, `"%s","%s","%s","%s","%s"`+"\n",
			"Account", "Login Name", "Password", "Web Site", "Comments",
		)

		keepass_site_list := dict.GetSections()
		for i := 0; i < len(keepass_site_list); i++ {
			keepass_site_name := keepass_site_list[i]
			keepass_site_id := dict[keepass_site_name]["LoginName"]
			keepass_site_salt := keepass_site_name + dict[keepass_site_name]["SiteSalt"]
			keepass_site_url := dict[keepass_site_name]["WebSite"]
			keepass_site_comments := dict[keepass_site_name]["Comments"]
			keepass_site_pwd := PwdGen(keepass_site_id, keepass_site_salt, *flagEncryptKey, *flagEncryptSalt)

			fmt.Fprintf(file, `"%s","%s","%s","%s","%s"`+"\n",
				keepass_site_name, keepass_site_id, keepass_site_pwd,
				keepass_site_url, keepass_site_comments,
			)
		}
	} else {
		for i := 0; i < flag.NArg(); i++ {
			password := PwdGen(flag.Arg(i), *flagSiteSalt, *flagEncryptKey, *flagEncryptSalt)
			fmt.Printf("%s\n", password)
		}
	}
}
