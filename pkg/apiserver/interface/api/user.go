package api

import (
	"github.com/1ch0/go-restful/pkg/apiserver/domain/service"
	apisv1 "github.com/1ch0/go-restful/pkg/apiserver/interface/api/dto/v1"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/bcode"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

type userAPIInterface struct {
	UserService service.UserService `inject:""`
}

func NewUserAPIInterface() Interface {
	return &userAPIInterface{}
}

func (c *userAPIInterface) GetWebServiceRoute() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(versionPerfix+"/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("api for user manage")

	tags := []string{"user"}

	ws.Route(ws.GET("/").To(c.listUser).
		Doc("list users").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		//Filter(c.RbacService.CheckPerm("user", "list")).
		Param(ws.QueryParameter("page", "query the page number").DataType("integer")).
		Param(ws.QueryParameter("pageSize", "query the page size number").DataType("integer")).
		Param(ws.QueryParameter("name", "fuzzy search based on name").DataType("string")).
		Param(ws.QueryParameter("email", "fuzzy search based on email").DataType("string")).
		Param(ws.QueryParameter("alias", "fuzzy search based on alias").DataType("string")).
		Returns(200, "OK", apisv1.ListUserResponse{}).
		Returns(400, "Bad Request", bcode.Bcode{}).
		Writes(apisv1.ListUserResponse{}))

	return ws
}

func (c *userAPIInterface) listUser(req *restful.Request, res *restful.Response) {
	//page, pageSize, err := utils.ExtractPagingParams(req, minPageSize, maxPageSize)
	//if err != nil {
	//	bcode.ReturnError(req, res, err)
	//	return
	//}
	//resp, err := c.UserService.ListUsers(req.Request.Context(), page, pageSize, apis.ListUserOptions{
	//	Name:  req.QueryParameter("name"),
	//	Alias: req.QueryParameter("alias"),
	//	Email: req.QueryParameter("email"),
	//})
	//if err != nil {
	//	bcode.ReturnError(req, res, err)
	//	return
	//}
	//if err := res.WriteEntity(resp); err != nil {
	//	bcode.ReturnError(req, res, err)
	//	return
	//}
	return
}
