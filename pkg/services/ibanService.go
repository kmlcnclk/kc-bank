package services

import (
	"math/rand"
	"time"
)

type IIbanService interface {
	GenerateIBAN(countryCode string, bankCodeLen, accountLen int) string
}

type IbanService struct {
	rng *rand.Rand
}

func NewIbanService() *IbanService {
	src := rand.NewSource(time.Now().UnixNano())

	return &IbanService{
		rng: rand.New(src),
	}
}

func (s *IbanService) GenerateIBAN(countryCode string, bankCodeLen, accountLen int) string {

	checkDigits := s.randomString(2, "0123456789")
	bankCode := s.randomString(bankCodeLen, "0123456789")
	accountNumber := s.randomString(accountLen, "0123456789")

	return countryCode + checkDigits + bankCode + accountNumber
}

func (s *IbanService) randomString(length int, charset string) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[s.rng.Intn(len(charset))]
	}
	return string(result)
}
