## Changelog (Current version: 3.0.2)

-----------------

### 3.0.2 (2016 Feb 29)

* [7d747a8] Merge pull request #13 from bitrise-io/develop
* [b5ee8a2] typo fix
* [9a0c3f0] v3.0.1

### 3.0.1 (2016 Feb 29)

* [b6b90b1] Merge pull request #12 from bitrise-io/develop
* [354a857] release configs
* [a1c0329] bitrise.yml format update
* [ca2e682] Merge pull request #11 from bitrise-io/develop
* [f21dc9b] default value in message input
* [bd56019] new input: message_format
* [aaa6823] step cleanup - removed old stuffs
* [b4e1088] gitignore fix
* [a24fdd5] Merge pull request #6 from bazscsa/patch-1
* [dea4664] Update step.yml
* [2141c7a] updated to V2 step.yml format

### 3.0.0 (2015 Aug 17)

* [d9e6506] Merge pull request #5 from rptabo/failure-color

### 2.1.0 (2015 Aug 13)

* [6b568bd] Allow different color for failed builds

### 2.0.0 (2015 Jul 14)

* [f299327] Merge pull request #4 from gkiki90/master
* [96afb2f] Merge branch 'release/code_style'
* [1121761] code style
* [a541c42] Merge branch 'master' of github.com:bitrise-io/steps-hipchat
* [ac9a91d] ci.sh
* [d4676fb] ci.sh fixes
* [80960fe] markdownlog added with depamn, ci.sh for testing
* [6893593] updated markdownlog
* [9876470] code cleaning
* [beb2b4b] Merge branch 'master' of github.com:gkiki90/steps-hipchat
* [da2ddad] translated to golang
* [c62c193] step translate to go
* [d62c1ee] step logic in go
* [61c7c78] log fixes
* [e4a41e8] logging function fix
* [a48e38d] formatted output logging
* [fee1c1b] output format
* [fc9b2c3] from name fix
* [55503e3] fixed from name in output
* [7593b5e] build failed message, and from name fix
* [d7d3c4f] fixed formatted output message
* [1a845cd] fixed error handling, fixed MESSAGE required
* [84ae7f7] fixed exit
* [5832179] code style, depman
* [d3d7bc9] step fixes mL - source back to bitrise-io
* [8d4b721] step fixes
* [a4c44b5] HIPCAHT_FROM_NAME is required
* [678065a] HIPCHAT_FROMNAME renamed to HIPCHAT_FROM_NAME
* [0c96e2c] fixed param setup
* [fb24224] yml fixes
* [57610e3] isBuildFailedMode handling
* [f0975cf] Update README.md
* [ede1067] further step.YML update
* [5c5d1c2] mL updates
* [865e229] README.md update
* [2d6a6af] Update README.md
* [c351bca] readme and step YML revision
* [53df481] Merge pull request #3 from erosdome/master
* [e1ab20d] Update step.yml
* [87c9e81] Update step.yml
* [ff6d524] Update README.md

### 1.1.6 (2014 Sep 05)

* [bea487d] formatted output : success summary -> includes roomID and better (?) formatting
* [889da45] Merge pull request #2 from erosdome/master

### 1.1.5 (2014 Sep 05)

* [43a861b] readme update
* [fb1708b] readme update
* [08b6608] handle invalid token & room ids; readme.md update
* [7a31b20] initial
* [7895ace] Merge branch 'release/1.1.4'

### 1.1.4 (2014 Jul 14)

* [a1952bb] newline
* [0a9d93c] Merge branch 'release/1.1.3'

### 1.1.3 (2014 Jul 14)

* [7bbf306] step_test moved from subfolder to root
* [b9bf7a1] Merge branch 'release/1.1.2'

### 1.1.2 (2014 Jul 11)

* [5dabd59] readme fix: hip chat color is actually already supported
* [9b52fc9] Merge branch 'release/1.1.1'

### 1.1.1 (2014 Jul 11)

* [e896807] simple renames to the new name, Bitrise
* [d834efa] Merge branch 'tomfurrier-master'
* [b3f9c35] Merge branch 'master' of https://github.com/tomfurrier/hipchat into tomfurrier-master
* [24df42b] sync with latest
* [7a16e01] Merge branch 'master' of https://github.com/tomfurrier/hipchat
* [3e72515] added tests
* [cca8eaa] Merge branch 'release/1.1.0'

### 1.1.0 (2014 Jun 17)

* [56b23da] completely rewritten, it's now pure bash script instead of ruby+gem; step descriptor 'step.yml' added; message color support
* [8f3abd6] Merge branch 'release/1.0.2'

### 1.0.2 (2014 Jun 17)

* [9ba47b6] removed not required ssl cert PEM download; step descriptor 'step.yml' added
* [7808dc1] Merge branch 'release/1.0.1'

### 1.0.1 (2014 May 08)

* [d6b01ac] read me with how to test/run locally + comment about the SSL fix
* [80df7a6] Merge branch 'release/1.0.0'
* [89ea468] Merge branch 'release/1.0.0' into develop

### 1.0.0 (2014 May 07)

* [b5f5d01] proper communication error reporting - return non 0 exit code if error happens
* [055b941] Merge branch 'release/0.9.4'

### 0.9.4 (2014 Apr 08)

* [a33b0fb] don't set --path - will be handled by the setup
* [502579c] Merge branch 'release/0.9.3'

### 0.9.3 (2014 Apr 08)

* [3647cc8] bundle install - specify bundle/gem folder path
* [bab99be] unused opt parse require removed
* [81de4f4] Merge branch 'release/0.9.2'

### 0.9.2 (2014 Apr 05)

* [e76176e] for testing : hip chat message can be specified through environment
* [4d47c0c] base hipchat step, it can only send a hardcoded "Quick test" message

-----------------

Updated: 2016 Feb 29