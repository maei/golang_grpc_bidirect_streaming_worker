package service

var GreetService greetServiceInterface = &greetService{}

type greetServiceInterface interface {
}

type greetService struct{}

func (*greetService) Greeting() {

}
