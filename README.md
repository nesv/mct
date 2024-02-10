# mct

Miek's Configuration Tool/Thing.

* [Config Management](https://miek.nl/2024/january/29/config-management/)
* [Config Management, part II: Microcode
  Language](https://miek.nl/2024/february/01/config-management-part-ii-microcode-language/)
  
> [!WARNING]
> This is super experimental.

## Commands

| INSTRUCTION | ARITY | ARGUMENTS         | REMARK             |
|-------------|-------|-------------------|--------------------|
| REM         | 1     | TEXT              | a comment          |
| MKDIR       | 2     | MODE PATH         | create a directory |
| COPY        | 2     | SRC-PATH DST-PATH | copy a file        |
| CHMOD       | 2     | MODE PATH         | set the mode       |
| CHOWN       | 2     | USER/ID PATH      | set the owner      |
| CHGRP       | 2     | GROUP/ID PATH     | set the group      |
| RM          | 1     | PATH              | remove the file    |
| ADDPKG      | 1     | PACKAGE           | install a package  |
| DELPKG      | 1     | PACKAGE           | remove a package   |

## Actions

| ACTION | ARITY | ARGUMENTS    | REMARK            |
|--------|-------|--------------|-------------------|
| NOP    | 0     |              | do nothing        |
| SYSCTL | 2     | COMMAND UNIT | execute systemctl |
