package aaa

import "os"

func init() {
	os.Setenv("RUNEWIDTH_EASTASIAN", "true")
	os.Setenv("LC_CTYPE", "en_US.UTF-8")
}
