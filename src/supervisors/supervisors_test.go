package supervisors

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestSupervisors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Supervisors Suite")
}

func logPrintln(text ...interface{}) {
	fmt.Println(text...)
}

var _ = Describe("Supervisors", func() {

	It("initialize with empty map", func() {
		subject := New("subject", logPrintln)
		Expect(len(subject.workerList)).To(Equal(0))
		Expect(subject.exited).To(Equal(false))
	})

	It("exits monitoring routine", func() {
		subject := New("subject", logPrintln)
		subject.channel <- "exit"
		time.Sleep(time.Second * 1)
		Expect(subject.exited).To(Equal(true))
	})

	It("starts worker", func() {
		subject := New("subject", logPrintln)
		subject.StartWorker("worker1", func() {
			subject.log("in side of function")
		}, "argForWorker")

		time.Sleep(time.Second * 3)
		Expect(len(subject.workerList)).To(Equal(0))
	})

})
