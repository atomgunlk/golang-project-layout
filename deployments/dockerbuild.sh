#!/bin/sh
appcommit=$(git rev-parse --short HEAD)
appversion=$(cat appversion.txt)-$appcommit

aws ecr get-login-password --region ap-southeast-1 --profile jlc | docker login --username AWS --password-stdin 770138902588.dkr.ecr.ap-southeast-1.amazonaws.com
if [ $? -eq 0 ]; then
    echo "Login OK"
else
    echo "Login FAIL"
	return
fi

cd ../
echo "Building $appversion"
docker build --no-cache -t _your_app_:$appversion .
if [ $? -eq 0 ]; then
    echo "Build OK"
else
    echo "Build FAIL"
	return
fi

docker tag _your_app_:$appversion 770138902588.dkr.ecr.ap-southeast-1.amazonaws.com/_your_app_:$appversion
docker tag _your_app_:$appversion 770138902588.dkr.ecr.ap-southeast-1.amazonaws.com/_your_app_:latest
docker push 770138902588.dkr.ecr.ap-southeast-1.amazonaws.com/_your_app_:$appversion
if [ $? -eq 0 ]; then
    echo "Push OK"
else
    echo "push FAIL"
	return
fi
docker push 770138902588.dkr.ecr.ap-southeast-1.amazonaws.com/_your_app_:latest
if [ $? -eq 0 ]; then
    echo "Push OK"
else
    echo "push FAIL"
	return
fi

echo "Build version $appversion Success"
