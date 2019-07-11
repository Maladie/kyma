package middlewares

import (
	"github.com/kyma-project/kyma/components/connector-service/internal/certificates"
	"github.com/kyma-project/kyma/components/connector-service/internal/clientcontext"
	"github.com/kyma-project/kyma/components/connector-service/internal/externalapi/middlewares/runtimeregistry/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRuntimeHealthCheckMiddlewaret(t *testing.T) {
	url := "https://connector-service.kyma.local"

	defaultSubjectVals := certificates.CSRSubject{}
	extractor := clientcontext.NewContextExtractor(defaultSubjectVals)
	contextExtractor := extractor.CreateClusterClientContextService

	runtimeRegistryService := &mocks.RuntimeRegistryService{}

	t.Run("should not send status when certificate CN is default", func(t *testing.T) {
		//given
		runtimeRegistryService.On("ReportState", mock.Anything).Return(nil)

		runtimeHealthCheckMiddleware := NewRuntimeHealthCheckMiddleware(contextExtractor, runtimeRegistryService)

		clusterContext := clientcontext.ClusterContext{
			Group:  "testGroup",
			Tenant: "testTenant",
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		request = request.WithContext(clusterContext.ExtendContext(request.Context()))

		rr := httptest.NewRecorder()

		//when
		resultHandler := runtimeHealthCheckMiddleware.Middleware(handler)
		resultHandler.ServeHTTP(rr, request)

		//then
		assert.Equal(t, http.StatusOK, rr.Code)
		runtimeRegistryService.AssertNumberOfCalls(t, "ReportState", 0)
	})
}
