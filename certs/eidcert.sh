 #!/bin/bash

wget --quiet http://sk.ee/upload/files/JUUR-SK.PEM.cer
wget --quiet http://sk.ee/upload/files/EECCRCA.pem.cer
wget --quiet http://sk.ee/upload/files/ESTEID-SK%202007.PEM.cer
wget --quiet http://sk.ee/upload/files/ESTEID-SK%202011.pem.cer

cat JUUR-SK.PEM.cer EECCRCA.pem.cer 'ESTEID-SK 2007.PEM.cer' 'ESTEID-SK 2011.pem.cer' > id.crt