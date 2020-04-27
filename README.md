# Web-Service & CLI & JSON enforcement tool

## Policy File Restriction

1. JSON format. Be subset of Spec-defined object's attributes.

2. BlackList Rule.

3. DO NOT use "[]", "(B)", "(N)" as mask symbols for policy file attribute value.

## Getting Started

1. goto /Mask, run 'build.sh', create executable.

2. goto /Server, run 'build.sh', create executable.
   run executable with 'config.toml' including Port setting etc. (a copy of config.toml exist in /Server)

3. goto /Client, run 'build.sh', create executable.
   run executable with 'config.toml' including IP setting etc. (a copy of config.toml exist in /Client)

## How To Use

1. Browse "http(s)://ServerIP:ServerPort/" to get info when Server is running. 
