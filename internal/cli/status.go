package cli

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/linkding-cli/internal/linkding"
)

func (a *app) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show configuration and optionally verify API access",
		Long:  `Reports URL, whether a token is configured (never prints it), and allow_write. Calls GET /api/user/profile/ when a token is present.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.loadConfig(); err != nil {
				return err
			}
			out := map[string]any{
				"ok":              true,
				"url":             a.cfg.URL,
				"has_token":       a.cfg.Token != "",
				"token_from":      a.cfg.Source,
				"allow_write":     a.cfg.AllowWrite,
				"allow_write_via": a.cfg.AllowWriteSource,
			}
			if a.cfg.Token == "" || a.cfg.URL == "" {
				return a.emitValue(out)
			}
			c, err := a.client()
			if err != nil {
				return err
			}
			raw, err := c.Do(context.Background(), linkding.Request{
				Method: http.MethodGet,
				Path:   "/api/user/profile/",
			})
			if err != nil {
				return a.wrapErr(err)
			}
			var profile any
			_ = json.Unmarshal(raw, &profile)
			out["profile"] = profile
			return a.emitValue(out)
		},
	}
}

func (a *app) accessCmd() *cobra.Command {
	access := &cobra.Command{
		Use:   "access",
		Short: "View or change the local allow_write policy",
	}
	access.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Print current allow_write state",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.loadConfig(); err != nil {
				return err
			}
			return a.emitValue(map[string]any{
				"ok":              true,
				"allow_write":     a.cfg.AllowWrite,
				"allow_write_via": a.cfg.AllowWriteSource,
			})
		},
	})
	var confirm string
	write := &cobra.Command{
		Use:   "write",
		Short: "Enable allow_write for mutating commands (requires a TTY)",
		Long:  `Stores allow_write in the config directory allow_write file (mode 0600). Agents must not run this; the user must.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !a.isTTY() {
				return policyErr("linkding access write requires an interactive terminal")
			}
			if confirm != "write" {
				return usageErr(`re-run with --confirm write (type the word "write")`)
			}
			if err := linkding.WriteAllowWriteFile(true); err != nil {
				return runtimeErr(err)
			}
			return a.emitValue(map[string]any{"ok": true, "allow_write": true})
		},
	}
	write.Flags().StringVar(&confirm, "confirm", "", `must be the word "write"`)
	access.AddCommand(write)

	read := &cobra.Command{
		Use:   "read",
		Short: "Return to read-only mode (requires a TTY)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !a.isTTY() {
				return policyErr("linkding access read requires an interactive terminal")
			}
			if err := linkding.WriteAllowWriteFile(false); err != nil {
				return runtimeErr(err)
			}
			return a.emitValue(map[string]any{"ok": true, "allow_write": false})
		},
	}
	access.AddCommand(read)
	return access
}
