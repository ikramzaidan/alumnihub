package repository

import (
	"alumnihub/internal/models"
	"context"
	"database/sql"
)

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
