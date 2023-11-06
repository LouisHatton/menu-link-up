package routes

const (
	ProjectIdParam    = "projectId"
	ProjectUrlParam   = "/{projectId}"
	ProjectPathPrefix = "/projects"
	CreateProjectPath = ProjectPathPrefix
	CheckProjectPath  = ProjectPathPrefix + "/check"
	ProjectIdPath     = ProjectPathPrefix + ProjectUrlParam
)

const (
	FileIdParam     = "fileId"
	FileUrlParam    = "/{fileId}"
	FilesPathPrefix = "/files"
	CreateFilesPath = ProjectIdPath + FilesPathPrefix
	FilesListPath   = ProjectIdPath + FilesPathPrefix + "/list"
	FileIdPath      = ProjectIdPath + FilesPathPrefix + FileUrlParam
)
