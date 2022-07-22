package api

import (
	"github.com/1ch0/go-restful/pkg/apiserver/utils/bcode"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/1ch0/go-restful/pkg/apiserver/domain/service"
	apis "github.com/1ch0/go-restful/pkg/apiserver/interface/api/dto/v1"
)

type authenticationAPIInterface struct {
	AuthentcationService service.AuthenticationService
}

func NewAuthenticationAPIInterface() Interface {
	return &authenticationAPIInterface{}
}

func (c *authenticationAPIInterface) GetWebServiceRoute() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(versionPerfix+"/auth").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("api for authentication manage")

	tags := []string{"authentication"}

	ws.Route(ws.POST("/login").To(c.login).
		Doc("hanle login request").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(apis.LoginRequest{}).
		Returns(200, "", apis.LoginResponse{}).
		Returns(400, "", bcode.Bcode{}).
		Writes(apis.LoginResponse{}))

	return ws
}

func (c *authenticationAPIInterface) login(req *restful.Request, res *restful.Response) {
	var loginReq apis.LoginRequest
	if err := req.ReadEntity(&loginReq); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
	base, err := c.AuthentcationService.Login(req.Request.Context(), loginReq)
	if err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
	if err := res.WriteEntity(base); err != nil {
		bcode.ReturnError(req, res, err)
		return
	}
}
