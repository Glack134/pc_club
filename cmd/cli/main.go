package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gameadmin",
		Short: "Gaming Access Control CLI",
	}

	grantCmd := &cobra.Command{
		Use:   "grant [user] [pc] [minutes]",
		Short: "Grant access to PC",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Granting %s access to %s for %s minutes\n", args[0], args[1], args[2])
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke <pc_id>",
		Short: "Revoke access from PC",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Revoking access from PC %s\n", args[0])
		},
	}
	rootCmd.AddCommand(grantCmd, revokeCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
