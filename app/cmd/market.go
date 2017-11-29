package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	pbm "diviner/protos/member"
	pbs "diviner/protos/service"
)

func maketInvoke(fcn string, cmd *cobra.Command) {
	switch fcn {
	case "query":
		if marketID == "" {
			log.Fatalf("market id error: %q\n", marketID)
		}

		req, err := pbs.NewQueryRequest(priv, marketID)
		if err != nil {
			log.Fatalf("genreate query request error: %v\n", err)
		}

		resp, err := client.QueryMarket(context.Background(), req)
		log.Println("market query", resp, err)
	case "create":
		if eventID == "" {
			log.Fatalf("event id error: %q\n", eventID)
		}

		if number <= 0.0 {
			log.Fatalf("number must be larger than 0: %v", number)
		}

		member, err := pbm.NewMember(priv, 0.0)
		if err != nil {
			log.Fatalf("generate member structure error: %v\n", err)
		}

		req, err := pbs.NewMarketCreateRequest(priv, member.Id, eventID, number, isFund)
		if err != nil {
			log.Fatalf("genreate market create request error: %v\n", err)
		}

		resp, err := client.CreateMarket(context.Background(), req)
		log.Println("market create", resp, err)
	} // switch fcn
}

// NewMarketCmd ...
func NewMarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "market [query|create]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fcn := args[0]
			if fcn != "query" && fcn != "create" {
				log.Fatalf("command error: %s\n", fcn)
			}

			maketInvoke(fcn, cmd)
		},
	}

	cmd.Flags().StringVar(&marketID, "id", "", "market id")
	cmd.Flags().StringVar(&eventID, "event", "", "event id")
	cmd.Flags().Float64Var(&number, "number", 0.0, "number for fund or liquidity")
	cmd.Flags().BoolVar(&isFund, "fund", false, "specify [number] is fund or not")

	return cmd
}
