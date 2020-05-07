# Web-Service & CLI & JSON enforcement tool

## Policy File Restriction

1. JSON format. Be subset of Spec-defined object's attributes.

2. BlackList Rule.

3. DO NOT use "[]", "(B)", "(N)" as mask symbol for policy file attribute value.

## Getting Started

1. Goto /Mask, run `build.sh`, make jm executable.

2. Goto /Server, run `build.sh`, make server executable and a copy of Client config is pushed to Client folder automatically.

3. Goto /Client, run `build.sh`, make client executable.

   IMPORTANT: Building Server must be prior to Building Client.

## How To Use

1. Browse "http(s)://ServerIP:Port/" to get info when Server is running. 

2. Mask & Client executable should be fetched by `wget`.
   
   e.g. `wget IP:Port/mask-linux64` to get Mask executable on Linux64.

   `wget -O config.toml IP:Port/mask-config` to get its configure.

   e.g. `wget IP:Port/client-linux64` to get Client executable on Linux64.

   `wget -O config.toml IP:Port/client-config` to get its configure.

3. For `/clean.sh`, if `rmdb` as first argument is provided, policy-database will be removed!


