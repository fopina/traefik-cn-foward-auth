package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fopina/traefik-cn-foward-auth/pkg/validator"
	"github.com/spf13/cobra"
)

type runOptions struct {
	BindAddress          string
	HeaderName           string
	AllowHeaderName      string
	AllowHeaderSeparator string
	Debug                bool
	Raw                  bool
	cmd                  *cobra.Command
}

func defaultRunOptions() *runOptions {
	return &runOptions{BindAddress: ":8080", HeaderName: "X-Forwarded-Tls-Client-Cert-Info", AllowHeaderName: "X-Allow-CN", AllowHeaderSeparator: ","}
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

	cmd.Flags().StringVarP(&o.BindAddress, "bind-addr", "b", o.BindAddress, `Address to bind the web server to`)
	cmd.Flags().StringVarP(&o.HeaderName, "header", "n", o.HeaderName, `Name of the header with values to be validated`)
	cmd.Flags().StringVarP(&o.AllowHeaderName, "allow-header", "a", o.AllowHeaderName, `Name of the header that will container a list of allowed common names - or raw values, if --raw`)
	cmd.Flags().StringVarP(&o.AllowHeaderSeparator, "allow-header-separator", "s", o.AllowHeaderSeparator, `Separator character of the values in allow-header. Use special value "json" if you prefer to use a JSON string array to specify the list`)
	cmd.Flags().BoolVar(&o.Debug, "debug", false, `Log all failed validations to help debug`)
	cmd.Flags().BoolVar(&o.Raw, "raw", false, `By default, values in --allow-header are expected to be the common names, eg: for "CN=mobile01,OU=..." it should have "mobile01". Using --raw no such parsing is done and allowed values are expected to match exactly the value sent in --header`)

	return cmd
}

func (o *runOptions) run(cmd *cobra.Command, args []string) error {
	if o.Debug {
		jsonData, err := json.Marshal(o)
		if err != nil {
			cmd.Println("Error encoding options:", err)
		}
		cmd.Printf("Options: %v\n", string(jsonData))
	}
	o.cmd = cmd
	http.HandleFunc("/", o.handler)
	cmd.Printf("Starting server on %s\n", o.BindAddress)
	if err := http.ListenAndServe(o.BindAddress, nil); err != nil {
		return err
	}
	return nil
}

func (o *runOptions) handler(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get(o.HeaderName)
	allowed := r.Header.Get(o.AllowHeaderName)
	// forbidden by default (if either of the headers is missing)
	if (value != "") && (allowed != "") {
		if o.validateValue(value, allowed) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
			return
		}
	}
	if o.Debug {
		o.cmd.Printf("REJECT: Value - %s / Allowed - %s\n", value, allowed)
	}
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "Forbidden")
}

func (o *runOptions) validateValue(value, allowed string) bool {
	if o.Raw {
		return validator.ValidateValue(value, allowed, o.AllowHeaderSeparator)
	}
	return validator.ValidateCommonName(value, allowed, o.AllowHeaderSeparator)
}

// Execute invokes the command.
func Execute(version string) error {
	return newRootCmd(version).Execute()
}
