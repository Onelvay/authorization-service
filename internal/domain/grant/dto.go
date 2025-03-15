package grant

type Response struct {
	AccessToken string `json:"accessToken"`
}

type Request struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
