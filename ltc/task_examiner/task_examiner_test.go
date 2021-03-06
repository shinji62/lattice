package task_examiner_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/lattice/ltc/task_examiner"
	"github.com/cloudfoundry-incubator/receptor"
	"github.com/cloudfoundry-incubator/receptor/fake_receptor"
)

var _ = Describe("TaskExaminer", func() {

	var (
		fakeReceptorClient *fake_receptor.FakeClient
		taskExaminer       task_examiner.TaskExaminer
	)

	BeforeEach(func() {
		fakeReceptorClient = &fake_receptor.FakeClient{}
		taskExaminer = task_examiner.New(fakeReceptorClient)
	})

	Describe("TaskStatus", func() {

		BeforeEach(func() {
			getTaskResponse := receptor.TaskResponse{
				TaskGuid:      "boop",
				State:         receptor.TaskStateCompleted,
				CellID:        "cell-01",
				Failed:        false,
				FailureReason: "",
				Result:        "some-result",
			}
			fakeReceptorClient.GetTaskReturns(getTaskResponse, nil)
		})

		It("returns a task status", func() {
			taskInfo, err := taskExaminer.TaskStatus("boop")

			Expect(err).ToNot(HaveOccurred())

			Expect(taskInfo.TaskGuid).To(Equal("boop"))
			Expect(taskInfo.State).To(Equal(receptor.TaskStateCompleted))
			Expect(taskInfo.CellID).To(Equal("cell-01"))
			Expect(taskInfo.Failed).To(BeFalse())
			Expect(taskInfo.FailureReason).To(BeEmpty())
			Expect(taskInfo.Result).To(Equal("some-result"))

			Expect(fakeReceptorClient.GetTaskCallCount()).To(Equal(1))
			Expect(fakeReceptorClient.GetTaskArgsForCall(0)).To(Equal("boop"))
		})

		Context("when the receptor returns errors", func() {
			It("returns exists false for TaskNotFound", func() {
				receptorError := receptor.Error{Type: receptor.TaskNotFound, Message: "could not locate this"}
				fakeReceptorClient.GetTaskReturns(receptor.TaskResponse{}, receptorError)

				_, err := taskExaminer.TaskStatus("boop1")

				Expect(err).To(MatchError(task_examiner.TaskNotFoundErrorMessage))
			})

			It("bubbles up error for receptor Error anything but TaskNotFound", func() {
				receptorError := receptor.Error{Type: receptor.TaskGuidAlreadyExists, Message: "could not locate this"}
				fakeReceptorClient.GetTaskReturns(receptor.TaskResponse{}, receptorError)

				_, err := taskExaminer.TaskStatus("boop1")

				Expect(err).To(MatchError(receptorError))
			})

			It("bubbles up error for non-receptor Error", func() {
				fakeReceptorClient.GetTaskReturns(receptor.TaskResponse{}, errors.New("you done goofed"))

				_, err := taskExaminer.TaskStatus("boop")

				Expect(err).To(MatchError("you done goofed"))
			})
		})

	})
})
