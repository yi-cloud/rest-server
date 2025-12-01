package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/yi-cloud/rest-server/pkg/server"
	"net/http"
)

func GetRoles(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	page, err := roles.List(client, roles.ListOpts{
		Name:    c.Query("name"),
		Filters: c.QueryMap("filters"),
	}).AllPages(context.TODO())

	var ret any
	if page != nil {
		if ok, _ := page.IsEmpty(); !ok {
			ret = page.GetBody()
		}
	}
	GinResponseData(c, ret, err)
}

func GetRole(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	var ret any
	url := client.ServiceURL("roles", c.Param("id"))
	resp, err := client.Get(context.TODO(), url, &ret, nil)
	GinResponseData(c, ret, err, resp.StatusCode)
}

func postRoleProcess(c *gin.Context, body any, parts ...string) {
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	var ret any
	url := client.ServiceURL(parts...)
	resp, err := client.Post(context.TODO(), url, map[string]any{"role": body}, &ret, &gophercloud.RequestOpts{
		OkCodes: []int{201, 204},
	})
	GinResponseData(c, ret, err, resp.StatusCode)
}

func CreateRole(c *gin.Context) {
	var body Role
	postRoleProcess(c, body, "roles")
}

func UpdateRole(c *gin.Context) {
	var body UpdateRoleOpts
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	var ret any
	url := client.ServiceURL("roles", c.Param("id"))
	resp, err := client.Patch(context.TODO(), url, map[string]any{"role": body}, &ret, nil)
	GinResponseData(c, ret, err, resp.StatusCode)
}

func DeleteRole(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	url := client.ServiceURL("roles", c.Param("id"))
	resp, err := client.Delete(context.TODO(), url, nil)
	GinResponseData(c, nil, err, resp.StatusCode)
}

func AddRoleRoutes(s *server.ApiServer) {
	g := s.AddGroup("/roles", nil)
	s.AddRoute(g, "GET", "", GetRoles)
	s.AddRoute(g, "GET", "/:id", GetRole)
	s.AddRoute(g, "POST", "", CreateRole)
	s.AddRoute(g, "PATCH", "/:id", UpdateRole)
	s.AddRoute(g, "DELETE", "/:id", DeleteRole)
}
