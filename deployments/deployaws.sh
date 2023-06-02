#!/bin/sh
appcommit=$(git rev-parse --short HEAD)
appversion=$(cat appversion.txt)-$appcommit
echo $appversion
# cd awscdk
# cdk deploy TEST-YOUR-APP-Service \
# 	--profile jlc \
# 	--parameters appVersion=$appversion \

git add -A && git commit -m "chore: release version $appversion"
# git push origin main
