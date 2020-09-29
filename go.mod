module github.com/oslokommune/okctl

go 1.14

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible // indirect

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6

require (
	github.com/99designs/keyring v1.1.5
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Microsoft/hcsshim/test v0.0.0-20200826032352-301c83a30e7c // indirect
	github.com/Netflix/go-expect v0.0.0-20180615182759-c93bf25de8e8
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/Shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d // indirect
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/aws-sdk-go v1.34.10
	github.com/awslabs/goformation/v4 v4.13.1
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/bshuster-repo/logrus-logstash-hook v1.0.0 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/containerd/cgroups v0.0.0-20200824123100-0b889c03f102 // indirect
	github.com/containerd/containerd v1.4.0
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/foolin/pagser v0.1.5
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-git/go-git/v5 v5.1.0
	github.com/go-kit/kit v0.10.0
	github.com/go-ozzo/ozzo-validation/v4 v4.2.2
	github.com/gofrs/flock v0.8.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.2
	github.com/google/go-github/v32 v32.1.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gosimple/slug v1.9.0
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/hinshun/vt10x v0.0.0-20180616224451-1954e6464174
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jarcoal/httpmock v1.0.6
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miekg/dns v1.0.14
	github.com/mishudark/errors v0.0.0-20190221111348-b16f7e94bb58
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/opencontainers/selinux v1.6.0 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.9.1
	github.com/rancher/k3d/v3 v3.0.1
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/sanity-io/litter v1.3.0
	github.com/sebdah/goldie/v2 v2.5.1
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.3.4
	github.com/spf13/cobra v1.0.1-0.20200629195214-2c5a0d300f8b
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/theckman/yacspin v0.8.0
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.7 // indirect
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20200926100807-9d91bd62050c // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/ini.v1 v1.60.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	helm.sh/helm/v3 v3.3.0
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/cli-runtime v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)
