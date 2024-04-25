package data

import (
	"log"
	"time"
)

func init() {
	start := time.Now().UnixMicro()
	initGlobal()
	initTemp()
	initDiffsrc()
	log.Println("init data cost:", time.Now().UnixMicro()-start)
}

func Exit() {
	exitTemp()
}
