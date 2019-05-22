package httphelpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestParseCertHeader(t *testing.T) {

	t.Run("Should return valid CertInfo", func(t *testing.T) {
		//given
		r, _ := http.NewRequest("GET", "", nil)
		r.Header.Set(clientCertHeader, "Hash=f4cf22fb633d4df500e371daf703d4b4d14a0ea9d69cd631f95f9e6ba840f8ad;Subject=\"CN=test-application,OU=OrgUnit,O=Organization,L=Waldorf,ST=Waldorf,C=DE\";URI=spiffe://cluster.local/ns/kyma-integration/sa/default;Hash=6d1f9f3a6ac94ff925841aeb9c15bb3323014e3da2c224ea7697698acf413226;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account")

		//when
		certInfo, e := ParseCertificateHeader(*r, "Organization", "OrgUnit")

		require.NoError(t, e)

		//then
		assert.Equal(t, "f4cf22fb633d4df500e371daf703d4b4d14a0ea9d69cd631f95f9e6ba840f8ad", certInfo.Hash)
		assert.Equal(t, "CN=test-application,OU=OrgUnit,O=Organization,L=Waldorf,ST=Waldorf,C=DE", certInfo.Subject)
		assert.Equal(t, "spiffe://cluster.local/ns/kyma-integration/sa/default", certInfo.URI)
	})

	t.Run("Should return error when unable to find matching CertInfo", func(t *testing.T) {
		//given
		r, _ := http.NewRequest("GET", "", nil)
		r.Header.Set(clientCertHeader, "Hash=f4cf22fb633d4df500e371daf703d4b4d14a0ea9d69cd631f95f9e6ba840f8ad;Subject=\"\";URI=spiffe://cluster.local/ns/kyma-integration/sa/default;Hash=6d1f9f3a6ac94ff925841aeb9c15bb3323014e3da2c224ea7697698acf413226;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account")

		//when
		_, e := ParseCertificateHeader(*r, "Organization", "OrgUnit")

		require.Error(t, e)
	})
}
