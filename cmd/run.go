package cmd

import (
	"fmt"
	"net/http"

	"github.com/fopina/traefik-cn-foward-auth/pkg/validator"
	"github.com/spf13/cobra"
)

type runOptions struct {
	bindAddress          string
	headerName           string
	allowHeaderName      string
	allowHeaderSeparator string
	debug                bool
	cmd                  *cobra.Command
}

func defaultRunOptions() *runOptions {
	return &runOptions{bindAddress: ":8080", headerName: "X-Forwarded-Tls-Client-Cert-Info", allowHeaderName: "X-Allow-CN", allowHeaderSeparator: ","}
}

func newRootCmd(version string) *cobra.Command {
	o := defaultRunOptions()

	cmd := &cobra.Command{
		Use:          "traefik-cn-foward-auth",
		Short:        "Run HTTP service that validates headers",
		SilenceUsage: true,
		RunE:         o.run,
		Version:      version,
	}

	cmd.Flags().StringVarP(&o.bindAddress, "bind-addr", "b", o.bindAddress, `Address to bind the web server to`)
	cmd.Flags().StringVarP(&o.headerName, "header", "n", o.headerName, `Name of the header with values to be validated`)
	cmd.Flags().StringVarP(&o.allowHeaderName, "allow-header", "a", o.allowHeaderName, `Name of the header that will container a list of allowed values`)
	cmd.Flags().StringVarP(&o.allowHeaderSeparator, "allow-header-separator", "s", o.allowHeaderSeparator, `Separator character of the values in allow-header. Use special value "json" if you prefer to use a JSON string array to specify the list`)
	cmd.Flags().BoolVar(&o.debug, "debug", true, `Log all failed validations to help debug`)

	return cmd
}

func (o *runOptions) run(cmd *cobra.Command, args []string) error {
	o.cmd = cmd
	http.HandleFunc("/", o.handler)
	cmd.Printf("Starting server on %s\n", o.bindAddress)
	if err := http.ListenAndServe(o.bindAddress, nil); err != nil {
		return err
	}
	return nil
}

func (o *runOptions) handler(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get(o.headerName)
	allowed := r.Header.Get(o.allowHeaderName)
	// forbidden by default (if either of the headers is missing)
	if (value != "") && (allowed != "") {
		if validator.ValidateValue(value, allowed, o.allowHeaderSeparator) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
			return
		}
	}
	if o.debug {
		o.cmd.Printf("REJECT: Value - %s / Allowed - %s\n", value, allowed)
	}
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "Forbidden")
}

// Execute invokes the command.
func Execute(version string) error {
	return newRootCmd(version).Execute()
}
