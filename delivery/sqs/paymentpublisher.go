package sqs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pevin/pevin-golang-training-beginner/model"
)

type IPaymentPublisher interface {
	Publish(model.Payment) (err error)
}

type PaymentPublisher struct {
	Sqs      *sqs.SQS
	QueueUrl *string
}

func NewPaymentPublisher(sess *session.Session, queue *string) (p PaymentPublisher, err error) {
	sqsClient := sqs.New(sess)

	res, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})

	if err != nil {
		return
	}

	p = PaymentPublisher{Sqs: sqsClient, QueueUrl: res.QueueUrl}

	return

}

func (p PaymentPublisher) Publish(payment model.Payment) (err error) {
	messageBody, _ := json.Marshal(payment)

	_, err = p.Sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
		QueueUrl:    p.QueueUrl,
	})

	return
}
