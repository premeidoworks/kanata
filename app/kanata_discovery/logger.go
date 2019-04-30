package main

import "log"

func logErr(v ...interface{}) {
	log.Println(v...)
}

func logInfo(v ...interface{}) {
	log.Panicln(v...)
}
