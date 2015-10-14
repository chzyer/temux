#!/bin/bash
if [[ ! -d '.git' ]]; then
	echo "please call in main directory" >&2
	exit 1
fi
if [[ "$1" == "" ]]; then
	echo "usage: $0 <github.com/xxx/xxx>" >&2
	exit 1
fi
git submodule add https://$1 deps/src/$1
