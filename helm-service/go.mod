module github.com/keptn/keptn/helm-service

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.13.1-0.20220228110412-d80b7a4ab103
	github.com/keptn/kubernetes-utils v0.13.1-0.20220223120201-02afb43c03d7
	github.com/kinbiko/jsonassert v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.22.7
	k8s.io/apimachinery v0.22.7
	k8s.io/cli-runtime v0.22.7
	k8s.io/client-go v0.22.7
	k8s.io/kubectl v0.22.7
	sigs.k8s.io/yaml v1.3.0
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
