package disneyland

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DisneylandStorageConfig struct {
	DatabaseURI string `json:"db_uri"`
}

type DisneylandStorage struct {
	db     *sql.DB
	Config DisneylandStorageConfig
}

func NewDisneylandStorage(db_uri string) (*DisneylandStorage, error) {
	ret := &DisneylandStorage{
		Config: DisneylandStorageConfig{DatabaseURI: db_uri},
	}

	err := ret.Connect()

	return ret, err
}

func (storage *DisneylandStorage) Connect() error {
	db, err := sql.Open("postgres", storage.Config.DatabaseURI)
	storage.db = db
	return err
}

func (storage *DisneylandStorage) CreateJob(job *Job, creator User) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	created_job := &Job{}
	err = tx.QueryRow(`
		INSERT INTO jobs (project, status, metadata, creator, input, output, kind)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, project, status, metadata, input, output, kind;`,
		job.Project, job.Status, job.Metadata, creator.Username, job.Input, job.Output, job.Kind,
	).Scan(
		&created_job.Id,
		&created_job.Project,
		&created_job.Status,
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

func (storage *DisneylandStorage) GetJob(id uint64) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}
	job := &Job{}

	strQuery := `SELECT id, project, status,  metadata, input, output, kind
					FROM jobs
					WHERE id=$1;`
	err = tx.QueryRow(strQuery, id).Scan(
		&job.Id,
		&job.Project,
		&job.Status,
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

func (storage *DisneylandStorage) ListJobs(project string, kind string) (*ListOfJobs, error) {
	strQuery := `SELECT id, project, status, metadata, input, output, kind
			  FROM jobs
			  WHERE project=$1 and kind=$2;`

	rows, err := storage.db.Query(strQuery, project, kind)
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

func (storage *DisneylandStorage) UpdateJob(job *Job) (*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	result_job := &Job{}
	err = tx.QueryRow(`
		UPDATE jobs
		SET
			status=$1,
			metadata=$2,
			input=$3,
			output=$4,
			kind=$5
		WHERE id=$6 and project=$7
		RETURNING id, project, status, metadata, input, output, kind;`,
		job.Status,
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

func (storage *DisneylandStorage) PullJobs(how_many uint32, kind string) ([]*Job, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}
	result_jobs := []*Job{}
	sqlStr := `
	WITH pulled_pts as (
		SELECT id, project, kind
		FROM jobs
		WHERE status=$1 and kind=$2
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	)
	UPDATE jobs pts
	SET status=$4
	FROM pulled_pts
	WHERE pulled_pts.id=pts.id and pulled_pts.project=pts.project and pulled_pts.kind=pts.kind
	RETURNING pts.id, pts.project, pts.status, pts.metadata, pts.input, pts.output, pts.kind;`

	rows, err := tx.Query(sqlStr, Job_PENDING, kind, how_many, Job_PULLED)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		job := Job{}
		err = rows.Scan(&job.Id, &job.Project, &job.Status, &job.Metadata, &job.Input, &job.Output, &job.Kind)
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
