package repository

import (
	"alumnihub/internal/models"
	"context"
)

func (m *PostgresDBRepo) GetProfiles() ([]*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT ap.id, ap.alumni_id, ap.user_id, u.username, a.name 
			from alumni_profile ap JOIN users u ON ap.user_id = u.id JOIN alumni a ON ap.alumni_id = a.id
			ORDER BY ap.id DESC LIMIT 4;
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*models.Profile

	for rows.Next() {
		var profile models.Profile
		err := rows.Scan(
			&profile.ID,
			&profile.AlumniID,
			&profile.UserID,
			&profile.UserUsername,
			&profile.UserName,
		)
		if err != nil {
			return nil, err
		}

		photo, err := m.GetUserPhotoByID(profile.UserID)
		if err != nil {
			return nil, err
		}

		profile.Photo = photo

		profiles = append(profiles, &profile)
	}

	return profiles, nil

}

func (m *PostgresDBRepo) InsertProfile(profile models.Profile) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into alumni_profile (user_id, alumni_id)
			values ($1, $2) returning id`

	var userID int

	err := m.DB.QueryRowContext(ctx, stmt,
		profile.UserID,
		profile.AlumniID,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *PostgresDBRepo) UpdateProfile(profile models.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update alumni_profile set bio = $1, location = $2, sm_facebook = $3, sm_instagram = $4, sm_twitter = $5, sm_tiktok = $6
			where user_id = $7`

	_, err := m.DB.ExecContext(ctx, stmt,
		profile.Bio,
		profile.Location,
		profile.Facebook,
		profile.Instagram,
		profile.Twitter,
		profile.Tiktok,
		profile.UserID,
	)

	if err != nil {
		return err
	}

	stmt = `update users set photo = $1
			where id = $2`

	_, err = m.DB.ExecContext(ctx, stmt,
		profile.Photo,
		profile.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetProfileByAlumniID(id int) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, alumni_id, user_id,
			COALESCE(bio, ''), 
			COALESCE(location, ''), 
			COALESCE(sm_facebook, ''), 
			COALESCE(sm_instagram, ''), 
			COALESCE(sm_twitter, ''), 
			COALESCE(sm_tiktok, '')
			from alumni_profile where alumni_id = $1`

	var profile models.Profile
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&profile.ID,
		&profile.AlumniID,
		&profile.UserID,
		&profile.Bio,
		&profile.Location,
		&profile.Facebook,
		&profile.Instagram,
		&profile.Twitter,
		&profile.Tiktok,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil

}

func (m *PostgresDBRepo) GetProfileByUserID(id int) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select alumni_id, user_id, 
			COALESCE(bio, ''), 
			COALESCE(location, ''), 
			COALESCE(sm_facebook, ''), 
			COALESCE(sm_instagram, ''), 
			COALESCE(sm_twitter, ''), 
			COALESCE(sm_tiktok, '')
			from alumni_profile where user_id = $1`

	var profile models.Profile
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&profile.AlumniID,
		&profile.UserID,
		&profile.Bio,
		&profile.Location,
		&profile.Facebook,
		&profile.Instagram,
		&profile.Twitter,
		&profile.Tiktok,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (m *PostgresDBRepo) UpdateAdminProfile(profile models.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update admin_profile set bio = $1, name = $2, sm_facebook = $3, sm_instagram = $4, sm_twitter = $5, sm_tiktok = $6
			where user_id = $7`

	_, err := m.DB.ExecContext(ctx, stmt,
		profile.Bio,
		profile.UserName,
		profile.Facebook,
		profile.Instagram,
		profile.Twitter,
		profile.Tiktok,
		profile.UserID,
	)

	if err != nil {
		return err
	}

	stmt = `update users set photo = $1
			where id = $2`

	_, err = m.DB.ExecContext(ctx, stmt,
		profile.Photo,
		profile.UserID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAdminProfileByUserID(id int) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select user_id, name, 
			COALESCE(bio, ''),
			COALESCE(sm_facebook, ''), 
			COALESCE(sm_instagram, ''), 
			COALESCE(sm_twitter, ''), 
			COALESCE(sm_tiktok, '')
			from admin_profile where user_id = $1`

	var profile models.Profile
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&profile.UserID,
		&profile.UserName,
		&profile.Bio,
		&profile.Facebook,
		&profile.Instagram,
		&profile.Twitter,
		&profile.Tiktok,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (m *PostgresDBRepo) InsertAlumniEducation(education models.AlumniEducation) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into alumni_educations (user_id, school_name, school_degree, school_study_major, start_year, end_year, currently_studying, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt,
		education.UserID,
		education.School,
		education.Degree,
		education.StudyMajor,
		education.StartYear,
		education.EndYear,
		education.CurrentlyStudying,
		education.CreatedAt,
		education.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAlumniEducations(id int) ([]*models.AlumniEducation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, school_name, school_degree, school_study_major, start_year, end_year, currently_studying, created_at, updated_at from alumni_educations where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var educations []*models.AlumniEducation

	for rows.Next() {
		var education models.AlumniEducation
		err := rows.Scan(
			&education.ID,
			&education.UserID,
			&education.School,
			&education.Degree,
			&education.StudyMajor,
			&education.StartYear,
			&education.EndYear,
			&education.CurrentlyStudying,
			&education.CreatedAt,
			&education.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		educations = append(educations, &education)
	}

	return educations, nil
}

func (m *PostgresDBRepo) GetAlumniEducation(id int) (*models.AlumniEducation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, school_name, school_degree, school_study_major, start_year, end_year, currently_studying, created_at, updated_at from alumni_educations where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var education models.AlumniEducation

	err := row.Scan(
		&education.ID,
		&education.UserID,
		&education.School,
		&education.Degree,
		&education.StudyMajor,
		&education.StartYear,
		&education.EndYear,
		&education.CurrentlyStudying,
		&education.CreatedAt,
		&education.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &education, nil
}

func (m *PostgresDBRepo) DeleteAlumniEducations(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from alumni_educations where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertAlumniJob(alumnijob models.AlumniJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into alumni_jobs (user_id, position, company, company_location, employment_type, start_year, end_year, currently_working, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := m.DB.ExecContext(ctx, stmt,
		alumnijob.UserID,
		alumnijob.Position,
		alumnijob.Company,
		alumnijob.CompanyLocation,
		alumnijob.EmploymentType,
		alumnijob.StartYear,
		alumnijob.EndYear,
		alumnijob.CurrentlyWorking,
		alumnijob.CreatedAt,
		alumnijob.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAlumniJobs(id int) ([]*models.AlumniJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, position, company, company_location, employment_type, start_year, end_year, currently_working, created_at, updated_at from alumni_jobs where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumnijobs []*models.AlumniJob

	for rows.Next() {
		var alumnijob models.AlumniJob
		err := rows.Scan(
			&alumnijob.ID,
			&alumnijob.UserID,
			&alumnijob.Position,
			&alumnijob.Company,
			&alumnijob.CompanyLocation,
			&alumnijob.EmploymentType,
			&alumnijob.StartYear,
			&alumnijob.EndYear,
			&alumnijob.CurrentlyWorking,
			&alumnijob.CreatedAt,
			&alumnijob.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		alumnijobs = append(alumnijobs, &alumnijob)
	}

	return alumnijobs, nil
}

func (m *PostgresDBRepo) GetAlumniJob(id int) (*models.AlumniJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, position, company, company_location, employment_type, start_year, end_year, currently_working, created_at, updated_at from alumni_jobs where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var alumnijob models.AlumniJob

	err := row.Scan(
		&alumnijob.ID,
		&alumnijob.UserID,
		&alumnijob.Position,
		&alumnijob.Company,
		&alumnijob.CompanyLocation,
		&alumnijob.EmploymentType,
		&alumnijob.StartYear,
		&alumnijob.EndYear,
		&alumnijob.CurrentlyWorking,
		&alumnijob.CreatedAt,
		&alumnijob.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &alumnijob, nil
}

func (m *PostgresDBRepo) DeleteAlumniJobs(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from alumni_jobs where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
