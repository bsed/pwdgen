// Copyright 2014 MortalSkulD@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Auto Generated By 'go generate', DONOT EDIT!!!

// +build ingore
// +build windows

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

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf16"
	"unsafe"
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

const (
	base58Table = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func base58Hash(ba []byte) []byte {
	sha := sha256.New()
	sha2 := sha256.New() // hash twice
	ba = sha.Sum(ba)
	return sha2.Sum(ba)
}

func EncodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}

	// Expected size increase from base58Table conversion is approximately 137%, use 138% to be safe
	ri := len(ba) * 138 / 100
	ra := make([]byte, ri+1)

	x := new(big.Int).SetBytes(ba) // ba is big-endian
	x.Abs(x)
	y := big.NewInt(58)
	m := new(big.Int)

	for x.Sign() > 0 {
		x, m = x.DivMod(x, y, m)
		ra[ri] = base58Table[int32(m.Int64())]
		ri--
	}

	// Leading zeroes encoded as base58Table zeros
	for i := 0; i < len(ba); i++ {
		if ba[i] != 0 {
			break
		}
		ra[ri] = '1'
		ri--
	}
	return ra[ri+1:]
}

func DecodeBase58(ba []byte) []byte {
	if len(ba) == 0 {
		return nil
	}

	x := new(big.Int)
	y := big.NewInt(58)
	z := new(big.Int)
	for _, b := range ba {
		v := strings.IndexRune(base58Table, rune(b))
		z.SetInt64(int64(v))
		x.Mul(x, y)
		x.Add(x, z)
	}
	xa := x.Bytes()

	// Restore leading zeros
	i := 0
	for i < len(ba) && ba[i] == '1' {
		i++
	}
	ra := make([]byte, i+len(xa))
	copy(ra[i:], xa)
	return ra
}

func EncodeBase58Check(ba []byte) []byte {
	// add 4-byte hash check to the end
	hash := base58Hash(ba)
	ba = append(ba, hash[:4]...)
	ba = EncodeBase58(ba)
	return ba
}

func DecodeBase58Check(ba []byte) bool {
	ba = DecodeBase58(ba)
	if len(ba) < 4 || ba == nil {
		return false
	}

	k := len(ba) - 4
	hash := base58Hash(ba[:k])
	for i := 0; i < 4; i++ {
		if hash[i] != ba[k+i] {
			return false
		}
	}
	return true
}

// GetPasswd returns the password read from the terminal without echoing input.
// The returned byte array does not include end-of-line characters.
func GetPasswd() []byte {
	return getPasswd(false)
}

// GetPasswdMasked returns the password read from the terminal, echoing asterisks.
// The returned byte array does not include end-of-line characters.
func GetPasswdMasked() []byte {
	return getPasswd(true)
}

// getPasswd returns the input read from terminal.
// If masked is true, typing will be matched by asterisks on the screen.
// Otherwise, typing will echo nothing.
func getPasswd(masked bool) []byte {
	var pass, bs, mask []byte
	if masked {
		bs = []byte("\b \b")
		mask = []byte("*")
	}

	for {
		if v := getch(); v == 127 || v == 8 {
			if l := len(pass); l > 0 {
				pass = pass[:l-1]
				os.Stdout.Write(bs)
			}
		} else if v == 13 || v == 10 {
			break
		} else {
			pass = append(pass, v)
			os.Stdout.Write(mask)
		}
	}
	println()
	return pass
}

var defaultConfigFile = `
# Copyright 2012 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#
# Gen KeePass CSV 1.x Format File
#
# keepass_site_name := section_name
# keepass_site_id := section["LoginName"]
# keepass_site_salt := section_name + section["SiteSalt"]
# keepass_site_url := section["WebSite"]
# keepass_site_comments := section["Comments"]
# keepass_site_pwd := pwdgen(keepass_site_id, keepass_site_salt, *encrypt_key, *encrypt_salt)
#
# output:
# *.ini -> *.keepass1x.csv
#
# http://github.org/chai2010/pwdgen
# http://keepass.info
#

# -----------------------------------------------------------------------------

[gmail]
LoginName = abc
WebSite = http://mail.google.com
Comments = abc@gmail.com
SiteSalt = 0x5f3759df

[taobao]
LoginName = 123
WebSite = http://www.taobao.com
Comments = 123@gmail.com
SiteSalt = 

[taobao.A]
LoginName = uer_a
WebSite = http://www.taobao.com
Comments = UserA@taobao.com
SiteSalt = 

[taobao.B]
LoginName = uer_b
WebSite = http://www.taobao.com
Comments = UserB@taobao.com
SiteSalt = 

[taobao.C]
LoginName = uer_c
WebSite = http://www.taobao.com
Comments = UserC@taobao.com
SiteSalt = 

# -----------------------------------------------------------------------------
# --encrypt_key=111 --encrypt-salt=fuckcsdn

[site]
LoginName = id3
WebSite = 
Comments = site0:N9V9FyMJ8tkScBN5
SiteSalt = 0 

[site1]
LoginName = id1
WebSite = 
Comments = 641gNmCY9YFNAQ1p
SiteSalt = 

[site2]
LoginName = id0
WebSite = 
Comments = 4HLACkWRCyDHtqtx
SiteSalt = 

[site3]
LoginName = id0
WebSite = 
Comments = 3eHtu74rMFdeRaVk
SiteSalt = 

[site4]
LoginName = id0
WebSite = 
Comments = 5DSxs623Rciz7bab
SiteSalt = 

[site5]
LoginName = id0
WebSite = 
Comments = 3cfiPrcjdrhwAgM1
SiteSalt = 

[site6]
LoginName = id0
WebSite = 
Comments = 5Las25BPXCjvtywo
SiteSalt = 

[site7]
LoginName = id0
WebSite = 
Comments = 1GK1x3GnRxLSH6DT
SiteSalt = 

[site8]
LoginName = id0
WebSite = 
Comments = 3VqQSgsRRQTeR6vL
SiteSalt = 

[site9]
LoginName = id0
WebSite = 
Comments = 5hHdKchVRPeJkFjU
SiteSalt = 

# -----------------------------------------------------------------------------

`[1:]

var (
	modkernel32        = syscall.NewLazyDLL("kernel32.dll")
	procReadConsole    = modkernel32.NewProc("ReadConsoleW")
	procGetConsoleMode = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = modkernel32.NewProc("SetConsoleMode")
)

func getch() byte {
	var mode uint32
	pMode := &mode
	procGetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pMode)))

	var echoMode, lineMode uint32
	echoMode = 4
	lineMode = 2
	var newMode uint32
	newMode = mode ^ (echoMode | lineMode)

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(newMode))

	line := make([]uint16, 1)
	pLine := &line[0]
	var n uint16
	procReadConsole.Call(
		uintptr(syscall.Stdin), uintptr(unsafe.Pointer(pLine)), uintptr(len(line)),
		uintptr(unsafe.Pointer(&n)),
	)

	// For some reason n returned seems to big by 2 (Null terminated maybe?)
	if n > 2 {
		n -= 2
	}

	b := []byte(string(utf16.Decode(line[:n])))

	procSetConsoleMode.Call(uintptr(syscall.Stdin), uintptr(mode))

	// Not sure how this could happen, but it did for someone
	if len(b) > 0 {
		return b[0]
	} else {
		return 13
	}
}