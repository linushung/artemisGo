package postgres

type Identity struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=30,min=6"`
}

type LoginReq struct{ Identity }

type RegisterReq struct {
	Identity
	Username string `json:"username" binding:"required,alphanum,min=3"`
}

type UpdateReq struct {
	/* Remove email from update request because Artemis use email as primary key for poster table */
	// Email    string `json:"email"`
	Username string `json:"username" binding:"omitempty,alphanum,min=3"`
	Password string `json:"password" binding:"omitempty,max=30,min=6"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

type ArticleReq struct {
	Title       string   `json:"title" binding:"required,max=20"`
	Description string   `json:"description" binding:"required,max=50"`
	Body        string   `json:"body" binding:"required,max=200"`
	Tags        []string `json:"tagList" binding:"omitempty,oneof=angularjs reactjs vuejs"`
}
