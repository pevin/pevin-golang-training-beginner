package cmd

import (
	"context"
	"fmt"
)

type CronCommand struct {
	Ctx context.Context
}

func (c *CronCommand) Execute() {
	fmt.Println("running cron...")
	fmt.Println("expiring payment codes which has expiry date less than current date...")

	pcUsecase := initPaymentCodeUsecase()
	err := pcUsecase.ExpireWithPassDueExpiryDate(c.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("expired payment codes done.")
}
