package forms

type Login struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Username   string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Password   string `json:"password" binding:"required,min=12,max=75"`
	Name       string `json:"name" binding:"required,min=3,max=32,alphanum"`
	Email      string `json:"email" binding:"required,email"`
	PublicKey  string `json:"public_key" binding:"required"`
	PrivateKey string `json:"private_key" binding:"required"`
}
