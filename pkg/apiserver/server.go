package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/1ch0/go-restful/pkg/apiserver/domain/service"

	"github.com/1ch0/go-restful/pkg/apiserver/interface/api"
	restfulSpec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"

	"github.com/1ch0/go-restful/pkg/apiserver/config"
	"github.com/1ch0/go-restful/pkg/apiserver/utils"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/container"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/log"
	pkgUtils "github.com/1ch0/go-restful/pkg/utils"
)

type APIServer interface {
	Run(context.Context, chan error) error
	BuildRestfulConfig() (*restfulSpec.Config, error)
}

type restServer struct {
	webContainer  *restful.Container
	beanContainer *container.Container
	cfg           config.Config
}

func New(cfg config.Config) (a APIServer) {
	s := &restServer{
		webContainer:  restful.NewContainer(),
		beanContainer: container.NewContainer(),
		cfg:           cfg,
	}
	return s
}

func (s *restServer) buildIoCContainer() error {
	// domain
	if err := s.beanContainer.Provides(service.InitServiceBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the service bean to the container: %w", err)
	}

	// interfaces
	if err := s.beanContainer.Provides(api.InitAPIBean()...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	return nil
}

func (s *restServer) Run(ctx context.Context, errChan chan error) error {

	s.RegisterAPIRoute()

	return s.startHTTP(ctx)
}

func (s *restServer) BuildRestfulConfig() (*restfulSpec.Config, error) {
	if err := s.buildIoCContainer(); err != nil {
		return nil, err
	}
	config := s.RegisterAPIRoute()
	return &config, nil
}

func (s *restServer) RegisterAPIRoute() restfulSpec.Config {
	/* **************************************************************  */
	/* *************       Open API Route Group     *****************  */
	/* **************************************************************  */
	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization", "RefreshToken"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		CookiesAllowed: true,
		Container:      s.webContainer,
	}
	s.webContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	s.webContainer.Filter(s.webContainer.OPTIONSFilter)
	s.webContainer.Filter(s.OPTIONSFilter)

	// Add request log
	s.webContainer.Filter(s.requestLog)

	// Register all custom api
	for _, handler := range api.GetRegisterAPIInterface() {
		s.webContainer.Add(handler.GetWebServiceRoute())
	}

	config := restfulSpec.Config{
		WebServices:                   s.webContainer.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}
	s.webContainer.Add(restfulSpec.NewOpenAPIService(config))
	return config
}

func (s *restServer) OPTIONSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if req.Request.Method != "OPTIONS" {
		chain.ProcessFilter(req, resp)
		return
	}
	resp.AddHeader(restful.HEADER_AccessControlAllowCredentials, "true")
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "1ch0 go-restful api doc",
			Description: "1ch0 go-restful api doc",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "1ch0 go-restful",
					Email: "github1ch0@163.com",
					URL:   "https://1ch0.github.io/",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache License 2.0",
					URL:  "https://github.com/1ch0/go-restful/blob/main/LICENSE",
				},
			},
			Version: "v1",
		},
	}
}

func (s *restServer) requestLog(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if req.HeaderParameter("Upgrade") == "websocket" && req.HeaderParameter("Connection") == "Upgrade" {
		chain.ProcessFilter(req, resp)
		return
	}
	start := time.Now()
	c := utils.NewResponseCapture(resp.ResponseWriter)
	resp.ResponseWriter = c
	chain.ProcessFilter(req, resp)
	takeTime := time.Since(start)
	log.Logger.With(
		"clientIP", pkgUtils.Sanitize(utils.ClientIP(req.Request)),
		"path", pkgUtils.Sanitize(req.Request.URL.Path),
		"method", req.Request.Method,
		"status", c.StatusCode(),
		"time", takeTime.String(),
		"reponseSize", len(c.Bytes()),
	).Infof("request log")
}

func (s *restServer) startHTTP(ctx context.Context) error {
	// Start HTTP apiserver
	log.Logger.Infof("HTTP APIs are being served on: %s, ctx: %s", s.cfg.BindAddr, ctx)
	server := &http.Server{Addr: s.cfg.BindAddr, Handler: s.webContainer}
	return server.ListenAndServe()
}
