package optimus

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
    storage, err := NewOptimusStorage(TestsConfig.DatabaseURI)
    checkTestErr(err, t)

    if storage == nil {
        t.Fail()
    }

    job := &Job{
        Project:     "test_project",
        Status:      Job_PENDING,
        Coordinate:  "[0,0]",
        MetricValue: "9.0",
        Metadata:    `{"a": 123}`,
    }

    created_job, err := storage.CreateJob(job, User{Username: "tester"})
    checkTestErr(err, t)

    if created_job == nil {
        t.Fail()
    }

    if created_job.Project != job.Project {
        t.Fail()
    }

    if created_job.Coordinate != job.Coordinate {
        t.Fail()
    }

    if created_job.Metadata != job.Metadata {
        t.Fail()
    }

    if created_job.MetricValue != job.MetricValue {
        t.Fail()
    }
}
