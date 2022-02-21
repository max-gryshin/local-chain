package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	ErrorAuthCheckTokenFail       = 20001
	ErrorAuthCheckTokenTimeout    = 20002
	ErrorAuthToken                = 20003
	ErrorAuth                     = 20004
	ErrorAuthCheckCredentialsFail = 20005

	ErrorUploadSaveImageFail    = 30001
	ErrorUploadCheckImageFail   = 30002
	ErrorUploadCheckImageFormat = 30003

	UserExists = 409
)
