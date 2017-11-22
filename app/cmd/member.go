package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	pbs "diviner/protos/service"
)

func memberInvoke(fcn string) {
	switch fcn {
	case "query":
		req, err := pbs.NewQueryRequest(priv, member.Id)
		if err != nil {
			log.Fatalf("genreate query request error: %v\n", err)
		}

		resp, err := client.QueryMember(context.Background(), req)
		log.Println("member query", resp, err)
	case "create":
		req, err := pbs.NewMemberCreateRequest(priv)
		if err != nil {
			log.Fatalf("genreate member create request error: %v\n", err)
		}

		resp, err := client.CreateMember(context.Background(), req)
		log.Println("member create", resp, err)
	} // switch fcn
}

func NewMemberCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "member [query|create]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "query" && args[0] != "create" {
				log.Fatalf("command error: %s\n", args[0])
			}

			memberInvoke(args[0])
		},
	}

	return cmd
}
