# Web-Service & CLI & JSON enforcement tool

## Policy File Restriction

1. JSON format. Be subset of Spec-defined object's attributes.

2. BlackList Rule.

3. DO NOT use "[]", "(B)", "(N)" as mask symbol for policy file attribute value.

## Getting Started

1. One-Step. Run `build.sh`, automatically generate all executables.

   [jm] in ./Mask/build/your-os/
   
   [server] in ./Server/build/your-os/
   
   [client] in ./Client/build/your-os/

   IMPORTANT: If manually build each sub-project, building Server must be prior to building Client.

## How To Use

1. Browse "http(s)://ServerIP:Port/" to get info when Server is running. 

2. Mask & Client executable should be fetched by `wget`.
   
   e.g. `wget IP:Port/mask-linux64` to get Mask executable on Linux64.

   `wget -O config.toml IP:Port/mask-config` to get its configure.

   e.g. `wget IP:Port/client-linux64` to get Client executable on Linux64.

   `wget -O config.toml IP:Port/client-config` to get its configure.

3. For `/clean.sh`, if `rmdb` provided as first argument, policy-database would be removed !!


