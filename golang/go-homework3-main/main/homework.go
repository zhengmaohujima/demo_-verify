package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 模型
type User struct {
	ID        uint
	Name      string
	Email     string
	PostCount int // 文章数量统计
	Posts     []Post
	CreatedAt time.Time
}

type Post struct {
	ID            uint
	Title         string
	Content       string
	UserID        uint
	CommentCount  int
	CommentStatus string // 有评论 / 无评论
	Comments      []Comment
	CreatedAt     time.Time
}

type Comment struct {
	ID        uint
	Content   string
	PostID    uint
	CreatedAt time.Time
}

func main() {
	// 连接数据库
	//db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}

	dsn := "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database tables created successfully!")
	QueryUserPostsWithComments(db, 1)
	QueryMostCommentedPost(db)
}

func QueryUserPostsWithComments(db *gorm.DB, userID uint) (*User, error) {
	var user User

	err := db.
		Preload("Posts").
		Preload("Posts.Comments").
		First(&user, userID).Error

	if err != nil {
		return nil, err
	}
	fmt.Println(user)
	return &user, nil
}

func QueryMostCommentedPost(db *gorm.DB) (*Post, error) {
	var post Post

	err := db.
		Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&post).Error

	if err != nil {
		return nil, err
	}
	fmt.Println(post)

	return &post, nil
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64

	// 查询该文章剩余评论数
	err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error
	if err != nil {
		return err
	}

	// 如果没有评论了，更新文章状态
	if count == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Updates(map[string]interface{}{
				"comment_count":  0,
				"comment_status": "无评论",
			}).Error
	}

	// 否则只更新数量
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_count", count).Error
}
