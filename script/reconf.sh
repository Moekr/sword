#!/bin/bash

if [ -d /opt/sword/data.bak ]
then
rm -rf /opt/sword/data.bak
fi
cp -R /opt/sword/data /opt/sword/data.bak
systemctl restart sword.server.service
