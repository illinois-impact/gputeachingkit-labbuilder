#!/bin/zsh
rm -fr ./filter
go build 
pandoc --filter ./filter -t html -o test.html ./test.markdown |&/Users/abduld/Code/go/bin/pp
