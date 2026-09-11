package ozon

import (
	"context"
	"net/http"
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

type Roles struct {
	client *ozonCore.Client
}

// 获取当前密钥的角色/方法列表与密钥到期日期。
// 密钥有效期为 3 个月，到期后需在 OZON 后台生成新密钥。
// 可通过 expires_at 获取密钥到期日期。
type GetRolesResponse struct {
	ozonCore.CommonResponse

	// 密钥到期日期
	ExpiresAt time.Time `json:"expires_at"`

	// 可用角色和方法信息
	Roles []GetRolesRole `json:"roles"`
}

type GetRolesRole struct {
	// 角色名称
	Name string `json:"name"`

	// 该角色可用方法列表
	Methods []string `json:"methods"`
}

// 获取角色和方式列表。
func (c Roles) GetRoles(ctx context.Context) (*GetRolesResponse, error) {
	url := "/v1/roles"

	resp := &GetRolesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
