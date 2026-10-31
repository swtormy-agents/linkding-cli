package cli

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/linkding-cli/internal/linkding"
)

func (a *app) userCmd() *cobra.Command {
	user := &cobra.Command{
		Use:   "user",
		Short: "User settings from the API",
	}
	user.AddCommand(&cobra.Command{
		Use:   "profile",
		Short: "GET /api/user/profile/",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/user/profile/",
			})
		},
	})
	return user
}
