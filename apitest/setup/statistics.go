package setup

import (
	"sync"
	"testing"
	"time"
)

// TestStatistics represents the statistics of a specific test
type TestStatistics struct {
	Name         string
	SetupTime    time.Duration
	TeardownTime time.Duration
}

// Clone returns an exact deep copy of the object
func (ts *TestStatistics) Clone() *TestStatistics {
	if ts == nil {
		return nil
	}
	return &TestStatistics{
		Name:         ts.Name,
		SetupTime:    ts.SetupTime,
		TeardownTime: ts.TeardownTime,
	}
}

// StatisticsRecorder represents the statistics recorder and computer
type StatisticsRecorder struct {
	lock   *sync.Mutex
	tests  []*TestStatistics
	byName map[string]*TestStatistics
}

// Statistics represents the final computed statistics
type Statistics struct {
	Tests            []*TestStatistics
	MinSetupTime     time.Duration
	MinSetupTimeTest string

	MaxSetupTime     time.Duration
	MaxSetupTimeTest string

	MinTeardownTime     time.Duration
	MinTeardownTimeTest string

	MaxTeardownTime     time.Duration
	MaxTeardownTimeTest string

	AvgSetupTime    time.Duration
	AvgTeardownTime time.Duration
}

// NewStatisticsRecorder constructs a new statistics recorder instance
func NewStatisticsRecorder() *StatisticsRecorder {
	return &StatisticsRecorder{
		lock:   &sync.Mutex{},
		tests:  make([]*TestStatistics, 0),
		byName: make(map[string]*TestStatistics),
	}
}

// Set allows to safely modify a certain tests statistics.
// It will automatically create a test if it's not yet registered
func (sr *StatisticsRecorder) Set(
	t *testing.T,
	mutator func(*TestStatistics),
) {
	testName := t.Name()
	sr.lock.Lock()
	stats, defined := sr.byName[testName]
	if defined {
		mutator(stats)
		stats.Name = testName
	} else {
		newTestStats := &TestStatistics{}
		mutator(newTestStats)
		newTestStats.Name = testName
		sr.byName[testName] = newTestStats
		sr.tests = append(sr.tests, newTestStats)
	}
	sr.lock.Unlock()
}

// Compute will compute and return the final statistics based on the recordings
func (sr *StatisticsRecorder) Compute() *Statistics {
	sr.lock.Lock()
	copiedTests := make([]*TestStatistics, len(sr.tests))

	for i, tst := range sr.tests {
		copiedTests[i] = tst.Clone()
	}
	sr.lock.Unlock()

	if len(copiedTests) < 0 {
		return &Statistics{
			Tests: make([]*TestStatistics, 0),
		}
	}

	// Compute statistics
	var minSetupTimeTest, maxSetupTimeTest,
		minTeardownTimeTest, maxTeardownTimeTest string
	var minSetupTime, maxSetupTime, avgSetupTime time.Duration
	var minTeardownTime, maxTeardownTime, avgTeardownTime time.Duration

	for i, testStats := range copiedTests {
		// Determine min setup time
		if minSetupTime == 0 || testStats.SetupTime < minSetupTime {
			minSetupTime = testStats.SetupTime
			minSetupTimeTest = testStats.Name
		}
		// Determine max setup time
		if maxSetupTime == 0 || testStats.SetupTime > maxSetupTime {
			maxSetupTime = testStats.SetupTime
			maxSetupTimeTest = testStats.Name
		}

		// Determine min teardown time
		if minTeardownTime == 0 || testStats.TeardownTime < minTeardownTime {
			minTeardownTime = testStats.TeardownTime
			minTeardownTimeTest = testStats.Name
		}
		// Determine max teardown time
		if maxTeardownTime == 0 || testStats.TeardownTime > maxTeardownTime {
			maxTeardownTime = testStats.TeardownTime
			maxTeardownTimeTest = testStats.Name
		}

		// Determine average setup time
		if avgSetupTime == 0 {
			avgSetupTime = testStats.SetupTime
		} else {
			avgSetupTime = avgSetupTime + (testStats.SetupTime-avgSetupTime)/time.Duration(i)
		}
		// Determine average teardown time
		if avgTeardownTime == 0 {
			avgTeardownTime = testStats.TeardownTime
		} else {
			avgTeardownTime = avgTeardownTime + (testStats.TeardownTime-avgTeardownTime)/time.Duration(i)
		}
	}

	return &Statistics{
		Tests: copiedTests,

		MinSetupTime:     minSetupTime,
		MinSetupTimeTest: minSetupTimeTest,

		MaxSetupTime:     maxSetupTime,
		MaxSetupTimeTest: maxSetupTimeTest,

		MinTeardownTime:     minTeardownTime,
		MinTeardownTimeTest: minTeardownTimeTest,

		MaxTeardownTime:     maxTeardownTime,
		MaxTeardownTimeTest: maxTeardownTimeTest,

		AvgSetupTime:    avgSetupTime,
		AvgTeardownTime: avgTeardownTime,
	}
}
