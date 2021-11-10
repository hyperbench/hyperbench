module github.com/meshplus/hyperbench

require (
	github.com/Shopify/sarama v1.23.1 // indirect
	github.com/ethereum/go-ethereum v1.10.12
	github.com/fsouza/go-dockerclient v1.4.4 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/gobuffalo/logger v1.0.4 // indirect
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hyperledger/fabric v1.4.3
	github.com/hyperledger/fabric-amcl v0.0.0-20190902191507-f66264322317 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e
	github.com/influxdata/tdigest v0.0.1
	github.com/json-iterator/go v1.1.11
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/meshplus/crypto-standard v0.1.1 // indirect
	github.com/meshplus/gosdk v0.1.1
	github.com/mholt/archiver/v3 v3.5.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pingcap/failpoint v0.0.0-20191029060244-12f4ac2fd11d
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/sykesm/zap-logfmt v0.0.2 // indirect
	github.com/yuin/gluamapper v0.0.0-20150323120927-d836955830e7
	github.com/yuin/gopher-lua v0.0.0-20190206043414-8bfc7677f583
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	layeh.com/gopher-luar v1.0.8-0.20190807124245-b07e371a3bb0
)

replace layeh.com/gopher-luar => github.com/layeh/gopher-luar v1.0.8-0.20190807124245-b07e371a3bb0

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a

go 1.13
