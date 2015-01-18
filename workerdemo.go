package main

import (
    "fmt"
    "time"
)

var exitCall chan bool

var wakeupCall chan bool
var work chan int

var stats map[string]int

func main() {
    exitCall = make(chan bool, 1)
    wakeupCall = make(chan bool, 1)
    work = make(chan int, 1)

    stats = make(map[string]int)
    stats["con1"] = 0
    stats["con2"] = 0
    stats["con3"] = 0

    go workProvider()

    go consumer("con1", 1 * time.Second)
    go consumer("con2", 1 * time.Second)
    go consumer("con3", 1 * time.Second)

    _ = <-exitCall
    time.Sleep(1 * time.Second) //give consumers time to process last batch, if any
    fmt.Println(stats)


}

func workProvider() {
    workid := 0

    for {
        //sleep until someone wakes me up
        _ = <-wakeupCall
        //push some work to the consumer
        work <- workid
        workid++

        if workid == 30 {
            exitCall <- true
        }

    }

}

func consumer(myname string, sleepTime time.Duration) {
    fmt.Println(myname, " entry with ", sleepTime)
    for {
        wakeupCall <- true
        workid := <-work
        time.Sleep(sleepTime)
        fmt.Println(myname, " received ", workid)
        stats[myname]++
    }
}

