package api

import (
	"net/http"

	apisv1 "github.com/1ch0/go-restful/pkg/apiserver/interface/api/dto/v1"
	"github.com/emicklei/go-restful/v3"
)

var versionPerfix = "/api/v1"

type Interface interface {
	GetWebServiceRoute() *restful.WebService
}

var registeredAPIInterface []Interface

func RegisterAPIInterface(ws Interface) {
	registeredAPIInterface = append(registeredAPIInterface, ws)
}

func GetRegisterAPIInterface() []Interface {
	return registeredAPIInterface
}

func returns200(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", apisv1.SimpleResponse{Status: "ok"})
}

func returns500(b *restful.RouteBuilder) {
	b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}

func InitAPIBean() []interface{} {
	// Authentication
	RegisterAPIInterface(NewAuthenticationAPIInterface())
	RegisterAPIInterface(NewUserAPIInterface())

	var beans []interface{}
	for i := range registeredAPIInterface {
		beans = append(beans, registeredAPIInterface[i])
	}
	return beans
}
