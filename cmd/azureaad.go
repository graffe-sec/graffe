package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var reviewAadUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Reviews Azure Active Directory user security configuration",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := azureAuth()
		fmt.Println(auth)
		fmt.Println(err)
	},
}

func init() {
	azureReviewCmd.AddCommand(reviewAadUsersCmd)
}
