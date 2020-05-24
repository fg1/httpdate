package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
	"time"
)

// #include <errno.h>
// #include <stdint.h>
// #include <time.h>
//
// int settime(int64_t tv_sec, int64_t tv_nsec)
// {
//     struct timespec ts;
//     ts.tv_sec = tv_sec;
//     ts.tv_nsec = tv_nsec;
//     int r = clock_settime (CLOCK_REALTIME, &ts);
//     if (r == 0) {
//         return r;
//     }
//     return errno;
// }
import "C"

var url string
var systohc bool

func init() {
	flag.StringVar(&url, "url", "http://google.com", "URL to use for getting date")
	flag.BoolVar(&systohc, "systohc", false, "Set hardware clock from the system clock using hwclock")
	flag.Parse()
}

func main() {
	// Performs HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Parse date header
	t, err := time.Parse(http.TimeFormat, resp.Header.Get("Date"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)

	// Set system clock
	ret := syscall.Errno(C.settime(C.long(t.Unix()), 0))
	if ret != 0 {
		log.Fatal("Error setting date: ", ret.Error())
	}

	if systohc {
		cmd := exec.Command("hwclock", "--systohc")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
