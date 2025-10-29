package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	pb "grpc/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var jwtToken string

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	bookClient := pb.NewBookServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)

	r := gin.Default()

	r.SetTrustedProxies([]string{"127.0.0.1"})
	// auth

	r.POST("/auth/sign-up", func(ctx *gin.Context) {
		var user pb.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := userClient.SignUp(ctx, &user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"id": res.Id})
	})

	r.GET("/auth/sign-in", func(ctx *gin.Context) {
		var user pb.SignInRequest
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := userClient.SignIn(ctx, &user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(res)
		jwtToken = res.Token
		fmt.Println(jwtToken)
		ctx.JSON(http.StatusCreated, gin.H{"token": res.Token})
	})

	// books
	r.GET("/books", func(ctx *gin.Context) {
		mdCtx := withAuthMetadata(context.Background())
		res, err := bookClient.GetBooks(mdCtx, &pb.Empty{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"books": res.Books})
	})

	r.GET("/books/:id", func(ctx *gin.Context) {
		mdCtx := withAuthMetadata(context.Background())
		idParam := ctx.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		res, err := bookClient.GetBook(mdCtx, &pb.BookId{Id: uint32(id)})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"book": res})
	})

	r.POST("/books", func(ctx *gin.Context) {
		mdCtx := withAuthMetadata(context.Background())
		var book pb.Book
		if err := ctx.ShouldBindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := bookClient.CreateBook(mdCtx, &book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"id": res.Id})
	})

	r.PUT("/books/:id", func(ctx *gin.Context) {
		mdCtx := withAuthMetadata(context.Background())
		var book pb.Book
		idParam := ctx.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		if err := ctx.ShouldBindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		book.Id = uint32(id)
		res, err := bookClient.UpdateBook(mdCtx, &book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"book": res})
	})

	r.DELETE("/books/:id", func(ctx *gin.Context) {
		mdCtx := withAuthMetadata(context.Background())
		idParam := ctx.Param("id")
		id, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		_, err = bookClient.DeleteBook(mdCtx, &pb.BookId{Id: uint32(id)})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "book deleted"})
	})

	r.Run(":5000")
}

func withAuthMetadata(ctx context.Context) context.Context {
	if jwtToken == "" {
		return ctx
	}
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + jwtToken,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
