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

func (m *PostgresDBRepo) GetProfileByAlumniID(id int) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, alumni_id, user_id 
			from alumni_profile where alumni_id = $1`

	var profile models.Profile
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&profile.ID,
		&profile.AlumniID,
		&profile.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil

}

func (m *PostgresDBRepo) GetProfileByUserID(id int) (*models.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, alumni_id, user_id 
			from alumni_profile where user_id = $1`

	var profile models.Profile
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&profile.ID,
		&profile.AlumniID,
		&profile.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil

}

func (m *PostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, title, slug, body, status, created_at, updated_at, published_at from articles order by id`

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
				SELECT id, title, slug, body, status, created_at, updated_at, published_at
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

	stmt := `insert into articles (title, slug, body, status, created_at, updated_at, published_at)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
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
				status = $4, updated_at = $5, published_at = $6 
				where id = $7`

	_, err := m.DB.ExecContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
		article.Status,
		article.UpdatedAt,
		article.PublishedAt,
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

	query := `select id, title, description, has_time_limit, start_date, end_date, created_at, updated_at 
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
				SELECT id, title, description, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
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
				SELECT id, question_text, type, created_at, updated_at
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
				SELECT id, title, description, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
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

func (m *PostgresDBRepo) ShowFormAnswers(id int) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, description, has_time_limit, start_date, end_date, created_at, updated_at
				FROM forms
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var form models.Form

	err := row.Scan(
		&form.ID,
		&form.Title,
		&form.Description,
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
				SELECT id, question_text, type, created_at, updated_at
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
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

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

		questions = append(questions, &question)
	}

	form.Questions = questions

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

	stmt := `update forms set title = $1, description = $2, has_time_limit = $3, start_date = $4, end_date = $5, updated_at = $6
				where id = $7`

	_, err := m.DB.ExecContext(ctx, stmt,
		form.Title,
		form.Description,
		form.HasTimeLimit,
		form.StartDate,
		form.EndDate,
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

func (m *PostgresDBRepo) Question(id int) (*models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, form_id, question_text, type, created_at, updated_at
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

	return &question, nil
}

func (m *PostgresDBRepo) QuestionsByForm(id int) ([]*models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, form_id, question_text, type, created_at, updated_at
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

		questions = append(questions, &question)
	}

	return questions, nil
}

func (m *PostgresDBRepo) InsertQuestion(question models.Question) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into questions (form_id, question_text, type, created_at, updated_at)
			values ($1, $2, $3, $4, $5) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		question.FormID,
		question.Question,
		question.Type,
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

	stmt := `update questions set question_text = $1, type = $2, updated_at = $3
				where id = $4`

	_, err := m.DB.ExecContext(ctx, stmt,
		question.Question,
		question.Type,
		question.UpdatedAt,
		question.ID,
	)

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
			isAdmin,
		)
		if err != nil {
			return nil, err
		}

		if isAdmin == false {
			name, err := m.GetAlumniNameByID(forum.UserID)
			if err != nil {
				return nil, err
			}

			forum.UserName = name
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

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

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

func (m *PostgresDBRepo) InsertComment(comment models.Comment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into replies (id, forum_id, user_id, reply_text, published_at)
			values ($1, $2, $3, $4, $5) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		comment.ID,
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

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (m *PostgresDBRepo) InsertLike(like models.Like) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into likes (forum_id, user_id, created_at)
			values ($1, $2, $3)`

	_, err := m.DB.ExecContext(ctx, stmt,
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
