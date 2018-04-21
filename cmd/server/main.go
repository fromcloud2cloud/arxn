package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fromcloud2cloud/arxn/arxn/common"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/julienschmidt/httprouter"
)

// Configuration Constants.
const (
	ServerPort     = "SERVER_PORT"
	AwsRegion      = "AWS_REGION"
	AwsRole        = "AWS_ROLE"
	DefaultTimeout = 5
	CloudName      = "CLOUD_NAME"
)

func main() {
	// Server variable declarations.
	var (
		router     *httprouter.Router
		serverPort string
		fmtPort    string
	)
	preFlightCheck()
	serverPort = os.Getenv(ServerPort)
	fmtPort = fmt.Sprintf("0.0.0.0:%s", serverPort)
	router = httprouter.New()
	router.POST("/transfer", handleTransfer)
	http.ListenAndServe(fmtPort, router)
}

func handleTransfer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	cloudName := os.Getenv(CloudName)
	source := r.FormValue("sourceName")
	destination := r.FormValue("destinationName")
	req, err := common.NewTransferRequest(source, destination)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
	}
	ctx, cancel = context.WithCancel(context.Background())
	switch cloudName {
	case common.Aws:
		submitAwsTransfer(ctx, req)
	case common.Gcp:
		submitGcpTransfer(ctx, req)
	case common.Azure:
		submitAzureTransfer(ctx, req)
	}
	defer cancel()
}

// submitAwsTransfer - Send request to SQS.
func submitAwsTransfer(ctx context.Context, req *common.TransferRequest) {

}

// submitGcpTransfer - Send request to PubSub.
func submitGcpTransfer(ctx context.Context, req *common.TransferRequest) {

}

// submitAzureTransfer - Send request to Azure Queue Service.
func submitAzureTransfer(ctx context.Context, req *common.TransferRequest) {

}

// newAwsSession - Get a new AWS Session.
func newAwsSession(assumedRole, region, service string) (*session.Session, *credentials.Credentials) {
	sess := session.Must(session.NewSession())
	if assumedRole != "" {
		creds := stscreds.NewCredentials(sess, assumedRole)
		return sess, creds
	}

	return sess, nil
}

// preFlightCheck - Validate all required env vars are there.
func preFlightCheck() {
	required := []string{ServerPort, AwsRegion, CloudName}
	for _, name := range required {
		if _, isThere := os.LookupEnv(name); !isThere {
			log.Panicf("Required env var %s missing!", name)
		}
	}
}
