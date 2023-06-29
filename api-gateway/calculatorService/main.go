package main

import (
	calculator "api-gateway/calculatorService/kitex_gen/calculator/calculatorservice"
	"log"
)

func main() {
	svr := calculator.NewServer(new(CalculatorServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
