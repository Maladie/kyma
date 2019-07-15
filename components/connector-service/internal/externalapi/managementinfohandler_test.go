package externalapi

import (
	"context"
	"encoding/json"
	"github.com/kyma-project/kyma/components/connector-service/internal/identitymapper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyma-project/kyma/components/connector-service/internal/apperrors"
	"github.com/kyma-project/kyma/components/connector-service/internal/clientcontext"
	"github.com/kyma-project/kyma/components/connector-service/internal/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	applicationKey = "application"
	tenantKey      = "tenant"
	groupKey       = "group"
	runtimeIDKey   = "runtimeID"

	runtimeID = "b45e3a4d-1213-4b7b-a0d0-8fbc580074cf"
	tenant    = "test-tenant"
	group     = "test-group"

	protectedBaseURL      = "https://gateway.kyma.local/v1/applications"
	expectedRenewalsURL   = "https://gateway.kyma.local/v1/applications/certificates/renewals"
	expectedRevocationURL = "https://gateway.kyma.local/v1/applications/certificates/revocations"

	expectedExtensions   = ""
	expectedKeyAlgorithm = "rsa2048"
)

func TestManagementInfoHandler_GetManagementInfo(t *testing.T) {

	t.Run("should successfully get management info response for application", func(t *testing.T) {
		//given
		expectedMetadataURL := "https://metadata.base.path/application/v1/metadata/services"
		expectedEventsURL := "https://events.base.path/application/v1/events"

		applicationContext := clientcontext.ClientContext{
			Tenant: tenant,
			Group:  group,
			ID:     appName,
		}

		extApplicationCtx := &clientcontext.ExtendedApplicationContext{
			ClientContext: applicationContext,
			RuntimeURLs: clientcontext.RuntimeURLs{
				MetadataURL: "https://metadata.base.path/application/v1/metadata/services",
				EventsURL:   "https://events.base.path/application/v1/events",
			},
		}

		connectorClientExtractor := func(ctx context.Context) (clientcontext.ClientCertContextService, apperrors.AppError) {
			return dummyClientCertCtx{extApplicationCtx}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/applications/management/info", nil)
		require.NoError(t, err)

		infoHandler := NewManagementInfoHandler(connectorClientExtractor, protectedBaseURL, identitymapper.ToApplicationIdentity)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetManagementInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse mgmtInfoReponse
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		urls := infoResponse.URLs
		receivedContext := infoResponse.ClientIdentity.(map[string]interface{})
		certificateInfo := infoResponse.CertificateInfo

		assert.Equal(t, expectedMetadataURL, urls.MetadataURL)
		assert.Equal(t, expectedEventsURL, urls.EventsURL)
		assert.Equal(t, expectedRenewalsURL, urls.RenewCertURL)
		assert.Equal(t, expectedRevocationURL, urls.RevokeCertURL)
		assert.Equal(t, appName, receivedContext[applicationKey])
		assert.Equal(t, group, receivedContext[groupKey])
		assert.Equal(t, tenant, receivedContext[tenantKey])
		assert.Equal(t, strSubject, certificateInfo.Subject)
		assert.Equal(t, expectedExtensions, certificateInfo.Extensions)
		assert.Equal(t, expectedKeyAlgorithm, certificateInfo.KeyAlgorithm)
	})

	t.Run("should successfully get management info response for runtime", func(t *testing.T) {
		//given
		runtimeContext := &clientcontext.ClientContext{
			Tenant: tenant,
			Group:  group,
			ID:     runtimeID,
		}

		connectorClientExtractor := func(ctx context.Context) (clientcontext.ClientCertContextService, apperrors.AppError) {
			return dummyClientCertCtx{runtimeContext}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/runtimes/management/info", nil)
		require.NoError(t, err)

		infoHandler := NewManagementInfoHandler(connectorClientExtractor, protectedBaseURL, identitymapper.ToRuntimeIdentity)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetManagementInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse mgmtInfoReponse
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		urls := infoResponse.URLs
		receivedContext := infoResponse.ClientIdentity.(map[string]interface{})
		certificateInfo := infoResponse.CertificateInfo

		assert.Equal(t, expectedRevocationURL, urls.RevokeCertURL)
		assert.Equal(t, expectedRenewalsURL, urls.RenewCertURL)
		assert.Equal(t, runtimeID, receivedContext[runtimeIDKey])
		assert.Equal(t, group, receivedContext[groupKey])
		assert.Equal(t, tenant, receivedContext[tenantKey])
		assert.Equal(t, strSubject, certificateInfo.Subject)
		assert.Equal(t, expectedExtensions, certificateInfo.Extensions)
		assert.Equal(t, expectedKeyAlgorithm, certificateInfo.KeyAlgorithm)
	})

	t.Run("should return 500 when failed to extract context", func(t *testing.T) {
		//given
		clientContextService := func(ctx context.Context) (clientcontext.ClientCertContextService, apperrors.AppError) {
			return nil, apperrors.Internal("error")
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/applications/management/info", nil)
		require.NoError(t, err)

		infoHandler := NewManagementInfoHandler(clientContextService, protectedBaseURL, identitymapper.ToApplicationIdentity)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetManagementInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var errorResposne httperrors.ErrorResponse
		err = json.Unmarshal(responseBody, &errorResposne)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
