package sqs

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pevin/pevin-golang-training-beginner/model"
)

// Integration test with AWS SQS
func TestPaymentPublisher_Publish(t *testing.T) {
	queue, region, endpoint := getSqsEnvVars()

	sess, err := session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	})

	if err != nil {
		panic(err)
	}

	type args struct {
		payment model.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				payment: model.Payment{
					Name:          "test-name",
					Amount:        "100",
					TransactionId: "test-transaction-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPaymentPublisher(sess, &queue)

			if err != nil {
				panic(err)
			}

			_, err = p.Sqs.PurgeQueue(&sqs.PurgeQueueInput{
				QueueUrl: p.QueueUrl,
			})

			if err != nil {
				panic(err)
			}

			if err := p.Publish(tt.args.payment); (err != nil) != tt.wantErr {
				t.Errorf("PaymentPublisher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			var timeout int64
			timeout = 5

			msgResult, err := p.Sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				QueueUrl:            p.QueueUrl,
				MaxNumberOfMessages: aws.Int64(10),
				VisibilityTimeout:   &timeout,
			})

			if err != nil {
				panic(err)
			}

			if len(msgResult.Messages) == 0 {
				t.Errorf("PaymentPublisher.Publish() received message was 0")
			}

			_, err = p.Sqs.PurgeQueue(&sqs.PurgeQueueInput{
				QueueUrl: p.QueueUrl,
			})

		})
	}
}

func getSqsEnvVars() (queue string, region string, endpoint string) {
	queue = os.Getenv("SQS_QUEUE_NAME")
	if queue == "" {
		queue = "payment_queue"
	}

	region = os.Getenv("SQS_AWS_REGION")
	if region == "" {
		region = "ap-southeast-2"
	}

	endpoint = os.Getenv("SQS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4566/000000000000/payment_queue"
	}

	return
}
