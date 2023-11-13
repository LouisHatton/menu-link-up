package routes

const (
	FileIdParam     = "fileId"
	FileUrlParam    = "/{fileId}"
	FilesPathPrefix = "/files"
	CreateFilesPath = FilesPathPrefix
	FileIdPath      = FilesPathPrefix + FileUrlParam
	FileIdLinkPath  = FileIdPath + "/link"
)
