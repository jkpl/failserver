package main

import (
	"errors"
	"fmt"
	"github.com/pinterest/bender"
	"github.com/pinterest/bender/hist"
	bhttp "github.com/pinterest/bender/http"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func SyntheticHttpRequests(n int) chan interface{} {
	c := make(chan interface{}, 100)
	go func() {
		for i := 0; i < n; i++ {
			req, err := http.NewRequest("POST", "http://localhost:8080/", nil)
			if err != nil {
				panic(err)
			}
			c <- req
		}
		close(c)
	}()
	return c
}

func bodyValidator(request interface{}, resp *http.Response) error {
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
	}

	return nil
}

func main() {
	intervals := bender.ExponentialIntervalGenerator(10)
	requests := SyntheticHttpRequests(100)
	exec := bhttp.CreateExecutor(nil, nil, bodyValidator)
	recorder := make(chan interface{}, 100)

	bender.LoadTestThroughput(intervals, requests, exec, recorder)

	l := log.New(os.Stdout, "", log.LstdFlags)
	h := hist.NewHistogram(60000, int(time.Millisecond))
	bender.Record(recorder, bender.NewLoggingRecorder(l), bender.NewHistogramRecorder(h))
	fmt.Println(h)

}
