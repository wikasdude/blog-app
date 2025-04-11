package config

// API messages
const (
	MessagePostCreated   = "Post created successfully"
	MessagePostUpdated   = "Post updated successfully"
	MessagePostDeleted   = "Post deleted successfully"
	MessagePostFetched   = "Post fetched successfully"
	MessagePostsFetched  = "Posts fetched successfully"
	MessagePostNotFound  = "Post not found with ID: %d"
	MessageInvalidID     = "Invalid ID"
	MessageInvalidURL    = "Invalid URL"
	MessageInvalidMethod = "Method not allowed"
	MessageInternalError = "Internal server error"
)

// Default pagination and sorting
const (
	DefaultPage  = 1
	DefaultLimit = 10
	DefaultSort  = "created_at"
	DefaultOrder = "desc"
)

// Valid sort fields
var ValidSortFields = map[string]bool{
	"title":      true,
	"created_at": true,
}
