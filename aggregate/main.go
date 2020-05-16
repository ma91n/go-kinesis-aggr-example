package main

import (
	"fmt"
	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"golang.org/x/sync/errgroup"
	"os"
)

var kc = kinesis.New(session.Must(session.NewSession()))

func handle(e events.KinesisEvent) error {
	fmt.Println("【Start Aggregation Lambda】", len(e.Records))

	var pr = producer.New(&producer.Config{
		StreamName: os.Getenv("KINESIS_STREAM"),
		Client:     kc,
	})

	eg := errgroup.Group{}

	pr.Start() // Producer用のgoroutine起動
	eg.Go(func() error {
		for r := range pr.NotifyFailures() {
			return r
		}
		return nil
	})

	for _, r := range e.Records {
		// TODO 取得したレコードに対する何かしらの処理。ここでは単純に集約して終わり
		if err := pr.Put(r.Kinesis.Data, r.Kinesis.PartitionKey); err != nil {
			return err
		}
	}
	pr.Stop() // 送信中のレコードのflushと、Producer goroutineの停止

	return eg.Wait()
}

func main() {
	lambda.Start(handle)
}
