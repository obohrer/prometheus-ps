package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mitchellh/go-ps"
	"github.com/stretchr/testify/assert"
)

type FakeProcess struct {
	name string
}

func (p FakeProcess) Pid() int {
	return 0
}

func (p FakeProcess) PPid() int {
	return 0
}

func (p FakeProcess) Executable() string {
	return p.name
}

func TestGroupByName(t *testing.T) {
	assert := assert.New(t)
	var a [4]ps.Process

	a[0] = FakeProcess{name: "a"}
	a[1] = FakeProcess{name: "b"}
	a[2] = FakeProcess{name: "a"}
	a[3] = FakeProcess{name: "c"}
	pl := make([]ps.Process, len(a))
	for i, v := range a {
		pl[i] = v
	}

	watchList := make(map[string]Void)
	watchListA := [...]string{"a", "c", "e", "x"}
	for _, x := range watchListA {
		watchList[x] = Void{}
	}

	grouped := groupByName(pl, watchList)
	assert.Equal(len(grouped), 2, "There should be only 2 processes")
	assert.Equal(grouped["a"], 2, "a is twice in the list")
	assert.Equal(grouped["c"], 1, "c is once in the list")
}

func TestBuildMetricName(t *testing.T) {
	// test special chars + upper case
	assert := assert.New(t)
	assert.Equal(buildMetricName("yo"), "process_yo_up")
	assert.Equal(buildMetricName("Yo"), "process_yo_up")
	assert.Equal(buildMetricName("HellO Test"), "process_hello_test_up")
}

func TestWriteProcessesMetrics(t *testing.T) {
	assert := assert.New(t)
	rr := httptest.NewRecorder()
	pl := make(map[string]int)
	pl["a"] = 2
	pl["c"] = 1
	writeProcessesMetrics(rr, pl)

	expectedOutputPartA := `# HELP process_a_up The process process_a_up is up
# TYPE process_a_up gauge
process_a_up 2`
	expectedOutputPartB := `# HELP process_c_up The process process_c_up is up
# TYPE process_c_up gauge
process_c_up 1`

	assert.True(strings.Contains(rr.Body.String(), expectedOutputPartA))
	assert.True(strings.Contains(rr.Body.String(), expectedOutputPartB))

}
