package forms

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewsForm struct {
	Content string `json:"content"`
}
