package optimus

import (
	"os"
	"testing"
)

func checkTestErr(err error, t *testing.T) {
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestPointCRUD(t *testing.T) {
	db_uri := os.Getenv("OPTIMUS_TEST_DB")
	storage, err := NewOptimusStorage(db_uri)
	checkTestErr(err, t)

	if storage == nil {
		t.Fail()
	}

	point := &Point{
		Project:     "test_project",
		Status:      Point_PENDING,
		Coordinate:  "[0,0]",
		MetricValue: "9.0",
		Metadata:    `{"a": 123}`,
	}

	created_point, err := storage.CreatePoint(point)
	checkTestErr(err, t)

	if created_point == nil {
		t.Fail()
	}

	if created_point.Project != point.Project {
		t.Fail()
	}

	if created_point.Coordinate != point.Coordinate {
		t.Fail()
	}

	if created_point.Metadata != point.Metadata {
		t.Fail()
	}

	if created_point.MetricValue != point.MetricValue {
		t.Fail()
	}
}
