package api

const (
	FileIdParam    = "fileId"
	FileUrlParam   = "/{" + FileIdParam + "}"
	FilesPath      = "/files"
	FileIdPath     = FilesPath + FileUrlParam
	FileIdLinkPath = FileIdPath + "/link"

	UserIdParam  = "userId"
	UserUrlParam = "/{" + UserIdParam + "}"
	UsersPath    = "/users"
	UserIdPath   = UsersPath + UserUrlParam
)
