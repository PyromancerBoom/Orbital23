package main

import (
	calculator "rpc_services/calculatorService/kitex_gen/calculator"
	"context"
	"errors"
	"strconv"
)

// CalculatorServiceImpl implements the last service interface defined in the IDL.
type CalculatorServiceImpl struct{}

// AddNums implements the CalculatorServiceImpl interface.
func (s *CalculatorServiceImpl) AddNums(ctx context.Context, request *calculator.CalculationRequest) (resp *calculator.CalculationResponse, err error) {
	result := request.A + request.B
	response := &calculator.CalculationResponse{
		Result: strconv.Itoa(result),
	}
	return response, nil
}

// SubNums implements the CalculatorServiceImpl interface.
func (s *CalculatorServiceImpl) SubNums(ctx context.Context, request *calculator.CalculationRequest) (resp *calculator.CalculationResponse, err error) {
	result := request.A - request.B
	response := &calculator.CalculationResponse{
		Result: strconv.Itoa(result),
	}
	return response, nil
}

// DivNums implements the CalculatorServiceImpl interface.
func (s *CalculatorServiceImpl) DivNums(ctx context.Context, request *calculator.CalculationRequest) (resp *calculator.CalculationResponse, err error) {
	if request.B == 0 {
		return nil, errors.New("division by zero error")
	}
	result := request.A / request.B
	response := &calculator.CalculationResponse{
		Result: strconv.Itoa(result),
	}
	return response, nil
}

// MultiplyNums implements the CalculatorServiceImpl interface.
func (s *CalculatorServiceImpl) MultiplyNums(ctx context.Context, request *calculator.CalculationRequest) (resp *calculator.CalculationResponse, err error) {
	result := request.A * request.B
	response := &calculator.CalculationResponse{
		Result: strconv.Itoa(result),
	}
	return response, nil
}