package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/linkding-cli/internal/linkding"
)

func (a *app) bookmarksCmd() *cobra.Command {
	bm := &cobra.Command{
		Use:   "bookmarks",
		Short: "List, search, and manage bookmarks",
	}
	bm.AddCommand(
		a.bookmarksListCmd(false),
		a.bookmarksListCmd(true),
		a.bookmarksSearchCmd(),
		a.bookmarksGetCmd(),
		a.bookmarksCheckCmd(),
		a.bookmarksCreateCmd(),
		a.bookmarksUpdateCmd(),
		a.bookmarksArchiveCmd(true),
		a.bookmarksArchiveCmd(false),
		a.bookmarksDeleteCmd(),
	)
	return bm
}

func (a *app) bookmarksListCmd(archived bool) *cobra.Command {
	var limit, offset int
	var modifiedSince, addedSince, bundle string
	path := "/api/bookmarks/"
	short := "List bookmarks"
	use := "list"
	if archived {
		path = "/api/bookmarks/archived/"
		short = "List archived bookmarks"
		use = "list-archived"
	}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if offset > 0 {
				q.Set("offset", strconv.Itoa(offset))
			}
			if modifiedSince != "" {
				q.Set("modified_since", modifiedSince)
			}
			if addedSince != "" {
				q.Set("added_since", addedSince)
			}
			if bundle != "" {
				q.Set("bundle", bundle)
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   path,
				Query:  q,
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results (API default 100)")
	cmd.Flags().IntVar(&offset, "offset", 0, "pagination offset")
	cmd.Flags().StringVar(&modifiedSince, "modified-since", "", "ISO 8601 timestamp")
	cmd.Flags().StringVar(&addedSince, "added-since", "", "ISO 8601 timestamp")
	cmd.Flags().StringVar(&bundle, "bundle", "", "bundle id filter")
	return cmd
}

func (a *app) bookmarksSearchCmd() *cobra.Command {
	var limit, offset int
	var qStr string
	cmd := &cobra.Command{
		Use:   "search QUERY",
		Short: "Search bookmarks (same logic as the web UI)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			phrase := args[0]
			if qStr != "" {
				phrase = qStr
			}
			q := url.Values{}
			q.Set("q", phrase)
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if offset > 0 {
				q.Set("offset", strconv.Itoa(offset))
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/bookmarks/",
				Query:  q,
			})
		},
	}
	cmd.Flags().StringVar(&qStr, "q", "", "search phrase (alternative to positional QUERY)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	cmd.Flags().IntVar(&offset, "offset", 0, "pagination offset")
	return cmd
}

func (a *app) bookmarksGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Retrieve one bookmark by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/bookmarks/" + id + "/",
			})
		},
	}
}

func (a *app) bookmarksCheckCmd() *cobra.Command {
	var rawURL string
	cmd := &cobra.Command{
		Use:   "check URL",
		Short: "Check whether a URL is already bookmarked",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := rawURL
			if len(args) == 1 {
				u = args[0]
			}
			if u == "" {
				return usageErr("provide URL as argument or --url")
			}
			q := url.Values{}
			q.Set("url", u)
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/bookmarks/check/",
				Query:  q,
			})
		},
	}
	cmd.Flags().StringVar(&rawURL, "url", "", "URL to check")
	return cmd
}

func (a *app) bookmarksCreateCmd() *cobra.Command {
	var bodyFile, bodyJSON string
	var disableScraping bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a bookmark (or update if URL already exists)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.requireWrite("bookmarks create"); err != nil {
				return err
			}
			body, err := a.readBody(bodyFile, bodyJSON)
			if err != nil {
				return err
			}
			if len(body) == 0 {
				return usageErr("provide --json or --body with bookmark fields (at least url)")
			}
			q := url.Values{}
			if disableScraping {
				q.Set("disable_scraping", "true")
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodPost,
				Path:   "/api/bookmarks/",
				Query:  q,
				Body:   body,
			})
		},
	}
	cmd.Flags().StringVar(&bodyJSON, "json", "", "bookmark JSON object")
	cmd.Flags().StringVar(&bodyFile, "body", "", "path to JSON file (- for stdin)")
	cmd.Flags().BoolVar(&disableScraping, "disable-scraping", false, "disable automatic title/description scraping")
	return cmd
}

func (a *app) bookmarksUpdateCmd() *cobra.Command {
	var bodyFile, bodyJSON string
	var patch bool
	cmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update a bookmark (PUT or PATCH with --patch)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.requireWrite("bookmarks update"); err != nil {
				return err
			}
			body, err := a.readBody(bodyFile, bodyJSON)
			if err != nil {
				return err
			}
			if len(body) == 0 {
				return usageErr("provide --json or --body with fields to update")
			}
			method := http.MethodPut
			if patch {
				method = http.MethodPatch
			}
			return a.do(context.Background(), linkding.Request{
				Method: method,
				Path:   "/api/bookmarks/" + args[0] + "/",
				Body:   body,
			})
		},
	}
	cmd.Flags().StringVar(&bodyJSON, "json", "", "bookmark JSON object")
	cmd.Flags().StringVar(&bodyFile, "body", "", "path to JSON file (- for stdin)")
	cmd.Flags().BoolVar(&patch, "patch", false, "use PATCH instead of PUT")
	return cmd
}

func (a *app) bookmarksArchiveCmd(archive bool) *cobra.Command {
	use, pathSuffix, action := "archive", "archive/", "bookmarks archive"
	if !archive {
		use, pathSuffix, action = "unarchive", "unarchive/", "bookmarks unarchive"
	}
	return &cobra.Command{
		Use:   use + " ID",
		Short: use + " a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.requireWrite(action); err != nil {
				return err
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/api/bookmarks/%s/%s", args[0], pathSuffix),
			})
		},
	}
}

func (a *app) bookmarksDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.requireWrite("bookmarks delete"); err != nil {
				return err
			}
			return a.do(context.Background(), linkding.Request{
				Method: http.MethodDelete,
				Path:   "/api/bookmarks/" + args[0] + "/",
			})
		},
	}
}

func (a *app) readBody(filePath, jsonInline string) ([]byte, error) {
	if jsonInline != "" {
		if !json.Valid([]byte(jsonInline)) {
			return nil, usageErr("invalid --json")
		}
		return []byte(jsonInline), nil
	}
	if filePath == "" {
		return nil, nil
	}
	var b []byte
	var err error
	if filePath == "-" {
		b, err = readAllLimit(a.in, 4<<20)
	} else {
		b, err = os.ReadFile(filePath)
	}
	if err != nil {
		return nil, runtimeErr(err)
	}
	if !json.Valid(b) {
		return nil, usageErr("body is not valid JSON")
	}
	return b, nil
}

func readAllLimit(r io.Reader, max int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > max {
		return nil, fmt.Errorf("body too large")
	}
	return data, nil
}
