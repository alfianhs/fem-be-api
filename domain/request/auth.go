package request

type SuperadminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MemberRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MemberLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type ResendEmailVerificationRequest struct {
	Email string `json:"email"`
}
