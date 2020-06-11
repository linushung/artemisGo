package postgres

type role string

const (
	Admin   role = "ADMIN"
	User    role = "USER"
	Visitor role = "VISITOR"
	Unknown role = "UNKNOWN"
)

/* Ref: https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/ */

type Poster struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"-"`
	Image    string `json:"image omitempty"`
	Bio      string `json:"bio omitempty"`
	Token    string `json:"-"`
}

type Follower struct {
	Email    string `json:"email"`
	Follower string `json:"follower"`
}

type Profile struct {
	Username  string `json:"username"`
	Image     string `json:"image omitempty"`
	Bio       string `json:"bio omitempty"`
	Following bool   `json:"following"`
}

type Article struct {
	Slug           string   `json:"slug" db:"slug"`
	Title          string   `json:"title" db:"title"`
	Description    string   `json:"description" db:"description"`
	Body           string   `json:"body" db:"body"`
	CreateTime     int64    `json:"createdAt" db:"created_time"`
	UpdateTime     int64    `json:"updatedAt" db:"modified_time"`
	Favorite       bool     `json:"favorited" db:"favorite"`
	FavoritesCount bool     `json:"favoritesCount" db:"favorite_count"`
	TagId          int64    `json:"-" db:"tagId"`
	Tags           []string `json:"tagList"`
	Author         Profile  `json:"author"`
}
