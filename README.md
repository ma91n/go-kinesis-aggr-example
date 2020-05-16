# go-kinesis-aggr-example
Go Kinesis Aggregation Format Example

## Deploy

AWS Profileは `my_profile` で作成済みとする。デフォルトリージョンも設定済みとする。

### Kinesis

```bash
# Kinesisの作成
aws kinesis --profile=my_profile create-stream --stream-name aggregate --shard-count 1
aws kinesis --profile=my_profile create-stream --stream-name deaggregate --shard-count 1
```

### Lambda

```bash
# Aggregation
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o aggregate/main ./aggregate/main.go
zip -j aggregate/main.zip aggregate/main

## (初回)
aws lambda create-function --profile=my_profile --function-name aggregate-lambda --runtime go1.x --handler main --zip-file fileb://aggregate/main.zip \
    --role <Your Lambda Role ARN> --environment 'Variables={KINESIS_STREAM="deaggregate"}'

## (2回目移行)
aws lambda update-function-code --profile my_profile --region ap-northeast-1 --function-name aggregate-lambda --zip-file fileb://aggregate/main.zip

# DeAggregation
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o deaggregate/main ./deaggregate/main.go
zip -j deaggregate/main.zip deaggregate/main

## (初回)
aws lambda create-function --profile=my_profile --function-name deaggregate-lambda --runtime go1.x --handler main --zip-file fileb://deaggregate/main.zip \
    --role <Your Lambda Role ARN>

## (2回目移行)
aws lambda update-function-code --profile my_profile --region ap-northeast-1 --function-name deaggregate-lambda --zip-file fileb://deaggregate/main.zip
```

### Lambda Event Source Mapping

```bash
aws lambda create-event-source-mapping --profile my_profile --event-source-arn arn:aws:kinesis:ap-northeast-1:<Your_AWS_Account_ID>:stream/aggregate  \
    --function-name  aggregate-lambda --starting-position TRIM_HORIZON

aws lambda create-event-source-mapping --profile my_profile --event-source-arn arn:aws:kinesis:ap-northeast-1:<Your_AWS_Account_ID>:stream/deaggregate  \
    --function-name deaggregate-lambda --starting-position TRIM_HORIZON
```

## Test

* aws kinesis --profile my_profile put-record --stream-name <steram_name> --partition-key xxx --data <data>

```bash
# 例
aws kinesis --profile my_profile put-record --stream-name aggregate --partition-key 123 --data MTIzNDU2Nzg5MA==
aws kinesis --profile my_profile put-record --stream-name aggregate --partition-key 124 --data MTIzNDU2Nzg5MB==
aws kinesis --profile my_profile put-record --stream-name aggregate --partition-key 125 --data MTIzNDU2Nzg5MC==
```
