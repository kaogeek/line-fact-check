package factcheck

type TypeMetadata string

const (
	TypeMetadataUserInfo = "META_USERINFO"
)

type Metadata[T any] struct {
	Type TypeMetadata `json:"type"`
	Data T            `json:"data"`
}

type UserInfo struct {
	UserType TypeUser `json:"user_type"`
	UserID   string   `json:"user_id"`
}
