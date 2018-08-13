package mongodbatlas

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/dghubble/sling"
)

const apiURL = "https://cloud.mongodb.com/api/atlas/v1.0/"

// Client is a MongoDB Atlas client for making MongoDB API requests.
type Client struct {
	sling               *sling.Sling
	Root                *RootService
	Whitelist           *WhitelistService
	Projects            *ProjectService
	Clusters            *ClusterService
	Containers          *ContainerService
	Peers               *PeerService
	DatabaseUsers       *DatabaseUserService
	Organizations       *OrganizationService
	AlertConfigurations *AlertConfigurationService
}

// SlingLogger implements Sling's Doer interface so that we can log request
// and response bodies
type SlingLogger struct {
	Client *http.Client
}

// Do implements Sling's Doer interface so that we can log request
// and response bodies
func (slingLogger *SlingLogger) Do(req *http.Request) (*http.Response, error) {
	out, err := os.Create("./logFile.log")
	logger := log.New(out, "", log.LstdFlags)
	logger.Println("============================================================")

	requestDump, err := httputil.DumpRequestOut(req, true)
	logger.Println("--> Request:")
	if err != nil {
		logger.Println(err)
	} else {
		logger.Println(string(requestDump))
	}

	resp, err := slingLogger.Client.Do(req)
	if err != nil {
		logger.Println(err)
	}

	defer resp.Body.Close()

	responseDump, err := httputil.DumpResponse(resp, true)
	logger.Println("--> Response:")
	if err != nil {
		logger.Println(err)
	} else {
		logger.Println(string(responseDump))
	}
	logger.Printf("\n\n")

	return resp, err
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(apiURL)
	if os.Getenv("DEBUG_SLING") != "" {
		base = base.Doer(&SlingLogger{Client: httpClient})
	}
	return &Client{
		sling:               base,
		Root:                newRootService(base.New()),
		Whitelist:           newWhitelistService(base.New()),
		Projects:            newProjectService(base.New()),
		Clusters:            newClusterService(base.New()),
		Containers:          newContainerService(base.New()),
		Peers:               newPeerService(base.New()),
		DatabaseUsers:       newDatabaseUserService(base.New()),
		Organizations:       newOrganizationService(base.New()),
		AlertConfigurations: newAlertConfigurationService(base.New()),
	}
}
