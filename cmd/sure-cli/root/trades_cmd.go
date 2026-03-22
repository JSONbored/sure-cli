package root

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/we-promise/sure-cli/internal/api"
	"github.com/we-promise/sure-cli/internal/output"
)

func newTradesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "trades", Short: "Trades"}

	var from, to string
	var account, symbol string
	var page, perPage int
	var limit int

	list := &cobra.Command{
		Use:   "list",
		Short: "List trades",
		Run: func(cmd *cobra.Command, args []string) {
			client := api.New()

			q := url.Values{}
			if from != "" {
				q.Set("from", from)
			}
			if to != "" {
				q.Set("to", to)
			}
			if account != "" {
				q.Set("account", account)
			}
			if symbol != "" {
				q.Set("symbol", symbol)
			}
			if page > 0 {
				q.Set("page", fmt.Sprintf("%d", page))
			}
			if perPage > 0 {
				q.Set("per_page", fmt.Sprintf("%d", perPage))
			}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}

			u := url.URL{Path: "/api/v1/trades", RawQuery: q.Encode()}
			path := u.String()

			var res any
			r, err := client.Get(path, &res)
			if err != nil {
				output.Fail("request_failed", err.Error(), nil)
			}
			if err := output.Print(format, output.Envelope{Data: res, Meta: &output.Meta{Status: r.StatusCode()}}); err != nil {
				output.Fail("output_failed", err.Error(), nil)
			}
		},
	}

	list.Flags().StringVar(&from, "from", "", "start date (YYYY-MM-DD)")
	list.Flags().StringVar(&to, "to", "", "end date (YYYY-MM-DD)")
	list.Flags().StringVar(&account, "account", "", "account id")
	list.Flags().StringVar(&symbol, "symbol", "", "symbol/ticker")
	list.Flags().IntVar(&page, "page", 1, "page number")
	list.Flags().IntVar(&perPage, "per-page", 25, "items per page (maps to per_page)")
	list.Flags().IntVar(&limit, "limit", 50, "max results")
	cmd.AddCommand(list)

	cmd.AddCommand(&cobra.Command{
		Use:   "show <id>",
		Short: "Show trade",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := api.New()
			var res any
			path := fmt.Sprintf("/api/v1/trades/%s", args[0])
			r, err := client.Get(path, &res)
			if err != nil {
				output.Fail("request_failed", err.Error(), nil)
			}
			if err := output.Print(format, output.Envelope{Data: res, Meta: &output.Meta{Status: r.StatusCode()}}); err != nil {
				output.Fail("output_failed", err.Error(), nil)
			}
		},
	})

	return cmd
}
