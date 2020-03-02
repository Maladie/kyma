package application

const (
	overridesTemplate = `global:
    domainName: {{ .DomainName }}
    eventServiceImage: {{ .EventServiceImage }}
    eventServiceTestsImage: {{ .EventServiceTestsImage }}
    applicationConnectivityValidatorImage: {{ .ApplicationConnectivityValidatorImage }}
    tenant: {{ .Tenant }}
    group: {{ .Group }}`
)

type OverridesData struct {
	DomainName                            string
	EventServiceImage                     string
	EventServiceTestsImage                string
	ApplicationConnectivityValidatorImage string
	Tenant                                string
	Group                                 string
}
