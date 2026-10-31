package cli

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/linkding-cli/internal/linkding"
)

func (a *app) bundlesCmd() *cobra.Command {
	bundles := &cobra.Command{
		Use:   "bundles",
		Short: "List bookmark bundles",
	}
	bundles.AddCommand(a.bundlesListCmd())
	return bundles
}

func (a *app) bundlesListCmd() *cobra.Command {
	var limit, offset int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bundles",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if offset > 0 {
				q.Set("offset", strconv.Itoa(offset))
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/bundles/",
				Query:  q,
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	cmd.Flags().IntVar(&offset, "offset", 0, "pagination offset")
	return cmd
}
