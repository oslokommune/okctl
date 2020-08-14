module github.com/oslokommune/okctl

go 1.14

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible // indirect

require (
	github.com/99designs/keyring v1.1.5
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/aws-sdk-go v1.34.4
	github.com/awslabs/goformation/v4 v4.13.1
	github.com/containerd/containerd v1.3.4
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/foolin/pagser v0.1.5
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-kit/kit v0.10.0
	github.com/golangci/golangci-lint v1.30.0 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.2.2
	github.com/gofrs/flock v0.7.1
	github.com/google/go-cmp v0.5.1
	github.com/google/uuid v1.1.1
	github.com/jarcoal/httpmock v1.0.6
	github.com/kr/pty v1.1.8 // indirect
	github.com/mishudark/errors v0.0.0-20190221111348-b16f7e94bb58
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.0 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rancher/k3d v1.7.0 // indirect
	github.com/rancher/k3d/v3 v3.0.0
	github.com/sanity-io/litter v1.3.0
	github.com/sebdah/goldie/v2 v2.5.1
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.3.4
	github.com/spf13/cobra v1.0.1-0.20200629195214-2c5a0d300f8b
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1 // indirect
	gopkg.in/ini.v1 v1.58.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.3.0
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/cli-runtime v0.18.4
	k8s.io/client-go v0.18.4
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)
