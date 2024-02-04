package forms

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewsForm struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TagForm struct {
	Tag   string `json:"tag"`
	Color string `json:"color"`
}
