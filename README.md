# Web-Service & CLI & JSON enforcement tool

## Policy File Restriction

1. JSON format. Be subset of Spec-defined object's attributes.

2. BlackList Rule.

3. DO NOT use "[]", "(B)", "(N)" as mask symbols for policy file's attribute value.

## Getting Started

1. goto /Mask, run 'build.sh', create executable.
   copy executable to where you want to use.
   run it as its usage.

2. goto /Server, run 'build.sh', create executable.
   run executable with 'config.toml' which has Port setting etc. (a copy of config.toml exist in /Server)

3. goto /Client, run 'build.sh', create executable.
   run executable with 'config.toml' which has IP setting etc. (a copy of config.toml exist in /Client)
