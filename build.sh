 #!/bin/bash

 set -e
 ORIGINALPATH=`pwd`

 cd Server && ./build.sh && cd $ORIGINALPATH && echo "Server built" 
 cd Enforcer && ./build.sh && cd $ORIGINALPATH && echo "Enforcer built"
 cd Client && ./build.sh && cd $ORIGINALPATH && echo "Client built"
