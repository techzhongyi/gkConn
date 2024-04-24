module gkConn

go 1.2

require github.com/gorilla/websocket v1.5.1 // indirect

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	golang.org/x/net v0.24.0 // indirect

)

require (
	github.com/IBM/sarama v1.43.1
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/sirupsen/logrus v1.9.0
	github.com/techzhongyi/comlibgo v0.0.0
)

replace github.com/techzhongyi/comlibgo => ../comlibgo
