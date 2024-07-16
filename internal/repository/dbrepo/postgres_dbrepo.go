package dbrepo

import (
	"alumnihub/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeOut = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into users (username, email, password, is_admin, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6) returning id`

	var userID int

	err := m.DB.QueryRowContext(ctx, stmt,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *PostgresDBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, password, is_admin, created_at, updated_at 
			from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, password, is_admin, created_at, updated_at 
			from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (m *PostgresDBRepo) GetUserUsernameByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select username from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.Username,
	)

	if err != nil {
		return "", err
	}

	return user.Username, nil
}

func (m *PostgresDBRepo) GetUserIDByUsername(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id from users where username = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
	)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (m *PostgresDBRepo) GetUserPhotoByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select photo from users where id = $1`

	var photo sql.NullString
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&photo)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	if photo.Valid {
		return photo.String, nil
	}

	return "", nil
}

func (m *PostgresDBRepo) AllAlumni() ([]*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, nisn, nis, name, gender, phone, graduation_year, class from alumni order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumnis []*models.Alumni

	for rows.Next() {
		var alumni models.Alumni
		err := rows.Scan(
			&alumni.ID,
			&alumni.NISN,
			&alumni.NIS,
			&alumni.Name,
			&alumni.Gender,
			&alumni.Phone,
			&alumni.Year,
			&alumni.Class,
		)
		if err != nil {
			return nil, err
		}

		alumnis = append(alumnis, &alumni)
	}

	return alumnis, nil
}

func (m *PostgresDBRepo) Alumni(id int) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, nisn, nis, name, gender, phone, graduation_year, class
				FROM alumni
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var alumni models.Alumni

	err := row.Scan(
		&alumni.ID,
		&alumni.NISN,
		&alumni.NIS,
		&alumni.Name,
		&alumni.Gender,
		&alumni.Phone,
		&alumni.Year,
		&alumni.Class,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func (m *PostgresDBRepo) InsertAlumni(alumni models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into alumni (nisn, nis, name, gender, phone, graduation_year, class)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		alumni.NISN,
		alumni.NIS,
		alumni.Name,
		alumni.Gender,
		alumni.Phone,
		alumni.Year,
		alumni.Class,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) UpdateAlumni(alumni models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update alumni set name = $1, gender = $2, phone = $3, graduation_year = $4, class = $5, nisn = $6, nis = $7 
			where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		alumni.Name,
		alumni.Gender,
		alumni.Phone,
		alumni.Year,
		alumni.Class,
		alumni.NISN,
		alumni.NIS,
		alumni.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteAlumni(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from alumni where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAlumniByNISN(nisn string) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, nisn, nis, name, gender, phone, graduation_year, class 
			from alumni where nisn = $1`

	var alumni models.Alumni
	row := m.DB.QueryRowContext(ctx, query, nisn)

	err := row.Scan(
		&alumni.ID,
		&alumni.NISN,
		&alumni.NIS,
		&alumni.Name,
		&alumni.Gender,
		&alumni.Phone,
		&alumni.Year,
		&alumni.Class,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil

}

func (m *PostgresDBRepo) GetAlumniNameByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select name from alumni where id = $1`

	var alumni models.Alumni
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&alumni.Name,
	)

	if err != nil {
		return "", err
	}

	return alumni.Name, nil
}

func (m *PostgresDBRepo) CountAlumni() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT COUNT(id) as count_alumni FROM alumni`

	row := m.DB.QueryRowContext(ctx, query)

	var count int

	err := row.Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, nil

}

