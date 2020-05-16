package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	deagg "github.com/awslabs/kinesis-aggregation/go/deaggregator"
)

func handle(e events.KinesisEvent) error {
	fmt.Println("【Start DeAggregation Lambda】", len(e.Records))

	krs := make([]*kinesis.Record, 0, len(e.Records))

	for _, r := range e.Records {
		krs = append(krs, &kinesis.Record{
			ApproximateArrivalTimestamp: aws.Time(r.Kinesis.ApproximateArrivalTimestamp.UTC()),
			Data:                        r.Kinesis.Data,
			EncryptionType:              &r.Kinesis.EncryptionType,
			PartitionKey:                &r.Kinesis.EncryptionType,
			SequenceNumber:              &r.Kinesis.SequenceNumber,
		})
	}

	dars, err := deagg.DeaggregateRecords(krs)
	if err != nil {
		return err
	}

	for _, r := range dars {
		// TODO de-aggregation後レコードに対する処理
		fmt.Println("input", string(r.Data))
	}
	return nil
}

func main() {
	lambda.Start(handle)
}
