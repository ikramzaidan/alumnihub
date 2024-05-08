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

func (m *PostgresDBRepo) AllAlumni() ([]*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, name, date_of_birth, place_of_birth, gender, phone from alumni order by id`

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
			&alumni.Name,
			&alumni.BirthDate,
			&alumni.BirthPlace,
			&alumni.Gender,
			&alumni.Phone,
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
				SELECT id, name, date_of_birth, place_of_birth, gender, phone
				FROM alumni
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var alumni models.Alumni

	err := row.Scan(
		&alumni.ID,
		&alumni.Name,
		&alumni.BirthDate,
		&alumni.BirthPlace,
		&alumni.Gender,
		&alumni.Phone,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
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
