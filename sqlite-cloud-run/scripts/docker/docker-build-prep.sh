#!/bin/bash

set -e

# https://linuxhint.com/return-string-bash-functions/
# https://stackoverflow.com/questions/3236871/how-to-return-a-string-value-from-a-bash-function

gatherFiles(){
    # Copy files to a tmp directory, since dockerfiles are kind of picky about what it will COPY into the image
    # remember to update the COPY instructions in the Dockerfile as well, if needed
    tmpdir=$(mktemp -d -t 'sqlite-cloud-run-docker-XXXXXX')

    cp scripts/docker/Dockerfile $tmpdir
    cp scripts/docker/run.sh $tmpdir
    cp go.* $tmpdir
    mkdir -p $tmpdir/cmd
    cp -r  cmd/server $tmpdir/cmd/

    echo $tmpdir
}
