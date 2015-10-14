#!/bin/bash -e
# usage: ./goverall.sh [func|html]
# default: generate cover.out
output_null=""
output_file=".cover.out"

list_cmd="go list ./..."
if [[ "$2" != "" ]]; then
	list_cmd="echo '$2'"
fi

mkdir -p .cover
eval "$list_cmd" | xargs -I% bash -c 'name="%"; go test % --coverprofile=.cover/${name//\//_}'$output_null
echo "mode: set" > $output_file
if [[ `ls .cover/* 2>/dev/null` == "" ]]; then
	exit
fi
cat .cover/* | grep -v mode >> $output_file
rm -r .cover

if [[ "$1" != "" ]]; then
	go tool cover -$1=$output_file
fi
