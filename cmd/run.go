/*
Copyright © 2022 Simon Emms <simon@simonemms.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
	"github.com/spf13/cobra"
)

var runOpts struct {
	Username    string
	Password    string
	GatewayURL  string
	Topic       string
	AsyncInvoke bool
	ContentType string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if runOpts.Topic == "" {
			return errors.New("topic not set")
		}

		creds := &auth.BasicAuthCredentials{
			User:     runOpts.Username,
			Password: runOpts.Password,
		}

		config := &types.ControllerConfig{
			RebuildInterval:          time.Millisecond * 1000,
			GatewayURL:               runOpts.GatewayURL,
			PrintResponse:            true,
			PrintResponseBody:        true,
			TopicAnnotationDelimiter: ",",
			AsyncFunctionInvocation:  runOpts.AsyncInvoke,
		}

		controller := types.NewController(creds, config)

		receiver := ResponseReceiver{}
		controller.Subscribe(&receiver)

		controller.BeginMapBuilder()

		additionalHeaders := http.Header{}
		additionalHeaders.Add("X-Served-By", "openfaas-pushbullet-connector")

		for {
			log.Printf("Invoking on topic %s - %s\n", runOpts.Topic, runOpts.GatewayURL)
			time.Sleep(2 * time.Second)
			data := []byte("test " + time.Now().String())
			controller.Invoke(runOpts.Topic, &data, additionalHeaders)
		}
	},
}

// ResponseReceiver enables connector to receive results from the
// function invocation
type ResponseReceiver struct{}

// Response is triggered by the controller when a message is
// received from the function invocation
func (ResponseReceiver) Response(res types.InvokerResponse) {
	if res.Error != nil {
		log.Printf("tester got error: %s", res.Error.Error())
	} else {
		log.Printf("tester got result: [%d] %s => %s (%d) bytes", res.Status, res.Topic, res.Function, len(*res.Body))
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&runOpts.Username, "username", "u", getEnvvar("OPENFAAS_USER", "admin"), "OpenFaaS username")
	runCmd.Flags().StringVarP(&runOpts.Password, "password", "p", getEnvvar("OPENFAAS_PASSWORD", ""), "OpenFaaS password")
	runCmd.Flags().StringVarP(&runOpts.GatewayURL, "gateway", "g", getEnvvar("OPENFAAS_GATEWAY", "http://127.0.0.1:8080"), "Gateway URL")
	runCmd.Flags().StringVarP(&runOpts.Topic, "topic", "t", getEnvvar("OPENFAAS_TOPIC", ""), "The topic name to/from which to publish/subscribe - this matches the annotation 'topic: <topic>' on the function")
	runCmd.Flags().BoolVar(&runOpts.AsyncInvoke, "async-invoke", false, "Invoke via queueing using NATS and the function's async endpoint")
	runCmd.Flags().StringVar(&runOpts.ContentType, "content-type", "application/json", "Response content type")
}
