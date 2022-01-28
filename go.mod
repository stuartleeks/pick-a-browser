module github.com/stuartleeks/pick-a-browser

go 1.17

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/google/uuid v1.3.0
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/jsonc v0.3.2
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/go-github/v30 v30.1.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tcnksm/go-gitconfig v0.1.2 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	golang.org/x/oauth2 v0.0.0-20181106182150-f42d05182288 // indirect
	google.golang.org/appengine v1.3.0 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/rhysd/go-github-selfupdate v1.2.3 => github.com/stuartleeks/go-github-selfupdate v1.2.4