func (m *PostgresDBRepo) CountAlumniAccount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT COUNT(id) as count_alumni FROM alumni_profile`

	row := m.DB.QueryRowContext(ctx, query)

	var count int

	err := row.Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, nil

}

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

func (m *PostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, title, slug, body, image, status, created_at, updated_at, published_at from articles order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Body,
			&article.Image,
			&article.Status,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *PostgresDBRepo) Article(id int) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, slug, body, image, status, created_at, updated_at, published_at
				FROM articles
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var article models.Article

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Body,
		&article.Image,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *PostgresDBRepo) ArticleBySlug(slug string) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, slug, body, image, status, created_at, updated_at, published_at
				FROM articles
				WHERE slug = $1
			`

	row := m.DB.QueryRowContext(ctx, query, slug)

	var article models.Article

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Body,
		&article.Image,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *PostgresDBRepo) InsertArticle(article models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into articles (title, slug, body, image, status, created_at, updated_at, published_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
		article.Image,
		article.Status,
		article.CreatedAt,
		article.UpdatedAt,
		article.PublishedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateArticle(article models.Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update articles set title = $1, slug = $2, body = $3,
				status = $4, updated_at = $5, published_at = $6, image = $7 
				where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
		article.Status,
		article.UpdatedAt,
		article.PublishedAt,
		article.Image,
		article.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteArticle(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from articles where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) AllForms() ([]*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, title, description, hidden, has_time_limit, start_date, end_date, created_at, updated_at 
				from forms order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []*models.Form

	for rows.Next() {
		var form models.Form
		err := rows.Scan(
			&form.ID,
			&form.Title,
			&form.Description,
			&form.Hidden,
			&form.HasTimeLimit,
			&form.StartDate,
			&form.EndDate,
			&form.CreatedAt,
			&form.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		forms = append(forms, &form)
	}

	return forms, nil
}

func (m *PostgresDBRepo) Form(id int) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, description, hidden, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
		&form.Hidden,
		&form.HasTimeLimit,
		&form.StartDate,
		&form.EndDate,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	query = `
				SELECT id, question_text, type, extension, created_at, updated_at
				FROM questions
				WHERE form_id = $1
				ORDER BY id
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID,
			&question.Question,
			&question.Type,
			&question.Extension,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		questions = append(questions, &question)
	}

	form.Questions = questions

	return &form, nil
}

func (m *PostgresDBRepo) ShowForm(id int) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, description, hidden, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
		&form.Hidden,
		&form.HasTimeLimit,
		&form.StartDate,
		&form.EndDate,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	fullQuestion, err := m.QuestionsByForm(form.ID)
	if err != nil {
		return nil, err
	}

	form.Questions = fullQuestion

	return &form, nil
}

func (m *PostgresDBRepo) InsertForm(form models.Form) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into forms (title, description, has_time_limit, start_date,
			end_date, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		form.Title,
		form.Description,
		form.HasTimeLimit,
		form.StartDate,
		form.EndDate,
		form.CreatedAt,
		form.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateForm(form models.Form) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update forms set title = $1, description = $2, has_time_limit = $3, start_date = $4, end_date = $5, hidden = $6, updated_at = $7
				where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		form.Title,
		form.Description,
		form.HasTimeLimit,
		form.StartDate,
		form.EndDate,
		form.Hidden,
		form.UpdatedAt,
		form.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteForm(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from forms where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) ShowFormAnswers(id int) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, description, hidden, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
		&form.Hidden,
		&form.HasTimeLimit,
		&form.StartDate,
		&form.EndDate,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	query = `
				SELECT id, question_text, type, extension, created_at, updated_at
				FROM questions
				WHERE form_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID,
			&question.Question,
			&question.Type,
			&question.Extension,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get answers for each question
		answerGroups, err := m.GroupAnswersByQuestion(form.ID, question.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}

		question.GroupAnswer = answerGroups

		questions = append(questions, &question)
	}

	form.Questions = questions

	return &form, nil
}

func (m *PostgresDBRepo) Question(id int) (*models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, form_id, question_text, type, extension, created_at, updated_at
				FROM questions
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var question models.Question

	err := row.Scan(
		&question.ID,
		&question.FormID,
		&question.Question,
		&question.Type,
		&question.Extension,
		&question.CreatedAt,
		&question.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	query = `
				SELECT id, question_id, option_text
				FROM options
				WHERE question_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var options []*models.Option
	for rows.Next() {
		var option models.Option
		err := rows.Scan(
			&option.ID,
			&option.QuestionID,
			&option.Option,
		)
		if err != nil {
			return nil, err
		}

		options = append(options, &option)
	}

	question.Options = options

	// Get answers for each question
	answerQuery := `
		SELECT id, user_id, form_id, question_id, answer_text
		FROM answers
		WHERE question_id = $1
	`

	answerRows, err := m.DB.QueryContext(ctx, answerQuery, question.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer answerRows.Close()

	var answers []*models.Answer
	for answerRows.Next() {
		var answer models.Answer
		err := answerRows.Scan(
			&answer.ID,
			&answer.UserID,
			&answer.FormID,
			&answer.QuestionID,
			&answer.Answer,
		)
		if err != nil {
			return nil, err
		}
		answers = append(answers, &answer)
	}

	question.Answers = answers

	return &question, nil
}

func (m *PostgresDBRepo) QuestionsByForm(id int) ([]*models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, form_id, question_text, type, extension, created_at, updated_at
				FROM questions
				WHERE form_id = $1
				ORDER BY id
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question

	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID,
			&question.FormID,
			&question.Question,
			&question.Type,
			&question.Extension,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		query = `
					SELECT id, question_id, option_text
					FROM options
					WHERE question_id = $1
				`

		questionRows, err := m.DB.QueryContext(ctx, query, question.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		defer questionRows.Close()

		var options []*models.Option
		for questionRows.Next() {
			var option models.Option
			err := questionRows.Scan(
				&option.ID,
				&option.QuestionID,
				&option.Option,
			)
			if err != nil {
				return nil, err
			}

			options = append(options, &option)
		}

		question.Options = options

		ext, err := m.GetQuestionExtension(question.ID)
		if err != nil {
			return nil, err
		}

		question.QuestionExtension = ext

		questions = append(questions, &question)
	}

	return questions, nil
}

func (m *PostgresDBRepo) InsertQuestion(question models.Question) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into questions (form_id, question_text, type, extension, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		question.FormID,
		question.Question,
		question.Type,
		question.Extension,
		question.CreatedAt,
		question.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateQuestion(question models.Question) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update questions set question_text = $1, type = $2, extension = $3, updated_at = $4
				where id = $5`

	_, err := m.DB.ExecContext(ctx, stmt,
		question.Question,
		question.Type,
		question.Extension,
		question.UpdatedAt,
		question.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteQuestion(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from questions where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) UpdateQuestionOptions(id int, options []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from options where question_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	for _, n := range options {
		stmt := `insert into options (question_id, option_text) values ($1, $2)`
		_, err := m.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostgresDBRepo) DeleteQuestionOptions(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from options where question_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetQuestionExtension(id int) (*models.Extension, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, question_id, followup_question_id, followup_option_value
				FROM questions_extension
				WHERE question_id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var ext models.Extension

	err := row.Scan(
		&ext.ID,
		&ext.QuestionID,
		&ext.FollowUpQuestion,
		&ext.FollowUpOption,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &ext, nil
}

func (m *PostgresDBRepo) UpdateQuestionExtension(extension *models.Extension) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from questions_extension where question_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, extension.QuestionID)
	if err != nil {
		return err
	}

	stmt = `insert into questions_extension (question_id, followup_question_id, followup_option_value)
			values ($1, $2, $3)`

	_, err = m.DB.ExecContext(ctx, stmt,
		extension.QuestionID,
		extension.FollowUpQuestion,
		extension.FollowUpOption,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteQuestionExtension(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from questions_extension where question_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertAnswers(answers []*models.Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	for _, answer := range answers {
		stmt := `insert into answers (user_id, form_id, question_id, answer_text) VALUES ($1, $2, $3, $4)`
		_, err := m.DB.ExecContext(ctx, stmt,
			answer.UserID,
			answer.FormID,
			answer.QuestionID,
			answer.Answer,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostgresDBRepo) GroupAnswersByQuestion(formID int, questionID int) ([]*models.GroupAnswer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT COUNT(id) as count_answers, answer_text, form_id, question_id
				FROM answers
				WHERE form_id = $1 AND question_id = $2
				GROUP BY answer_text, form_id, question_id`

	rows, err := m.DB.QueryContext(ctx, query, formID, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupAnswers []*models.GroupAnswer

	for rows.Next() {
		var groupAnswer models.GroupAnswer
		err := rows.Scan(
			&groupAnswer.Count,
			&groupAnswer.Answer,
			&groupAnswer.FormID,
			&groupAnswer.QuestionID,
		)
		if err != nil {
			return nil, err
		}

		groupAnswers = append(groupAnswers, &groupAnswer)
	}

	return groupAnswers, nil
}

func (m *PostgresDBRepo) GetAnswersByUser(id int) ([]*models.Answer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT id, user_id, form_id, question_id, answer_text
				FROM answers
				WHERE user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*models.Answer

	for rows.Next() {
		var answer models.Answer
		err := rows.Scan(
			&answer.ID,
			&answer.UserID,
			&answer.FormID,
			&answer.QuestionID,
			&answer.Answer,
		)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if err == sql.ErrNoRows {
			return nil, nil
		}

		answers = append(answers, &answer)
	}

	return answers, nil
}

func (m *PostgresDBRepo) AllForums() ([]*models.Forum, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT f.id, f.forum_text, f.user_id, f.published_at, u.username, u.is_admin
				FROM forums f
				LEFT JOIN users u ON f.user_id = u.id
				ORDER BY f.id DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []*models.Forum
	var isAdmin bool

	for rows.Next() {
		var forum models.Forum
		err := rows.Scan(
			&forum.ID,
			&forum.Forum,
			&forum.UserID,
			&forum.PublishedAt,
			&forum.UserUsername,
			&isAdmin,
		)
		if err != nil {
			return nil, err
		}

		likesNumber, err := m.GetForumLikesNumber(forum.ID)
		if err != nil {
			return nil, err
		}

		commentsNumber, err := m.GetForumCommentsNumber(forum.ID)
		if err != nil {
			return nil, err
		}

		comments, err := m.GetCommentsByForum(forum.ID)
		if err != nil {
			return nil, err
		}

		forum.LikesNumber = likesNumber
		forum.CommentsNumber = commentsNumber
		forum.Comments = comments

		if !isAdmin {
			alumniID, err := m.GetProfileByUserID(forum.UserID)
			if err != nil {
				return nil, err
			}

			name, err := m.GetAlumniNameByID(alumniID.AlumniID)
			if err != nil {
				return nil, err
			}

			photo, err := m.GetUserPhotoByID(forum.UserID)
			if err != nil {
				return nil, err
			}

			forum.UserName = name
			forum.UserPhoto = photo
		}

		forums = append(forums, &forum)
	}

	return forums, nil
}

func (m *PostgresDBRepo) Forum(id int) (*models.Forum, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, form_text, user_id, published_at
				FROM forums
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var forum models.Forum

	err := row.Scan(
		&forum.ID,
		&forum.Forum,
		&forum.UserID,
		&forum.PublishedAt,
	)

	if err != nil {
		return nil, err
	}

	comments, err := m.GetCommentsByForum(id)
	if err != nil {
		return nil, err
	}

	likes, err := m.GetLikesByForum(id)
	if err != nil {
		return nil, err
	}

	forum.Likes = likes
	forum.Comments = comments

	return &forum, nil
}

func (m *PostgresDBRepo) InsertForum(forum models.Forum) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into forums (forum_text, user_id, published_at)
			values ($1, $2, $3) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		forum.Forum,
		forum.UserID,
		forum.PublishedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) DeleteForum(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from forums where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetForumsByUser(id int) ([]*models.Forum, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT f.id, f.forum_text, f.user_id, f.published_at, u.username, u.is_admin
				FROM forums f
				LEFT JOIN users u ON f.user_id = u.id
				WHERE u.id = $1
				ORDER BY f.id DESC`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []*models.Forum
	var isAdmin bool

	for rows.Next() {
		var forum models.Forum
		err := rows.Scan(
			&forum.ID,
			&forum.Forum,
			&forum.UserID,
			&forum.PublishedAt,
			&forum.UserUsername,
			&isAdmin,
		)
		if err != nil {
			return nil, err
		}

		likesNumber, err := m.GetForumLikesNumber(forum.ID)
		if err != nil {
			return nil, err
		}

		commentsNumber, err := m.GetForumCommentsNumber(forum.ID)
		if err != nil {
			return nil, err
		}

		comments, err := m.GetCommentsByForum(forum.ID)
		if err != nil {
			return nil, err
		}

		forum.LikesNumber = likesNumber
		forum.CommentsNumber = commentsNumber
		forum.Comments = comments

		if !isAdmin {
			alumniID, err := m.GetProfileByUserID(forum.UserID)
			if err != nil {
				return nil, err
			}

			name, err := m.GetAlumniNameByID(alumniID.AlumniID)
			if err != nil {
				return nil, err
			}

			forum.UserName = name
		}

		forums = append(forums, &forum)
	}

	return forums, nil
}

func (m *PostgresDBRepo) GetForumLikesNumber(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT COUNT(id) as count_likes
				FROM likes
				WHERE forum_id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var likesNumber int

	err := row.Scan(
		&likesNumber,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return likesNumber, nil

}

func (m *PostgresDBRepo) GetForumCommentsNumber(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT COUNT(id) as count_replies
				FROM replies
				WHERE forum_id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var commentsNumber int

	err := row.Scan(
		&commentsNumber,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return commentsNumber, nil

}

func (m *PostgresDBRepo) InsertComment(comment models.Comment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into replies (forum_id, user_id, reply_text, published_at)
			values ($1, $2, $3, $4) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		comment.ForumID,
		comment.UserID,
		comment.Comment,
		comment.PublishedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) GetCommentsByForum(id int) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, forum_id, user_id, reply_text, published_at
				FROM replies
				WHERE forum_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ForumID,
			&comment.UserID,
			&comment.Comment,
			&comment.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		user, err := m.GetUserByID(comment.UserID)
		if err != nil {
			return nil, err
		}

		if !user.IsAdmin {
			profile, err := m.GetProfileByUserID(comment.UserID)
			if err != nil {
				return nil, err
			}

			name, err := m.GetAlumniNameByID(profile.AlumniID)
			if err != nil {
				return nil, err
			}

			comment.UserName = name
		} else {
			comment.UserName = ""
		}

		photo, err := m.GetUserPhotoByID(comment.UserID)
		if err != nil {
			return nil, err
		}

		comment.UserUsername = user.Username
		comment.UserPhoto = photo

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *PostgresDBRepo) InsertLike(like models.Like) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	// Query untuk mengecek apakah like sudah ada
	checkStmt := `SELECT COUNT(*) FROM likes WHERE forum_id = $1 AND user_id = $2`
	var count int
	err := m.DB.QueryRowContext(ctx, checkStmt, like.ForumID, like.UserID).Scan(&count)
	if err != nil {
		return err
	}

	// Jika sudah ada like dengan user_id dan forum_id yang sama, kembalikan error
	if count > 0 {
		return nil
	}

	// Query untuk menambahkan like
	insertStmt := `INSERT INTO likes (forum_id, user_id, created_at)
					VALUES ($1, $2, $3)`

	_, err = m.DB.ExecContext(ctx, insertStmt,
		like.ForumID,
		like.UserID,
		like.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteLike(userId int, forumId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from likes where user_id = $1 and forum_id = $2`

	_, err := m.DB.ExecContext(ctx, stmt, userId, forumId)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetLikesByUser(id int) ([]*models.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, forum_id, user_id, created_at
				FROM likes
				WHERE user_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var likes []*models.Like

	for rows.Next() {
		var like models.Like
		err := rows.Scan(
			&like.ID,
			&like.ForumID,
			&like.UserID,
			&like.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		likes = append(likes, &like)
	}

	return likes, nil
}

func (m *PostgresDBRepo) GetLikesByForum(id int) ([]*models.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, forum_id, user_id, created_at
				FROM likes
				WHERE forum_id = $1
			`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var likes []*models.Like

	for rows.Next() {
		var like models.Like
		err := rows.Scan(
			&like.ID,
			&like.ForumID,
			&like.UserID,
			&like.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		likes = append(likes, &like)
	}

	return likes, nil
}

func (m *PostgresDBRepo) AllJobs() ([]*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, job_position, company, job_location, job_type, min_salary, max_salary, closed, description, created_at, updated_at from jobs order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job

	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID,
			&job.UserID,
			&job.JobPosition,
			&job.Company,
			&job.JobLocation,
			&job.JobType,
			&job.MinSalary,
			&job.MaxSalary,
			&job.Closed,
			&job.Description,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (m *PostgresDBRepo) Job(id int) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, job_position, company, job_location, job_type, min_salary, max_salary, closed, description, created_at, updated_at from jobs where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var job models.Job

	err := row.Scan(
		&job.ID,
		&job.UserID,
		&job.JobPosition,
		&job.Company,
		&job.JobLocation,
		&job.JobType,
		&job.MinSalary,
		&job.MaxSalary,
		&job.Closed,
		&job.Description,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (m *PostgresDBRepo) InsertJob(job models.Job) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into jobs (user_id, job_position, company, job_location, job_type, min_salary, max_salary, description, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		job.UserID,
		job.JobPosition,
		job.Company,
		job.JobLocation,
		job.JobType,
		job.MinSalary,
		job.MaxSalary,
		job.Description,
		job.CreatedAt,
		job.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateJob(job models.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update jobs set job_position = $1, company = $2, job_location = $3, job_type = $4, min_salary = $5, max_salary = $6, description = $7, closed = $8, updated_at = $9 
			where id = $10`

	_, err := m.DB.ExecContext(ctx, stmt,
		job.JobPosition,
		job.Company,
		job.JobLocation,
		job.JobType,
		job.MinSalary,
		job.MaxSalary,
		job.Description,
		job.Closed,
		job.UpdatedAt,
		job.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteJob(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from jobs where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
