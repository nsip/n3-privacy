# Web-Service & CLI & JSON enforcement tool

## Policy File Restriction

1. JSON format. Be subset of Spec-defined object's attributes.

2. BlackList Rule.

3. DO NOT use "[]", "(B)", "(N)" as mask symbol for policy file attribute value.

## Getting Started

1. Goto /Mask, run 'build.sh', make executable.

2. Goto /Server, run 'build.sh', make executable and meanwhile push a copy of Client config to its folder.

3. Goto /Client, run 'build.sh', make executable.

4. Building Server should be prior to Building Client ! 

## How To Use

1. Browse "http(s)://ServerIP:Port/" to get info when Server is running. 

2. Mask & Client executable should be fetched by "wget".
