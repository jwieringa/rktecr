package main

import (
	"encoding/base64"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"os"
	"strings"
)

type DockerAuth struct {
	Credentials `json:"credentials"`
	Registeries []string `json:"registries"`
	RktKind     string   `json:"rktKind"`
	RktVersion  string   `json:"rktVersion"`
}

type Credentials struct {
	Password string `json:"password"`
	User     string `json:"user"`
}

func NewDockerAuth(authData *ecr.AuthorizationData) *DockerAuth {
	dockerAuth := &DockerAuth{
		Credentials: *ExtractCredentials(authData),
		Registeries: []string{
			strings.TrimPrefix(
				aws.StringValue(authData.ProxyEndpoint), "https://",
			),
		},
		RktKind:    "dockerAuth",
		RktVersion: "v1",
	}
	return dockerAuth
}

func ExtractCredentials(authData *ecr.AuthorizationData) *Credentials {
	decodedToken, err := base64.StdEncoding.DecodeString(
		aws.StringValue(authData.AuthorizationToken),
	)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)

	creds := &Credentials{
		User:     parts[0],
		Password: parts[1],
	}
	return creds
}

func FectchAuthorization(registery_id string) *ecr.AuthorizationData {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("error: ", err)
	}

	svc := ecr.New(sess)
	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{
			aws.String(registery_id),
		},
	}
	resp, err := svc.GetAuthorizationToken(params)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return resp.AuthorizationData[0]
}

func main() {
	authData := FectchAuthorization(
		os.Getenv("AWS_REGISTERY_ID"),
	)
	if err := json.NewEncoder(os.Stdout).Encode(
		NewDockerAuth(authData),
	); err != nil {
		log.Fatal(err)
	}
}
