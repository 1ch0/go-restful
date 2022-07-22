package api

import "github.com/emicklei/go-restful/v3"

type Interface interface {
	GetWebServiceRoute() *restful.WebService
}

var registeredAPIInterface []Interface

func GetRegisterAPIInterface() []Interface {
	return registeredAPIInterface
}
