package cli

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/linkding-cli/internal/linkding"
)

func (a *app) tagsCmd() *cobra.Command {
	tags := &cobra.Command{
		Use:   "tags",
		Short: "List and manage tags",
	}
	tags.AddCommand(a.tagsListCmd(), a.tagsGetCmd(), a.tagsCreateCmd())
	return tags
}

func (a *app) tagsListCmd() *cobra.Command {
	var limit, offset int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tags",
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
				Path:   "/api/tags/",
				Query:  q,
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	cmd.Flags().IntVar(&offset, "offset", 0, "pagination offset")
	return cmd
}

func (a *app) tagsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Retrieve one tag by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/tags/" + args[0] + "/",
			})
		},
	}
}

func (a *app) tagsCreateCmd() *cobra.Command {
	var bodyFile, bodyJSON, name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.requireWrite("tags create"); err != nil {
				return err
			}
			body, err := a.readBody(bodyFile, bodyJSON)
			if err != nil {
				return err
			}
			if len(body) == 0 {
				if name == "" {
					return usageErr("provide --name or --json {\"name\":\"...\"}")
				}
				b, err := marshal(map[string]string{"name": name})
				if err != nil {
					return runtimeErr(err)
				}
				body = b
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodPost,
				Path:   "/api/tags/",
				Body:   body,
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "tag name")
	cmd.Flags().StringVar(&bodyJSON, "json", "", "tag JSON object")
	cmd.Flags().StringVar(&bodyFile, "body", "", "path to JSON file (- for stdin)")
	return cmd
}
