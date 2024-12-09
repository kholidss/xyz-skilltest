package presentation

type (
	ReqRegisterUser struct {
		NIK        string `json:"nik"`
		FullName   string `json:"full_name"`
		LegalName  string `json:"legal_name"`
		POB        string `json:"place_of_birth"`
		DOB        string `json:"date_of_birth"`
		Salary     int    `json:"salary"`
		Password   string `json:"password"`
		FileKTP    *File
		FileSelfie *File
	}

	RespRegisterUser struct {
		UserID    string `json:"user_id"`
		FullName  string `json:"full_name"`
		LegalName string `json:"legal_name"`
		DOB       string `json:"date_of_birth"`
	}
)

type (
	ReqLoginUser struct {
		NIK      string `json:"nik"`
		Password string `json:"password"`
	}

	RespLoginUser struct {
		UserID      string `json:"user_id"`
		AccessToken string `json:"access_token"`
		FullName    string `json:"full_name"`
		LegalName   string `json:"legal_name"`
	}
)

type UserAuthData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	FullName    string `json:"full_name"`
	LegalName   string `json:"legal_name"`
}
