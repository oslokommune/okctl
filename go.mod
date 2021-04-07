module github.com/oslokommune/okctl

go 1.16

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible // indirect

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

require (
	github.com/99designs/keyring v1.1.6
	github.com/AlecAivazis/survey/v2 v2.2.9
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/Shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d // indirect
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/aws-sdk-go v1.38.9
	github.com/awslabs/goformation/v4 v4.18.3
	github.com/beevik/etree v1.1.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/bshuster-repo/logrus-logstash-hook v1.0.0 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/containerd/cgroups v0.0.0-20210114181951-8a68de567b68 // indirect
	github.com/containerd/containerd v1.4.4
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/containerd/fifo v0.0.0-20210129194248-f8e8fdba47ef // indirect
	github.com/containerd/ttrpc v1.0.2 // indirect
	github.com/containerd/typeurl v1.0.1 // indirect
	github.com/docker/cli v20.10.3+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.3+incompatible // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/evanphx/json-patch/v5 v5.2.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/foolin/pagser v0.1.5
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-git/go-billy/v5 v5.1.0
	github.com/go-git/go-git/v5 v5.3.0
	github.com/go-kit/kit v0.10.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/gofrs/flock v0.8.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gosimple/slug v1.9.0
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jarcoal/httpmock v1.0.8
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miekg/dns v1.1.41
	github.com/mishudark/errors v0.0.0-20210318113247-bd4e9ef2fc74
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/opencontainers/selinux v1.8.0 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/oslokommune/kaex v0.1.7
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.9.1
	github.com/rancher/k3d/v3 v3.4.0
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/sanity-io/litter v1.5.0
	github.com/sebdah/goldie/v2 v2.5.3
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/theckman/yacspin v0.8.0
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.7 // indirect
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	go.opencensus.io v0.22.6 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460 // indirect
	google.golang.org/grpc v1.35.0 // indirect
	gopkg.in/h2non/gock.v1 v1.0.16
	gopkg.in/ini.v1 v1.62.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/cli-runtime v0.20.5
	k8s.io/client-go v0.20.5
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/aws-iam-authenticator v0.5.2
	sigs.k8s.io/yaml v1.2.0
)
