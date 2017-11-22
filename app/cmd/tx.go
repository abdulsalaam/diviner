package cmd

import (
	"context"
	"log"
	"strconv"

	pbs "diviner/protos/service"

	"github.com/spf13/cobra"
)

func txInvoke(fcn, share string, volume float64) {
	req, err := pbs.NewTxRequest(priv, member.Id, fcn == "buy", share, volume)
	if err != nil {
		log.Fatalf("generate tx request error: %v\n", err)
	}

	resp, err := client.Tx(context.Background(), req)

	log.Println("tx ", fcn, resp, err)
}

func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "tx [buy|sell]",
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			fcn := args[0]
			share := args[1]
			volume, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				log.Fatalf("volume must be float64: %s, %v", args[2], err)
			}

			if fcn != "buy" && fcn != "sell" {
				log.Fatalf("command error: %s\n")
			}

			txInvoke(fcn, share, volume)
		},
	}

	return cmd
}
