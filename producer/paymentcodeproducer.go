package producer

import (
	"github.com/pevin/pevin-golang-training-beginner/model"
)

type IPaymentCodeMessageProducer interface {
	Produce(p *model.PaymentCode) (err error)
}

type PaymentCodeMessageProducer struct{}

func (r PaymentCodeMessageProducer) Produce(p *model.PaymentCode) (err error) {
	// this is a fake message producer
	return
}
