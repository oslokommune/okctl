module github.com/oslokommune/okctl

go 1.16

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible // indirect

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

require (
	github.com/99designs/keyring v1.1.6
	github.com/AlecAivazis/survey/v2 v2.3.2
	github.com/Masterminds/semver v1.5.0
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/asdine/storm/v3 v3.2.1
	github.com/aws/aws-sdk-go v1.43.28
	github.com/awslabs/goformation/v4 v4.19.5
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/containerd v1.6.15
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.3+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/evanphx/json-patch/v5 v5.6.0
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-git/go-billy/v5 v5.3.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-kit/kit v0.12.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gofrs/flock v0.8.1
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.12.0
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/jarcoal/httpmock v1.1.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/miekg/dns v1.1.45
	github.com/mishudark/errors v0.0.0-20210318113247-bd4e9ef2fc74
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/oslokommune/okctl-metrics-service v0.1.9
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.52.1
	github.com/rancher/k3d/v3 v3.4.0
	github.com/sanity-io/litter v1.5.2
	github.com/sebdah/goldie/v2 v2.5.3
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/afero v1.8.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/theckman/yacspin v0.13.12
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.7 // indirect
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	gopkg.in/h2non/gock.v1 v1.1.2
	gopkg.in/ini.v1 v1.66.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.7.2
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/cli-runtime v0.22.4
	k8s.io/client-go v0.22.5
	sigs.k8s.io/aws-iam-authenticator v0.5.9
	sigs.k8s.io/yaml v1.3.0
)
