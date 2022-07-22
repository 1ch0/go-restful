package api

import "github.com/emicklei/go-restful/v3"

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

func InitAPIBean() []interface{} {
	// Authentication
	RegisterAPIInterface(NewAuthenticationAPIInterface())
	//RegisterAPIInterface(NewUserAPIInterface())

	var beans []interface{}
	for i := range registeredAPIInterface {
		beans = append(beans, registeredAPIInterface[i])
	}
	beans = append(beans, NewAuthenticationAPIInterface())
	return beans
}
