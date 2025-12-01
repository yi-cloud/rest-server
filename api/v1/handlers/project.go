package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/yi-cloud/rest-server/pkg/server"
)

func ListAssignToUserOnProject(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	page, err := roles.ListAssignmentsOnResource(client, roles.ListAssignmentsOnResourceOpts{
		ProjectID: c.Param("id"),
		UserID:    c.Param("user_id"),
	}).AllPages(context.TODO())

	var ret any
	if page != nil {
		if ok, _ := page.IsEmpty(); !ok {
			ret = page.GetBody()
		}
	}

	GinResponseData(c, ret, err)
}

func AssignRoleToUserOnProject(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	r := roles.Assign(context.TODO(), client, c.Param("role_id"), roles.AssignOpts{
		ProjectID: c.Param("id"),
		UserID:    c.Param("user_id"),
	})
	GinResponseData(c, r.Body, r.Err, r.StatusCode)
}

func UnassignRoleToUserOnProject(c *gin.Context) {
	client, err := NewIdentityV3Client(c)
	if err != nil {
		GinResponseData(c, nil, err)
		return
	}

	r := roles.Unassign(context.TODO(), client, c.Param("role_id"), roles.UnassignOpts{
		ProjectID: c.Param("id"),
		UserID:    c.Param("user_id"),
	})
	GinResponseData(c, r.Body, r.Err, r.StatusCode)
}

func AddProjectRoutes(s *server.ApiServer) {
	g := s.AddGroup("/projects", nil)
	s.AddRoute(g, "GET", "/:id/users/:user_id/roles", ListAssignToUserOnProject)
	s.AddRoute(g, "PUT", "/:id/users/:user_id/roles/:role_id", AssignRoleToUserOnProject)
	s.AddRoute(g, "DELETE", "/:id/users/:user_id/roles/:role_id", UnassignRoleToUserOnProject)
}
