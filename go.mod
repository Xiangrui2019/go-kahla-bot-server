module github.com/xiangrui2019/go-kahla-bot-server

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Unknwon/com v0.0.0-20190321035513-0fed4efef755 // indirect
	github.com/go-macaron/binding v0.0.0-20170611065819-ac54ee249c27
	github.com/go-macaron/inject v0.0.0-20160627170012-d8a0b8677191 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.3
	github.com/gorilla/websocket v1.4.0
	github.com/malisit/kolpa v0.0.0-20171128201422-10c4658a62cb
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/stretchr/testify v1.3.0
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/macaron.v1 v1.3.2
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.26.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190123085648-057139ce5d2b
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3
	golang.org/x/net => github.com/golang/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180905080454-ebe1bf3edb33
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190328211700-ab21143f2384
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc => github.com/grpc/grpc-go v1.20.1
)
