package common

const (
	UserTable    = "users"
	PostTable    = "posts"
	CommentTable = "comments"
)

const ( // for session cookie
	SessionUserUUIDKey     = "uuid"
	SessionUserIdKey       = "userId"
	SessionUserNicknameKey = "name"
	SessionUserTagKey      = "tag"
)

const (
	SignupEndpoint = "/signup"
	LoginEndpoint  = "/login"
)

const (
	HomeEndpoint               = "/"
	LogoutEndpoint             = "/logout"
	PostFormEndpoint           = "/posts/form"
	PostCreateEndpoint         = "/posts/create"
	PostDetailEndpoint         = "/posts"
	PostDeleteEndpoint         = "/posts/delete"
	PostCommentsCreateEndpoint = "/posts/comments/create"
	PostCommentsUpdateEndpoint = "/posts/comments/update"
	PostCommentsDeleteEndpoint = "/posts/comments/delete"
)

const (
	ChatEndpoint = "/chat"
)
