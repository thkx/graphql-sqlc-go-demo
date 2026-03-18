package memory

import (
	"fmt"
	"sync"

	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	users        map[string]*model.User
	posts        map[string]*model.Post
	comments     map[string]*model.Comment
	indicators   map[string]*model.Indicator
	userIDSeq    int
	postIDSeq    int
	commentIDSeq int
}

func NewMemoryRepository() *MemoryRepository {
	r := &MemoryRepository{
		users:      make(map[string]*model.User),
		posts:      make(map[string]*model.Post),
		comments:   make(map[string]*model.Comment),
		indicators: make(map[string]*model.Indicator),
	}
	r.seedData()
	return r
}

func (r *MemoryRepository) Close() error { return nil }

func (r *MemoryRepository) seedData() {
	r.users["1"] = &model.User{ID: "1", Name: "test1", Email: "test1@qq.com", Password: "test1", Avatar: "/avatar/default.png", Gender: "m", Bio: "这是测试用户"}
	// Add seed posts
	r.posts["1"] = &model.Post{ID: "1", Title: "First post'", Content: "This is the content of the first post.", Pv: 0, Author: r.users["1"]}
	r.comments["1"] = &model.Comment{ID: "1", UserID: "1", Content: "这是一条测试评论", PostID: "1"}
	r.userIDSeq = 1
	r.postIDSeq = 1
	r.commentIDSeq = 1
}

func (r *MemoryRepository) ListUsers() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

func (r *MemoryRepository) GetUserByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user := r.users[id]
	return user, nil
}

func (r *MemoryRepository) GetUserByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var user model.User
	for _, v := range r.users {
		if v.GetEmail() == email {
			user = *v
		}
	}
	return &user, nil
}

func (r *MemoryRepository) CreateUser(user *model.User) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.userIDSeq = r.userIDSeq + 1
	r.users[fmt.Sprintf("%d", r.userIDSeq)] = user
	return nil
}

func (r *MemoryRepository) GetUserByName(name string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var user model.User
	for _, v := range r.users {
		if v.GetName() == name {
			user = *v
		}
	}
	return &user, nil
}

func (r *MemoryRepository) UpdateUser(user *model.User) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var row model.User
	for i, v := range r.users {
		if v.GetID() == user.ID {
			r.users[i] = user
		}
	}
	return &row, nil
}

func (r *MemoryRepository) ListPosts() ([]*model.Post, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()
	posts := make([]*model.Post, 0, len(r.posts))
	for _, u := range r.posts {
		posts = append(posts, u)
	}
	return posts, nil
}

func (r *MemoryRepository) CreatePost(post *model.Post) error {

	return nil
}

func (r *MemoryRepository) GetPostByID(id string) (*model.Post, error) {

	return nil, nil
}

func (r *MemoryRepository) ListPostsByTitle(title string) ([]*model.Post, error) {

	return nil, nil
}

func (r *MemoryRepository) UpdatePost(post *model.Post) (*model.Post, error) {

	return &model.Post{ID: post.ID}, nil
}

func (r *MemoryRepository) DeletePost(id string) error {
	return nil
}

func (r *MemoryRepository) ListComments() ([]*model.Comment, error) {

	return nil, nil
}

func (r *MemoryRepository) GetCommentByID(id string) (*model.Comment, error) {

	return nil, nil
}

func (r *MemoryRepository) ListCommentsByUserID(user_id string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MemoryRepository) ListCommentsByPostID(post_id string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MemoryRepository) ListCommentsByContent(content string) ([]*model.Comment, error) {

	return nil, nil
}

func (r *MemoryRepository) CreateComment(comment *model.Comment) error {

	return nil
}

func (r *MemoryRepository) UpdateComment(comment *model.Comment) (*model.Comment, error) {

	return &model.Comment{ID: comment.ID}, nil
}

func (r *MemoryRepository) DeleteComment(id string) error {

	return nil
}

func (r *MemoryRepository) CreateIndicator(ind *model.Indicator) error {

	return nil
}

func (r *MemoryRepository) ListIndicators() ([]*model.Indicator, error) {

	return nil, nil
}
