package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/SamuelsSantos/product-discount-service/discountcalculator/config"
	"github.com/SamuelsSantos/product-discount-service/discountcalculator/domain"
	"github.com/SamuelsSantos/product-discount-service/discountcalculator/domain/pb"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.NewConfig()
	log.Printf("User Service: %s", cfg.UserConfig.ToURL())
	log.Printf("Product Service: %s", cfg.ProductConfig.ToURL())
	log.Printf("Black Friday Date: %s", cfg.ToStringISO())

	userService := domain.NewUserClient(cfg.UserConfig.ToURL())
	productService := domain.NewProductClient(cfg.ProductConfig.ToURL())

	grpcServer := grpc.NewServer()
	calculator := domain.NewCalculatorServer(userService, productService, cfg)

	port := flag.Int("port", 50001, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	pb.RegisterDiscountCalculatorServiceServer(grpcServer, calculator)
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
