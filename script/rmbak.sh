#!/bin/bash

RESERVED=24

pushd /opt/sword/data/ &>/dev/null

total=$(ls -al | awk '{print $9}' | grep -e '^data-' | wc -l)
echo "Total $total backup files found"
count=$(expr ${total} - ${RESERVED})
echo "Prepare to remove $count backup files"
ls -al | awk '{print $9}' | grep -e '^data-' | sort | sed $(expr ${count} + 1)',$d' | xargs rm -f
echo "Finished"

popd &>/dev/null
