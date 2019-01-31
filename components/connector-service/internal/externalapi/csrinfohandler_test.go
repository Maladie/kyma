package externalapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyma-project/kyma/components/connector-service/internal/httperrors"

	"github.com/kyma-project/kyma/components/connector-service/internal/externalapi/mocks"

	"github.com/kyma-project/kyma/components/connector-service/internal/apperrors"
	"github.com/kyma-project/kyma/components/connector-service/internal/certificates"
	"github.com/kyma-project/kyma/components/connector-service/internal/clientcontext"
	tokenMocks "github.com/kyma-project/kyma/components/connector-service/internal/tokens/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	commonName     = "commonName"
	application    = "application"
	certificateURL = "https://host/v1/applications/certificates"
)

type dummyClientContext struct{}

func (dc dummyClientContext) ToJSON() ([]byte, error) {
	return []byte("test"), nil
}

func (dc dummyClientContext) GetApplication() string {
	return application
}

func (dc dummyClientContext) GetCommonName() string {
	return commonName
}

func TestInfoHandler_GetCSRInfo(t *testing.T) {

	url := fmt.Sprintf("/v1/applications/signingRequests/info?token=%s", token)

	subjectValues := certificates.CSRSubject{
		Country:            country,
		Organization:       organization,
		OrganizationalUnit: organizationalUnit,
		Locality:           locality,
		Province:           province,
	}

	dummyClientContext := dummyClientContext{}
	connectorClientExtractor := func(ctx context.Context) (clientcontext.ConnectorClientContext, apperrors.AppError) {
		return dummyClientContext, nil
	}

	t.Run("should successfully get csr info", func(t *testing.T) {
		// given
		newToken := "newToken"
		expectedSignUrl := fmt.Sprintf("https://%s/v1/applications/certificates?token=%s", host, newToken)

		expectedAPI := "dummyAPI"

		expectedCertInfo := certInfo{
			Subject:      fmt.Sprintf("OU=%s,O=%s,L=%s,ST=%s,C=%s,CN=%s", organizationalUnit, organization, locality, province, country, commonName),
			Extensions:   "",
			KeyAlgorithm: "rsa2048",
		}

		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return(newToken, nil)

		apiURLsGenerator := &mocks.APIUrlsGenerator{}
		apiURLsGenerator.On("Generate", dummyClientContext).Return(expectedAPI)

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(tokenRequestRaw))
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		// when
		infoHandler.GetCSRInfo(rr, req)

		// then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse infoResponse
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedSignUrl, infoResponse.CsrURL)
		assert.EqualValues(t, expectedAPI, infoResponse.API)
		assert.EqualValues(t, expectedCertInfo, infoResponse.CertificateInfo)
	})

	t.Run("should return 500 when failed to extract context", func(t *testing.T) {
		// given
		tokenCreator := &tokenMocks.Creator{}
		apiURLsGenerator := &mocks.APIUrlsGenerator{}

		errorExtractor := func(ctx context.Context) (clientcontext.ConnectorClientContext, apperrors.AppError) {
			return nil, apperrors.Internal("error")
		}

		infoHandler := NewCSRInfoHandler(tokenCreator, errorExtractor, apiURLsGenerator, certificateURL, subjectValues)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(tokenRequestRaw))
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		// when
		infoHandler.GetCSRInfo(rr, req)

		// then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var errorResponse httperrors.ErrorResponse
		err = json.Unmarshal(responseBody, &errorResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should return 500 when failed to replace token", func(t *testing.T) {
		// given
		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return("", apperrors.Internal("error"))

		apiURLsGenerator := &mocks.APIUrlsGenerator{}

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(tokenRequestRaw))
		require.NoError(t, err)
		rr := httptest.NewRecorder()

		// when
		infoHandler.GetCSRInfo(rr, req)

		// then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var errorResponse httperrors.ErrorResponse
		err = json.Unmarshal(responseBody, &errorResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should use headers if present", func(t *testing.T) {
		//given
		newToken := "newToken"
		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return(newToken, nil)

		apiURLsGenerator := NewApplicationApiUrlsStrategy("", "", "", "")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Add(BaseEventsPathHeader, "events.base.path")
		req.Header.Add(BaseMetadataPathHeader, "metadata.base.path")
		expectedMetadataUrl := "https://metadata.base.path/application/v1/metadata/services"
		expectedEventsUrl := "https://events.base.path/application/v1/events"

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetCSRInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse infoResponse
		infoResponse.API = &applicationApi{}
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		api := infoResponse.API.(*applicationApi)
		assert.Equal(t, expectedMetadataUrl, api.MetadataURL)
		assert.Equal(t, expectedEventsUrl, api.EventsURL)
	})

	t.Run("should use default value for empty header", func(t *testing.T) {
		//given
		newToken := "newToken"
		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return(newToken, nil)

		apiURLsGenerator := NewApplicationApiUrlsStrategy("metadata.default.path", "", "", "")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Add(BaseEventsPathHeader, "events.base.path")
		req.Header.Add(BaseMetadataPathHeader, "")
		expectedMetadataUrl := "https://metadata.default.path/application/v1/metadata/services"
		expectedEventsUrl := "https://events.base.path/application/v1/events"

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetCSRInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse infoResponse
		infoResponse.API = &applicationApi{}
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		api := infoResponse.API.(*applicationApi)
		assert.Equal(t, expectedMetadataUrl, api.MetadataURL)
		assert.Equal(t, expectedEventsUrl, api.EventsURL)
	})

	t.Run("should use default value for non-existing header", func(t *testing.T) {
		//given
		newToken := "newToken"
		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return(newToken, nil)

		apiURLsGenerator := NewApplicationApiUrlsStrategy("metadata.default.path", "", "", "")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Add(BaseEventsPathHeader, "events.base.path")
		expectedMetadataUrl := "https://metadata.default.path/application/v1/metadata/services"
		expectedEventsUrl := "https://events.base.path/application/v1/events"

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetCSRInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse infoResponse
		infoResponse.API = &applicationApi{}
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		api := infoResponse.API.(*applicationApi)
		assert.Equal(t, expectedMetadataUrl, api.MetadataURL)
		assert.Equal(t, expectedEventsUrl, api.EventsURL)
	})

	t.Run("should return empty url when both header and default value are empty", func(t *testing.T) {
		//given
		newToken := "newToken"
		tokenCreator := &tokenMocks.Creator{}
		tokenCreator.On("Replace", token, dummyClientContext).Return(newToken, nil)

		apiURLsGenerator := NewApplicationApiUrlsStrategy("", "", "", "")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Add(BaseEventsPathHeader, "events.base.path")
		req.Header.Add(BaseMetadataPathHeader, "")
		expectedMetadataUrl := ""
		expectedEventsUrl := "https://events.base.path/application/v1/events"

		infoHandler := NewCSRInfoHandler(tokenCreator, connectorClientExtractor, apiURLsGenerator, certificateURL, subjectValues)

		rr := httptest.NewRecorder()

		//when
		infoHandler.GetCSRInfo(rr, req)

		//then
		responseBody, err := ioutil.ReadAll(rr.Body)
		require.NoError(t, err)

		var infoResponse infoResponse
		infoResponse.API = &applicationApi{}
		err = json.Unmarshal(responseBody, &infoResponse)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		api := infoResponse.API.(*applicationApi)
		assert.Equal(t, expectedMetadataUrl, api.MetadataURL)
		assert.Equal(t, expectedEventsUrl, api.EventsURL)
	})
}
