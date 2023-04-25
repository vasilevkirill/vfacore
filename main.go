package main

import (
	"log"
	"main/ldap"
	"main/queue"
	radius "main/radius"
	tg "main/tg"
)

var rsecret = "los9diufmhspod9ufhmspdioufhspoidh"

func main() {
	config := LoadConfig()
	queue.InitQ()
	var LdapConfig ldap.Config
	err := config.UnmarshalKey("ldap", &LdapConfig)
	err = ldap.Run(LdapConfig)
	checkFatal(err)
	var tgConfig tg.Config
	err = config.UnmarshalKey("telegram", &tgConfig)
	checkFatal(err)
	err = tg.Run(tgConfig)
	checkFatal(err)

	var radiusConfig radius.Config
	err = config.UnmarshalKey("radius", &radiusConfig)
	checkFatal(err)
	err = radius.Run(radiusConfig)
	checkFatal(err)

}

func checkFatal(err error) {
	if err == nil {
		return
	}
	log.Panic(err)
}
