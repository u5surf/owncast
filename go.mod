module github.com/owncast/owncast

go 1.14

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/amalfra/etag v0.0.0-20190921100247-cafc8de96bc5
	github.com/aws/aws-sdk-go v1.40.0
	github.com/go-fed/activity v1.0.0 // indirect
	github.com/go-fed/httpsig v1.1.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/grafov/m3u8 v0.11.1
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/markbates/pkger v0.17.1
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/microcosm-cc/bluemonday v1.0.15
	github.com/mssola/user_agent v0.5.3
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/nareix/joy5 v0.0.0-20200712071056-a55089207c88
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oschwald/geoip2-golang v1.5.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/schollz/sqlite3dump v1.3.0
	github.com/shirou/gopsutil v3.21.8+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/yuin/goldmark v1.4.1
	golang.org/x/mod v0.5.0
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mvdan.cc/xurls v1.1.0
)

replace github.com/go-fed/activity => github.com/owncast/activity v1.0.0
