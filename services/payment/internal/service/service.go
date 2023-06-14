package service

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) ProcessPayment(payment *model.Payment) error {

	payment.Status = "success"

	return nil
}
