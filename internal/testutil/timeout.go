// Copyright 2019 Google Inc. All Rights Reserved.
// This file is available under the Apache license.
package testutil

import (
	"testing"
	"time"

	"github.com/golang/glog"
)

// DoOrTimeout runs a check function every interval until deadline, unless the
// check returns true.  The check should return false otherwise. If the check
// returns an error the check is immediately failed.
func DoOrTimeout(do func() (bool, error), deadline, interval time.Duration) (bool, error) {
	timeout := time.After(deadline)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			return false, nil
		case <-ticker.C:
			glog.V(2).Infof("tick")
			ok, err := do()
			glog.V(2).Infof("ok, err: %v %v", ok, err)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
			// otherwise wait and retry
		}
	}
}

// TimeoutTest returns a test function that executes f with a timeout, If the
// test does not complete in time the test is failed.  This lets us set a
// per-test timeout instead of the global `go test -timeout` coarse timeout.
func TimeoutTest(timeout time.Duration, f func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		done := make(chan bool)
		go func() {
			t.Helper()
			defer close(done)
			f(t)
		}()
		select {
		case <-time.After(timeout * RaceDetectorMultiplier):
			t.Fatal("timed out")
		case <-done:
		}
	}
}
