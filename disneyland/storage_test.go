package disneyland

import (
	"testing"
)

func checkTestErr(err error, t *testing.T) {
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestJobCRUD(t *testing.T) {
	initTestsConfig()
	storage, err := NewDisneylandStorage(TestsConfig.DatabaseURI)
	checkTestErr(err, t)

	if storage == nil {
		t.Fail()
	}

	job := &Job{
		Project:  "test_project",
		Status:   Job_PENDING,
		Metadata: `{"a": 123}`,
		Kind:     "kind_test",
	}

	createdJob, err := storage.CreateJob(job, User{Username: "tester"})
	checkTestErr(err, t)

	if createdJob == nil {
		t.Fail()
	}

	if createdJob.Project != job.Project {
		t.Fail()
	}

	if createdJob.Metadata != job.Metadata {
		t.Fail()
	}

}
