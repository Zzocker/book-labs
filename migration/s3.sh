# !/bin/bash


echo "[default]
aws_access_key_id = test
aws_secret_access_key = test" | tee credentials

docker exec -w /root localS3 mkdir .aws
docker cp credentials localS3:/root/.aws/credentials

SUCCESS=1
while [ $SUCCESS -eq 1 ];do
docker exec localS3 aws --endpoint-url=http://localhost:4566 s3 mb s3://dev
SUCCESS=$(echo $?)
done

rm credentials