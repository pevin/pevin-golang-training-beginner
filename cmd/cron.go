package cmd

import (
	"context"
	"fmt"
)

type CronCommand struct {
	Ctx context.Context
}

func (c *CronCommand) Execute() {
	// todo: Implement expire payment codes
	fmt.Println("running cron...")
	fmt.Println("expiring payment codes which has expiry date less than current date...")

	pcUsecase := initPaymentUsecase()
	err := pcUsecase.ExpireWithPassDueExpiryDate(c.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("expired payment codes done.")
}