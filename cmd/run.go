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
	"log"
	"net/http"
	"time"

	"github.com/MrSimonEmms/openfaas-pushbullet-connector/pkg/pushbullet"
	"github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/faas-provider/auth"
	"github.com/spf13/cobra"
)

var runOpts struct {
	Username        string
	Password        string
	GatewayURL      string
	AsyncInvoke     bool
	ContentType     string
	PushbulletToken string
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Pushbullet listener",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			ContentType:              runOpts.ContentType,
		}

		controller := types.NewController(creds, config)

		receiver := ResponseReceiver{}
		controller.Subscribe(&receiver)

		controller.BeginMapBuilder()

		additionalHeaders := http.Header{}
		additionalHeaders.Add("X-Served-By", "openfaas-pushbullet-connector")

		pb, err := pushbullet.New(runOpts.PushbulletToken)
		if err != nil {
			return err
		}
		defer pb.Close()

		receiveCount := 0
		msgCh := make(chan [3]string)

		pb.SetHandler(func(pb pushbullet.Pushbullet, topic string, push pushbullet.Push) {
			log.Print("Message incoming")
			msgCh <- [3]string{topic, push.Iden, push.Payload()}
		})

		go func() {
			for {
				incoming := <-msgCh

				topic := incoming[0]
				messageId := incoming[1]
				data := []byte(incoming[2])

				// Add a de-dupe header to the message
				additionalHeaders.Add("X-Message-Id", messageId)

				controller.Invoke(topic, &data, additionalHeaders)

				receiveCount++
			}
		}()

		select {}
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
	runCmd.Flags().BoolVar(&runOpts.AsyncInvoke, "async-invoke", false, "Invoke via queueing using NATS and the function's async endpoint")
	runCmd.Flags().StringVar(&runOpts.ContentType, "content-type", "application/json", "Response content type")
	runCmd.Flags().StringVar(&runOpts.PushbulletToken, "pushbullet-token", getEnvvar("PUSHBULLET_TOKEN", ""), "PushBullet token")
}
