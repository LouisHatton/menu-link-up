package api

const (
	FileIdParam    = "fileId"
	FileUrlParam   = "/{" + FileIdParam + "}"
	FilesPath      = "/files"
	FileIdPath     = FilesPath + FileUrlParam
	FileIdLinkPath = FileIdPath + "/link"
	CheckFilePath  = "/check-file"

	UserIdParam  = "userId"
	UserUrlParam = "/{" + UserIdParam + "}"
	UsersPath    = "/users"
	UserIdPath   = UsersPath + UserUrlParam
)
