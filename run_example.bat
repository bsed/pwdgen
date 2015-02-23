:: Copyright 2010 <MortalSkulD@gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

go build
:: ----------------------------------------------------------------------------

pwdgen --help
pwdgen --version

:: ----------------------------------------------------------------------------
:: encrypt_key: 111

pwdgen --site-salt=site0 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site0 --encrypt-key=111 --encrypt-salt=fuckcsdn id1
pwdgen --site-salt=site0 --encrypt-key=111 --encrypt-salt=fuckcsdn id3

pwdgen --site-salt=site1 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site1 --encrypt-key=111 --encrypt-salt=fuckcsdn id1

pwdgen --site-salt=site2 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site3 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site4 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site5 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site6 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site7 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site8 --encrypt-key=111 --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site9 --encrypt-key=111 --encrypt-salt=fuckcsdn id0

:: ----------------------------------------------------------------------------
:: encrypt_key: abc

pwdgen --site-salt=site0 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site0 --encrypt-key=abc --encrypt-salt=fuckcsdn id1
pwdgen --site-salt=site0 --encrypt-key=abc --encrypt-salt=fuckcsdn id3

pwdgen --site-salt=site1 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site1 --encrypt-key=abc --encrypt-salt=fuckcsdn id1

pwdgen --site-salt=site2 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site3 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site4 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site5 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site6 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site7 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site8 --encrypt-key=abc --encrypt-salt=fuckcsdn id0
pwdgen --site-salt=site9 --encrypt-key=abc --encrypt-salt=fuckcsdn id0

:: ----------------------------------------------------------------------------
:: KeePass

:: config.ini -> config.keepass1x.csv
pwdgen --encrypt-key=111 --encrypt-salt=fuckcsdn --keepass-config=config.ini

:: ----------------------------------------------------------------------------

PAUSE
