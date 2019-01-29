package response

type (
	Data struct {
		Data  interface{} `json:"data"`
		State bool        `json:"state"`
	}
	SendActiveCode struct {
		PhoneKey string `json:"phone_key"`
	}
	ValidationError struct {
		Error string `json:"error"`
	}
	ServerError struct {
		Message string `json:"message"`
	}
	UnauthorizedError struct {
		Error string `json:"error"`
	}
	DatabaseError struct {
		Message string `json:"message"`
	}
	CheckedUsername struct {
		Message string `json:"message"`
	}
	NewToken struct {
		Token     string `json:"token"`
		IsNewUser bool   `json:"is_new_user"`
	}
	Token struct {
		Token string `json:"token"`
	}
	NewUser struct {
		PhoneKey  string `json:"phone_key"`
		IsNewUser bool   `json:"is_new_user"`
	}
	CreateDocument struct {
		Message string `json:"message"`
		Key     string `json:"key"`
	}
	UpdateDocument struct {
		Message string `json:"message"`
		Key     string `json:"key"`
	}
	FindDocument struct {
		Document interface{} `json:"document"`
	}
	FindAllDocuments struct {
		Documents   []map[string]interface{} `json:"documents"`
		TotalPages  float64                  `json:"total_pages"`
		CurrentPage int64                    `json:"current_page"`
	}
	DestroyDocument struct {
		Message string `json:"message"`
	}
	EmptyDocument struct {
		Message string `json:"message"`
	}
)
