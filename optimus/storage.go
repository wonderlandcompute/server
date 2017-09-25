package optimus

import (
	"database/sql"
	_ "github.com/lib/pq"
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

func (storage *OptimusStorage) CreatePoint(point *Point, creator User) (*Point, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	created_point := &Point{}
	err = tx.QueryRow(`
		INSERT INTO points (project, status, coordinate, metric_value, metadata, creator)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, project, status, coordinate, metric_value, metadata;`,
		point.Project, point.Status, point.Coordinate, point.MetricValue, point.Metadata,
		creator.Username,
	).Scan(
		&created_point.Id,
		&created_point.Project,
		&created_point.Status,
		&created_point.Coordinate,
		&created_point.MetricValue,
		&created_point.Metadata,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return created_point, err
}

func (storage *OptimusStorage) GetPoint(id uint64, project string) (*Point, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	point := &Point{}
	err = tx.QueryRow(`
		SELECT id, project, status, coordinate, metric_value, metadata FROM points
		WHERE id=$1 AND project=$2;`,
		id, project,
	).Scan(
		&point.Id,
		&point.Project,
		&point.Status,
		&point.Coordinate,
		&point.MetricValue,
		&point.Metadata,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return point, err
}

func (storage *OptimusStorage) ListPoints(project string) (*ListOfPoints, error) {
	query := `
		SELECT id, project, status, coordinate, metric_value, metadata
		FROM points
		WHERE project=$1;
	`

	rows, err := storage.db.Query(query, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := &ListOfPoints{Points: []*Point{}}

	for rows.Next() {
		point := &Point{}

		err = rows.Scan(
			&point.Id,
			&point.Project,
			&point.Status,
			&point.Coordinate,
			&point.MetricValue,
			&point.Metadata,
		)

		if err != nil {
			return nil, err
		}
		ret.Points = append(ret.Points, point)
	}

	err = rows.Err()
	return ret, err
}

func (storage *OptimusStorage) UpdatePoint(point *Point) (*Point, error) {
	tx, err := storage.db.Begin()
	if err != nil {
		return nil, err
	}

	result_point := &Point{}
	err = tx.QueryRow(`
		UPDATE points
		SET
			status=$1,
			metric_value=$2,
			metadata=$3
		WHERE id=$4 and project=$5
		RETURNING id, project, status, coordinate, metric_value, metadata;`,
		point.Status,
		point.MetricValue,
		point.Metadata,
		point.Id,
		point.Project,
	).Scan(
		&result_point.Id,
		&result_point.Project,
		&result_point.Status,
		&result_point.Coordinate,
		&result_point.MetricValue,
		&result_point.Metadata,
	)
	if err != nil {
		return result_point, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return result_point, err
	}

	return result_point, err

}
