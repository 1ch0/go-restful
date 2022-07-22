package api

//import (
//	"github.com/1ch0/go-restful/pkg/apiserver/domain/service"
//	"github.com/emicklei/go-restful/v3"
//)
//
//type userAPIInterface struct {
//	UserService service.UserService `inject:""`
//}
//
//func NewUserAPiInterface() Interface {
//	return &userAPIInterface{}
//}
//
//func (c *userAPIInterface) GetWebServiceRoute() *restful.WebService {
//	ws := new(restful.WebService)
//	ws.Path(versionPerfix+"/users").
//		Consumes(restful.MIME_XML, restful.MIME_JSON).
//		Produces(restful.MIME_JSON, restful.MIME_XML).
//		Doc("api for user manage")
//
//	tags := []string{"user"}
//
//	ws.Route(ws.GET("/").To(c.listuser))
//}
