module github.com/ariefdarmawan/kmsg

go 1.16

replace github.com/ariefdarmawan/flexpg v0.7.0 => github.com/ariefdarmawan/flexpg v0.7.1

require (
	git.kanosolution.net/kano/dbflex v1.2.5
	git.kanosolution.net/kano/kaos v0.2.1
	github.com/ariefdarmawan/datahub v0.2.4
	github.com/sebarcode/codekit v0.1.1
	go.mongodb.org/mongo-driver v1.9.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
