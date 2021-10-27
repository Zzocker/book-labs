package delivery

import (
	"context"

	pb "github.com/Zzocker/book-labs/protos/book"
	"github.com/Zzocker/book-labs/protos/common"
)

type bookService struct {
	pb.UnimplementedBookServiceServer
}

func NewBookService() pb.BookServiceServer {
	return &bookService{}
}

func (b *bookService) CreateBook(context.Context, *pb.CreateBookRequest) (*common.EmptyResponse, error) {
	return nil, nil
}
func (b *bookService) GetBook(context.Context, *pb.ISBNRequest) (*pb.Book, error) {
	return nil, nil
}
func (b *bookService) UpdateBook(context.Context, *pb.CreateBookRequest) (*pb.Book, error) {
	return nil, nil
}
func (b *bookService) DeleteBook(context.Context, *pb.ISBNRequest) (*common.EmptyResponse, error) {
	return nil, nil
}
func (b *bookService) GetBookPic(context.Context, *pb.ISBNRequest) (*pb.BookCover, error) {
	return nil, nil
}
