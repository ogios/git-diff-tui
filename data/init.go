package data

import (
	"log"
	"time"
)

func init() {
	start := time.Now().UnixMicro()
	initGlobal()
	initTemp()
	log.Println("init data cost:", time.Now().UnixMicro()-start)
}

func Exit() {
	exitTemp()
}
