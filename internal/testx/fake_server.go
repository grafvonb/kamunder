package testx

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Predefined responses for collection endpoints
var collectionResponses = map[string]string{
	"/App/System": `{
    }`,
}

// Predefined responses for single-object endpoints
var objectResponses = map[string]string{
	"/v2/topology": `{
	  "Brokers": [
		{
		  "Host": "camunda-platform-c88-zeebe-0.camunda-platform-c88-zeebe",
		  "NodeId": 0,
		  "Partitions": [
			{
			  "Health": "healthy",
			  "PartitionId": 1,
			  "Role": "leader"
			}
		  ],
		  "Port": 26501,
		  "Version": "8.8.0"
		}
	  ],
	  "ClusterSize": 1,
	  "GatewayVersion": "8.8.0",
	  "PartitionsCount": 1,
	  "ReplicationFactor": 1,
	  "LastCompletedChangeId": ""
	}`,
}

var (
	onceFS   sync.Once
	sharedFS *FakeServer
)

type FakeServer struct {
	FS      *httptest.Server
	BaseURL string
}

// NewFakeServer returns a FakeServer with fixed responses
func NewFakeServer(t *testing.T) *FakeServer {
	t.Helper()
	onceFS.Do(func() {
		fs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if resp, ok := collectionResponses[r.URL.Path]; ok {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(resp))
				return
			}
			if resp, ok := objectResponses[r.URL.Path]; ok {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(resp))
				return
			}
			http.NotFound(w, r)
		}))
		sharedFS = &FakeServer{
			FS:      fs,
			BaseURL: fs.URL,
		}
	})

	require.NotNil(t, sharedFS)
	return sharedFS
}
