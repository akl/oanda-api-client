package oanda_api_client

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Param implements parameter store access method.
type Param string

const (
	envPlaceholder = "<ENV>"
	prefix         = "/Oanda/" + envPlaceholder
)

// ParamOandaAPIKey defines api key of oanda API.
var ParamOandaAPIKey = Param(prefix + "APIkey")

var (
	ssmClient     *ssm.SSM
	ssmClientOnce sync.Once
	env           = os.Getenv("ENVIRONMENT")
)

func initClient() {
	s := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	}))
	ssmClient = ssm.New(s)
}

func (p Param) FetchValue() string {
	if len(env) == 0 {
		return ""
	}
	resolved := strings.Replace(string(p), envPlaceholder, env, 1)
	ssmClientOnce.Do(initClient)
	output, err := ssmClient.GetParameter(&ssm.GetParameterInput{Name: aws.String(resolved)})
	if err != nil {
		log.Fatalf("failed to get parameter (key=%s): %v", resolved, err)
		return ""
	}
	if output.Parameter.Value == nil {
		log.Fatalf("failed to get parameter (key=%s): empty", resolved)
		return ""
	}
	return *output.Parameter.Value
}
