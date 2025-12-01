package handlers

type User struct {
	Name        string  `json:"name" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateUserOpts struct {
	Name        string  `json:"name" binding:"required"`
	Password    *string `json:"password,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ChangePassword struct {
	Password         string `json:"password" binding:"required"`
	OriginalPassword string `json:"original_password" binding:"required"`
}

type Role struct {
	Name        string         `json:"name" binding:"required"`
	Options     map[string]any `json:"options,omitempty"`
	Description *string        `json:"description,omitempty"`
}

type UpdateRoleOpts struct {
	Name        *string        `json:"name,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
	Description *string        `json:"description,omitempty"`
}

type TokenResponse struct {
	Catalogs  []any  `json:"catalogs"`
	ExpiresAt any    `json:"expires_at"`
	IsAdmin   bool   `json:"is_admin"`
	IssuedAt  any    `json:"issued_at"`
	Project   any    `json:"project"`
	Roles     []any  `json:"-"`
	User      any    `json:"user"`
	Token     string `json:"token,omitempty"`
}
