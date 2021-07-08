package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pevin/pevin-golang-training-beginner/cmd"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) >= 2 {
		c := os.Args[1]
		switch c {
		case "rest":
			command := cmd.WebCommand{}
			command.Execute()
			return
		case "cron":
			command := cmd.CronCommand{Ctx: context.Background()}
			command.Execute()
			return
		default:
			fmt.Println("Please enter rest or cron as argument.")
			return
		}
	}
	fmt.Println("Please enter arguments.")
}
