 #!/bin/bash

 set -e
 ORIGINALPATH=`pwd`

 cd Server && ./build.sh && cd $ORIGINALPATH && echo "Server built" 
 cd Mask && ./build.sh && cd $ORIGINALPATH && echo "Mask built"
 cd Client && ./build.sh && cd $ORIGINALPATH && echo "Client built"
