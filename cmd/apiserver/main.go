package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1ch0/go-restful/pkg/apiserver"
	"github.com/1ch0/go-restful/pkg/apiserver/config"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/log"
)

func main() {
	s := &Server{}
	flag.StringVar(&s.serverConfig.BindAddr, "bind-addr", "0.0.0.0:8080", "The bind address used to serve the http APIs.")
	flag.Parse()

	//if len(os.Args) > 2 && os.Args[1] == "build-swagger" {
	//	func() {
	//		swagger, err := s.buildSwagger()
	//		if err != nil {
	//			log.Logger.Fatal(err.Error())
	//		}
	//		outData, err := json.MarshalIndent(swagger, "", "\t")
	//		if err != nil {
	//			log.Logger.Fatal(err.Error())
	//		}
	//		swaggerFile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	//		if err != nil {
	//			log.Logger.Fatal(err.Error())
	//		}
	//		defer func() {
	//			if err := swaggerFile.Close(); err != nil {
	//				log.Logger.Error("close swagger file failur %s", err.Error())
	//			}
	//		}()
	//		_, err = swaggerFile.Write(outData)
	//		if err != nil {
	//			log.Logger.Fatal(err.Error())
	//		}
	//		fmt.Println("build swagger config file success")
	//	}()
	//	return
	//}

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := s.run(ctx, errChan); err != nil {
			errChan <- fmt.Errorf("failed to run apiserver: %w", err)
		}
	}()
	var term = make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.Logger.Infof("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		log.Logger.Errorf("Received an error: %s, exiting gracefully...", err.Error())
	}
	log.Logger.Infof("See you next time!")
}

type Server struct {
	serverConfig config.Config
}

func (s *Server) run(ctx context.Context, errChan chan error) error {
	log.Logger.Infof("1ch0/go-restful informatin: version: %s", "1.0.0")

	server := apiserver.New(s.serverConfig)

	return server.Run(ctx, errChan)
}

//func (s *Server) buildSwagger() (*spec.Swagger, error) {
//	server := apiserver.New(s.serverConfig)
//	config, err := server.BuildRestfulConfig()
//	if err != nil {
//		return nil, err
//	}
//	return restfulSpec.BuildSwagger(*config), nil
//}
