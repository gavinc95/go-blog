package db

import (
	"database/sql"
	"fmt"

	"github.com/gavinc95/go-blog/db/models"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

// Wrapper interface that handles all blog-related operations
type BlogStore interface {
	UserStore
	PostStore
	GetDB() *sql.DB // used for table creation/deletion
}

type store struct {
	db        *sql.DB
	idManager IDManager
}

type IDManager interface {
	UUID() string
}

type GenID struct{}

func (g *GenID) UUID() string {
	return uuid.New().String()
}

func NewBlogStore(db *sql.DB, idManager IDManager) *store {
	return &store{
		db:        db,
		idManager: idManager,
	}
}

func (m *store) GetDB() *sql.DB {
	return m.db
}

// a sub-interface that handles only user-related operations
type UserStore interface {
	//GetAllUsers() ([]*models.User, error)
	GetUser(id string) (*models.User, error)
	CreateUser(name, email string) (string, error)
	UpdateUser(id, name, email string) (string, error)
	DeleteUser(id string) (string, error)
}

// a sub-interface that handles only post-related operations
type PostStore interface {
	GetAllPosts(userID string) ([]*models.Post, error)
	GetPost(postID string) (*models.Post, error)
	CreatePost(userID, title, content string) (string, error)
	UpdatePost(postID, title, content string) (string, error)
	DeletePost(postID string) (string, error)
}

func (m *store) GetUser(id string) (*models.User, error) {
	row := m.db.QueryRow("SELECT * FROM users WHERE id = $1", id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, xerrors.Errorf("error finding user in db: %w", err)
	}

	return &user, nil
}

func (m *store) CreateUser(name, email string) (string, error) {
	id := m.idManager.UUID()
	// create a new user row
	_, err := m.db.Exec("INSERT INTO users(id, name, email) VALUES($1, $2, $3)",
		id, name, email)
	if err != nil {
		return id, xerrors.Errorf("error while inserting user: %w", err)
	}
	return id, nil
}

func (m *store) UpdateUser(id, name, email string) (string, error) {
	// check if the user ID already exists in the db
	user, err := m.GetUser(id)
	if err != nil {
		return id, xerrors.Errorf("failed to check for existing user: %w", err)
	}
	if user == nil {
		return id, fmt.Errorf("user doesn't exist - create one first")
	}

	// update the existing user
	if name != "" {
		_, err := m.db.Exec("UPDATE users SET name = $1 WHERE id = $2",
			name, id)
		if err != nil {
			return id, xerrors.Errorf("error while updating user: %w", err)
		}
	}

	if email != "" {
		_, err := m.db.Exec("UPDATE users SET email = $1 WHERE id = $2",
			email, id)
		if err != nil {
			return id, xerrors.Errorf("error while updating user: %w", err)
		}
	}

	return id, nil
}

func (m *store) DeleteUser(id string) (string, error) {
	// check if the user exists
	user, err := m.GetUser(id)
	if err != nil {
		return id, err
	}
	if user == nil {
		return id, fmt.Errorf("user does not exist for ID: %s", id)
	}

	_, err = m.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return id, xerrors.Errorf("error deleting user: %w", err)
	}

	return id, nil
}

func (m *store) GetAllPosts(userID string) ([]*models.Post, error) {
	rows, err := m.db.Query("SELECT * FROM posts WHERE user_id = $1", userID)
	if err != nil {
		return nil, xerrors.Errorf("failed to fetch posts for user: %w", err)
	}

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
		if err != nil {
			return nil, xerrors.Errorf("error parsing DB response: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (m *store) GetPost(postID string) (*models.Post, error) {
	row := m.db.QueryRow("SELECT * FROM posts WHERE id = $1", postID)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, xerrors.Errorf("error finding post in db: %w", err)
	}

	return &post, nil
}

func (m *store) CreatePost(userID, title, content string) (string, error) {
	postID := m.idManager.UUID()

	// create the post
	_, err := m.db.Exec("INSERT INTO posts(id, user_id, title, content) VALUES($1, $2, $3, $4)",
		postID, userID, title, content)
	if err != nil {
		return postID, xerrors.Errorf("error creating new post: %w", err)
	}
	return postID, nil
}

func (m *store) UpdatePost(postID, title, content string) (string, error) {
	// check to see if a post with the same postID already exists
	// NOTE call the helper function instead of m.GetPost(...) to avoid checking for an existing user again
	post, err := m.GetPost(postID)
	if err != nil {
		return postID, xerrors.Errorf("error getting post: %w", err)
	}
	if post == nil {
		return postID, fmt.Errorf("post doesn't exist for ID: %s", postID)
	}

	// update the existing post
	if title != "" {
		_, err = m.db.Exec("UPDATE posts SET title = $1 WHERE id = $2",
			title, postID)
		if err != nil {
			return postID, xerrors.Errorf("error while updating post: %w", err)
		}
	}

	if content != "" {
		_, err = m.db.Exec("UPDATE posts SET content = $1 WHERE id = $2",
			content, postID)
		if err != nil {
			return postID, xerrors.Errorf("error while updating post: %w", err)
		}
	}

	return postID, nil
}

func (m *store) DeletePost(postID string) (string, error) {
	// check if the post exists
	post, err := m.GetPost(postID)
	if err != nil {
		return postID, xerrors.Errorf("error getting post: %w", err)
	}
	if post == nil {
		return postID, fmt.Errorf("cannot delete post that doesn't exist")
	}

	_, err = m.db.Exec("DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		return postID, xerrors.Errorf("error deleting post: %w", err)
	}

	return postID, nil
}
