package optimus

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lib/pq"
)

type OptimusStorageConfig struct {
	DatabaseURI string `json:"db_uri"`
}

type OptimusStorage struct {
	db     *sql.DB
	Config OptimusStorageConfig
}

func NewOptimusStorage(db_uri string) (*OptimusStorage, error) {
	ret := &OptimusStorage{
		Config: OptimusStorageConfig{DatabaseURI: db_uri},
	}

	err := ret.Connect()

	return ret, err
}

func (storage *OptimusStorage) Connect() error {
	db, err := sql.Open("postgres", storage.Config.DatabaseURI)
	storage.db = db
	return err
}

func (storage *OptimusStorage) CreateJob(job *Job, creator User) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	created_job := &Job{}
	err = tx.QueryRow(`
		INSERT INTO jobs (project, status, coordinate, metric_value, metadata, creator, input, output, kind)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, project, status, coordinate, metric_value, metadata, input, output, kind;`,
		job.Project, job.Status, job.Coordinate, job.MetricValue, job.Metadata,
		creator.Username, job.Input, job.Output, job.Kind,
	).Scan(
		&created_job.Id,
		&created_job.Project,
		&created_job.Status,
		&created_job.Coordinate,
		&created_job.MetricValue,
		&created_job.Metadata,
		&created_job.Input,
		&created_job.Output,
		&created_job.Kind,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return created_job, err
}

func (storage *OptimusStorage) CreateMultipleJobs(jobs []*Job, creator User, project string) ([]*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}
	list_jobs := []*Job{}
	stmt, _ := tx.Prepare(pq.CopyIn("jobs", "project", "status", "coordinate", "metric_value", "metadata", "input", "output", "kind", "creator")) // MessageDetailRecord is the table name
	for _, job := range jobs {
		created_job := &Job{}
		job.Project = project
		_, err := stmt.Exec(job.Project, job.Status, job.Coordinate, job.MetricValue, job.Metadata, job.Input, job.Output, job.Kind, creator.Username)
		if err != nil {
			return nil, err
		}
		list_jobs = append(list_jobs, created_job)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return list_jobs, err
}

func (storage *OptimusStorage) GetJob(id uint64, project string) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	job := &Job{}
	err = tx.QueryRow(`
		SELECT id, project, status, coordinate, metric_value, metadata, input, output, kind FROM jobs
		WHERE id=$1 AND project=$2;`,
		id, project,
	).Scan(
		&job.Id,
		&job.Project,
		&job.Status,
		&job.Coordinate,
		&job.MetricValue,
		&job.Metadata,
		&job.Input,
		&job.Output,
		&job.Kind,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return job, err
}

func (storage *OptimusStorage) ListJobs(project string) (*ListOfJobs, error) {
	query := `SELECT id, project, status, coordinate, metric_value, metadata, input, output, kind
		FROM jobs
		WHERE project=$1;`

	rows, err := storage.db.Query(query, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := &ListOfJobs{Jobs: []*Job{}}

	for rows.Next() {
		job := &Job{}

		err = rows.Scan(
			&job.Id,
			&job.Project,
			&job.Status,
			&job.Coordinate,
			&job.MetricValue,
			&job.Metadata,
			&job.Input,
			&job.Output,
			&job.Kind,
		)

		if err != nil {
			return nil, err
		}
		ret.Jobs = append(ret.Jobs, job)
	}

	err = rows.Err()
	return ret, err
}

func (storage *OptimusStorage) UpdateJob(job *Job) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	result_job := &Job{}
	err = tx.QueryRow(`
		UPDATE jobs
		SET
			status=$1,
			metric_value=$2,
			metadata=$3,
			input=$4,
			output=$5,
			kind=$6
		WHERE id=$7 and project=$8
		RETURNING id, project, status, coordinate, metric_value, metadata, input, output, kind;`,
		job.Status,
		job.MetricValue,
		job.Metadata,
		job.Input,
		job.Output,
		job.Kind,
		job.Id,
		job.Project,
	).Scan(
		&result_job.Id,
		&result_job.Project,
		&result_job.Status,
		&result_job.Coordinate,
		&result_job.MetricValue,
		&result_job.Metadata,
		&result_job.Input,
		&result_job.Output,
		&result_job.Kind,
	)
	if err != nil {
		return result_job, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return result_job, err
	}

	return result_job, err
}

func (storage *OptimusStorage) PullJobs(how_many uint32) ([]*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	result_jobs := []*Job{}
	rows, err := tx.Query(`
		WITH pulled_pts as (
			SELECT id, project, kind
			FROM jobs
			WHERE status=$1
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)
		UPDATE jobs pts
		SET status=$3
		FROM pulled_pts
		WHERE pulled_pts.id=pts.id and pulled_pts.project=pts.project and pulled_pts.kind=pts.kind
		RETURNING pts.id, pts.project, pts.status, pts.coordinate, pts.metric_value, pts.metadata, pts.input, pts.output, pts.kind;`,
		Job_PENDING,
		how_many,
		Job_PULLED,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		job := Job{}
		err = rows.Scan(&job.Id, &job.Project, &job.Status, &job.Coordinate, &job.MetricValue, &job.Metadata, &job.Input, &job.Output, &job.Kind)
		if err != nil {
			return nil, err
		}

		result_jobs = append(result_jobs, &job)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return result_jobs, err
}
