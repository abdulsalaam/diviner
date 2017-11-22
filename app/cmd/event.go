package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	pbs "diviner/protos/service"
)

func eventInvoke(fcn string, cmd *cobra.Command) {
	switch fcn {
	case "query":
		if eventId == "" {
			log.Fatalf("event id error: %q\n", eventId)
		}

		req, err := pbs.NewQueryRequest(priv, eventId)
		if err != nil {
			log.Fatalf("genreate query request error: %v\n", err)
		}

		resp, err := client.QueryMarket(context.Background(), req)
		log.Println("event query", resp, err)
	case "create":
		if title == "" {
			log.Fatal("title is empty")
		}

		if len(outcomes) < 2 {
			log.Fatalf("length of outcomes error: %d", len(outcomes))
		}

		req, err := pbs.NewEventCreateRequest(priv, member.Id, title, outcomes)
		if err != nil {
			log.Fatalf("genreate event create request error: %v\n", err)
		}

		resp, err := client.CreateEvent(context.Background(), req)
		log.Println("market create", resp, err)
	} // switch fcn
}

func NewEventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "event [query|create]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fcn := args[0]
			if fcn != "query" && fcn != "create" {
				log.Fatalf("command error: %s\n")
			}

			eventInvoke(fcn, cmd)
		},
	}

	cmd.Flags().StringVar(&eventId, "id", "", "event id")
	cmd.Flags().StringVar(&title, "title", "", "event title")
	cmd.Flags().StringSliceVar(&outcomes, "outcome", []string{}, "outcomes for event")

	return cmd
}