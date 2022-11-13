// Filename : internal/data/forums.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"universityforum.miguelavila.net/internals/validator"
)

type Forum struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Version     int32     `json:"version"`
}

// define a ForumModel object that wraps a sql.DB connection pool
type ForumModel struct {
	DB *sql.DB
}

func ValidateForum(v *validator.Validator, forum *Forum) {
	// Use the Check() method to execute our validation checks
	v.Check(forum.Title != "", "title", "must be provided")
	v.Check(len(forum.Title) <= 200, "title", "must not be more than 200 bytes long")

	v.Check(forum.Description != "", "description", "must be provided")
	v.Check(len(forum.Description) <= 2000, "description", "must not be more than 2000 bytes long")
}

// Insert() allows us  to create a new Forum
func (m ForumModel) Insert(forum *Forum) error {
	query := `
		INSERT INTO forums (title, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`

	// Collect the data fields into a slice
	args := []interface{}{
		forum.Title, forum.Description,
	}
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&forum.ID, &forum.CreatedAt, &forum.Version)

}

// Get() allows us to retrieve a specific Forum
func (m ForumModel) Get(id int64) (*Forum, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, title, description, version
		FROM forums
		WHERE id = $1
	`
	// Declare a Forum variable to hold the returned data
	var forum Forum
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&forum.ID,
		&forum.CreatedAt,
		&forum.Title,
		&forum.Description,
		&forum.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &forum, nil
}

// Update() allows us to edit/alter a specific Forum
// Optimistic locking (version number)
func (m ForumModel) Update(forum *Forum) error {
	// Create the query
	query := `
		UPDATE forums
		SET title = $1, description = $2, version = version + 1
		WHERE id = $3
		AND version = $4
		RETURNING version
	`
	args := []interface{}{
		forum.Title,
		forum.Description,
		forum.ID,
		forum.Version,
	}

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&forum.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific Forum
func (m ForumModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM forums
		WHERE id = $1
	`

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// The GetAll() method retuns a list of all the forums sorted by id
func (m ForumModel) GetAll(title string, filters Filters) ([]*Forum, Metadata, error) {
	// Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, title, description,
		       version
		FROM forums
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortOrder())

	// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{title, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the resultset
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Forum data
	forums := []*Forum{}
	// Iterate over the rows in the resultset
	for rows.Next() {
		var forum Forum
		// Scan the values from the row into forum
		err := rows.Scan(
			&totalRecords,
			&forum.ID,
			&forum.CreatedAt,
			&forum.Title,
			&forum.Description,
			&forum.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Forum to our slice
		forums = append(forums, &forum)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of Forums
	return forums, metadata, nil
}
