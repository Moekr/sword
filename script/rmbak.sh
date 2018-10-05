#!/bin/bash

RESERVED=24

pushd /opt/sword/data/ &>/dev/null

total=$(ls -al | awk '{print $9}' | grep ^backup | wc -l)
echo "Total $total backup files found"
count=$(expr ${total} - ${RESERVED})
echo "Prepare to remove $count backup files"
ls -al | awk '{print $9}' | grep ^backup | sort | sed $(expr ${count} + 1)',$d' | xargs rm -f
echo "Finished"

popd &>/dev/null
