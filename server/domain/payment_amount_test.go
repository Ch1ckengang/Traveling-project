package domain

import "testing"

func TestPaymentAmountUsesVND(t *testing.T) {
	payment := NewPayment(1, 2500000)

	if err := payment.ValidateAmount(); err != nil {
		t.Fatalf("ValidateAmount() unexpected error = %v", err)
	}

	if got := payment.GetAmountInVND(); got != 2500000 {
		t.Fatalf("GetAmountInVND() = %v, want 2500000", got)
	}

	summary := payment.ToSummary()
	if summary.Amount != 2500000 {
		t.Fatalf("PaymentSummary.Amount = %v, want 2500000", summary.Amount)
	}
}

func TestPaymentAmountMinimumVND(t *testing.T) {
	payment := NewPayment(1, 4999)

	if err := payment.ValidateAmount(); err == nil {
		t.Fatal("ValidateAmount() expected error for amount below 5,000 VND")
	}
}
