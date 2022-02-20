package e

var MsgFlags = map[int]string{
	Success:                       "ok",
	Error:                         "fail",
	InvalidParams:                 "invalid params",
	ErrorAuthCheckTokenFail:       "Token Authentication failed",
	ErrorAuthCheckTokenTimeout:    "Token Timed out",
	ErrorAuthToken:                "Token authentication failed",
	ErrorAuth:                     "Token error",
	ErrorAuthCheckCredentialsFail: "Credentials authentication failed",
	ErrorUploadSaveImageFail:      "Failed to save picture",
	ErrorUploadCheckImageFail:     "Check image failed",
	ErrorUploadCheckImageFormat:   "Check picture error, picture format or size problem",
	UserExists:                    "user already exists",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
