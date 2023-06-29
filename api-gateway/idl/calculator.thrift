namespace Go calculator

struct CalculationRequest {
  1: i32 a;
  2: i32 b;
}

struct CalculationResponse {
  1: string Result;
}

service CalculatorService {
  CalculationResponse addNums(1: CalculationRequest request);
  CalculationResponse subNums(1: CalculationRequest request);
  CalculationResponse divNums(1: CalculationRequest request);
  CalculationResponse multiplyNums(1: CalculationRequest request);
}
