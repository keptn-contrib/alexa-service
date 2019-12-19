package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type KeptnEvent struct {
	Stage   string `json:"stage,omitempty"`
	Service string `json:"service,omitempty"`
}

type EvaluationDoneEvent struct {
	KeptnEvent
	Result string `json:"Result"`
}

type ProblemEvent struct {
	State          string `json:"State"`
	ProblemID      string `json:"ProblemID"`
	PID            string `json:"PID"`
	ProblemTitle   string `json:"ProblemTitle"`
	ImpactedEntity string `json:"ImpactedEntity"`
}

//keptnHandler : receives keptn events via http
func keptnHandler(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "alexa-service-go")

	if event.Type() == keptnevents.EvaluationDoneEventType {
		data := &EvaluationDoneEvent{}
		if err := event.DataAs(data); err != nil {
			logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
			return err
		}
		var msg string
		if data.Result == "pass" {
			msg = fmt.Sprintf("New Keptn event detected. EVALUATION DONE. has been reported for %s , in %s."+
			" The result of the evaluation was %s. Promoting artifact to next stage. ", data.Service, data.Result)
		} else {
			msg = fmt.Sprintf("New Keptn event detected. EVALUATION DONE. has been reported for %s , in %s."+
			" The result of the evaluation was %s. The artifact will not be promoted from %s to next stage. ", data.Service, data.Stage, data.Result)
		}
		go postAlexaNotification(msg, logger)
	}
	else if event.Type() == keptnevents.ConfigurationChangeEventType {
		data := &KeptnEvent{}
		if err := event.DataAs(data); err != nil {
			logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
			return err
		}
		logger.Info(fmt.Sprintf("Using AlexaConfig: Service:%s, Stage:%s, Result:%s", data.Service, data.Stage))
		go postAlexaNotification(fmt.Sprintf("New Keptn event detected. CONFIGURATION CHANGED, has been reported for %s , in %s . ", data.Service, data.Stage), logger)
	}
	else if event.Type() == keptnevents.DeploymentFinishedEventType {
		data := &KeptnEvent{}
		if err := event.DataAs(data); err != nil {
			logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
			return err
		}
		logger.Info(fmt.Sprintf("Using AlexaConfig: Service:%s, Stage:%s", data.Service, data.Stage))
		go postAlexaNotification(fmt.Sprintf("New Keptn event detected. DEPLOYMENT FINISHED, has been reported for %s , in %s. ", data.Service, data.Stage), logger)
	}
	else if event.Type() == keptnevents.TestsFinishedEventType {
		data := &KeptnEvent{}
		if err := event.DataAs(data); err != nil {
			logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
			return err
		}
		logger.Info(fmt.Sprintf("Using AlexaConfig: Service:%s, Stage:%s, Result:%s", data.Service, data.Stage))
		go postAlexaNotification(fmt.Sprintf("New Keptn event detected. TESTS FINISHED, has been reported for %s , in %s. ", data.Service, data.Stage), logger)
	}
	else if event.Type() == keptnevents.ProblemOpenEventType {
		data := &ProblemEvent{}
		if err := event.DataAs(data); err != nil {
			logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
			return err
		}
		logger.Info(fmt.Sprintf("Using AlexaConfig: Service:%s, Stage:%s, Result:%s", data.ProblemID, data.ProblemTitle, data.ImpactedEntity))
		var msg string
		if data.State == "OPEN" {
			msg = fmt.Sprintf("New problem reported by Dynatrace. P.I.D. %s . %s . The impact is %s . ", data.ProblemID, data.ProblemTitle, data.ImpactedEntity)
		} else if data.State == "RESOLVED" {
			msg = fmt.Sprintf("Existing problem P.I.D. %s . %s . Has been resolved and the problem is now closed in Dynatrace. ", data.ProblemID, data.ProblemTitle)
		} else {
			msg = fmt.Sprintf("New problem reported by Dynatrace. P.I.D. %s . %s . Has reported a new state. %s", data.ProblemID, data.ProblemTitle, data.State)
		}
		go postAlexaNotification(msg, logger)
	} else {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		logger.Error(fmt.Sprintf("Event type: %s", event.Type()))
		return errors.New(errorMsg)
	}

	return nil
}

func postAlexaNotification(alexaMessage string, logger *keptnutils.Logger) {
	url := os.Getenv("ALEXA_WEBHOOK_URL")
	logger.Info(fmt.Sprintf("URL:>", url))

	var jsonStr = []byte(`{"notification": "` + alexaMessage + `", "accessCode": "` +
		os.Getenv("ALEXA_ACCESS_TOKEN") + `", "title": "CONFIGURATION CHANGED" }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Info(fmt.Sprintf("response Status: %s", resp.Status))
	logger.Info(fmt.Sprintf("response Headers: %s", resp.Header))
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Info(fmt.Sprintf("response Body: %s", string(body)))
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {

	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, keptnHandler))

	return 0
}
