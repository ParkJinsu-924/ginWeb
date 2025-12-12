package common

const ( // for session cookie
	SessionUserIdKey       = "userId"
	SessionUserNicknameKey = "name"
)

const (
	SignupEndpoint = "/signup"
	LoginEndpoint  = "/login"
)

const (
	HomeEndpoint       = "/"
	LogoutEndpoint     = "/logout"
	PostFormEndpoint   = "/posts/form"
	PostCreateEndpoint = "/posts/create"
	PostDetailEndpoint = "/posts/:id"
	PostDeleteEndpoint = "/posts/delete"
)
