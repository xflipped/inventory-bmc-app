module github.com/foliagecp/inventory-bmc-app

go 1.19

replace git.fg-tech.ru/listware/proto => github.com/foliagecp/proto v0.1.6

replace git.fg-tech.ru/listware/cmdb => github.com/foliagecp/cmdb v0.1.5

replace git.fg-tech.ru/listware/go-core => github.com/foliagecp/go-core v0.1.6

// replace github.com/stmcginnis/gofish => github.com/Danile71/gofish v0.13.4-0.20230404103108-b33736b90757

require (
	git.fg-tech.ru/listware/cmdb v0.1.5
	git.fg-tech.ru/listware/go-core v0.1.6
	git.fg-tech.ru/listware/proto v0.1.6
	github.com/hashicorp/go-multierror v1.1.1
	github.com/koron/go-ssdp v0.0.3
	github.com/manifoldco/promptui v0.9.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stmcginnis/gofish v0.14.1-0.20230406231140-da46c1326ae9
	github.com/urfave/cli/v2 v2.24.3
	go.mongodb.org/mongo-driver v1.11.3
	go.uber.org/goleak v1.2.0
	google.golang.org/grpc v1.52.1
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/apache/flink-statefun/statefun-sdk-go/v3 v3.2.0 // indirect
	github.com/arangodb/go-driver v1.4.1 // indirect
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230111030713-bf00bc1b83b6 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.15.14 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
)
