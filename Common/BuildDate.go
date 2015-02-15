package Common

// +build cgo

/*
#include <time.h>
#include <stdlib.h>

int buildDate() {
	struct tm t;
	strptime(__DATE__, "%b %d %Y", &t);
	char buf[256] = {0};
	strftime(buf, sizeof(buf), "%Y%m%d", &t);
	return atoi(buf);
}

int buildTime() {
	struct tm t;
	strptime(__TIME__, "%H:%M:%S", &t);
	char buf[256] = {0};
	strftime(buf, sizeof(buf), "%H%M%S", &t);
	return atoi(buf);
}
*/
import "C"
import "flag"

func init() {
	flag.Uint64("--build--timestamp", BuildDateTime(), "build date-time, as a tag")
}

func BuildDateTime() uint64 {
	return uint64(C.buildDate())*1000000 + uint64(C.buildTime())
}
