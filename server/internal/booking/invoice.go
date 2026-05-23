package booking

import (
	"bytes"
	"fmt"
	"time"
	"travel-backend/domain"

	"github.com/jung-kurt/gofpdf"
)

// GenerateInvoicePDF creates a PDF invoice for a booking and returns the raw bytes
func GenerateInvoicePDF(booking *domain.Booking) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(190, 15, "TRAVELING INVOICE", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Booking Info
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(95, 8, "Booking Information", "B", 0, "L", false, 0, "")
	pdf.CellFormat(95, 8, "Customer Information", "B", 1, "R", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 11)
	
	// Row 1
	pdf.CellFormat(95, 6, fmt.Sprintf("Booking Code: %s", booking.BookingCode), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Name: %s", booking.FullName), "", 1, "R", false, 0, "")
	
	// Row 2
	pdf.CellFormat(95, 6, fmt.Sprintf("Date: %s", booking.CreatedAt.Format("02/01/2006")), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Email: %s", booking.Email), "", 1, "R", false, 0, "")

	// Row 3
	pdf.CellFormat(95, 6, fmt.Sprintf("Status: %s", booking.Status), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Phone: %s", booking.Phone), "", 1, "R", false, 0, "")

	pdf.Ln(15)

	// Tour Details
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "Tour Details", "B", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(100, 8, "Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Quantity", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, "Total Amount", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	var tourName string
	if booking.Tour.Name != "" {
		tourName = booking.Tour.Name
	} else {
		tourName = fmt.Sprintf("Tour ID: %d", booking.TourID)
	}

	desc := fmt.Sprintf("%s (Date: %s)", tourName, booking.TravelDate)
	if len(desc) > 40 {
		desc = desc[:37] + "..."
	}

	pdf.CellFormat(100, 10, desc, "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%d", booking.Quantity), "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("%s VND", formatNumber(booking.TotalAmount)), "1", 1, "R", false, 0, "")

	// Discount / Coupon
	if booking.DiscountAmount > 0 {
		pdf.CellFormat(130, 10, fmt.Sprintf("Discount (Coupon: %s)", booking.CouponCode), "1", 0, "R", false, 0, "")
		pdf.CellFormat(60, 10, fmt.Sprintf("-%s VND", formatNumber(booking.DiscountAmount)), "1", 1, "R", false, 0, "")
	}

	// Total
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 10, "Grand Total:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("%s VND", formatNumber(booking.TotalAmount)), "1", 1, "R", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(190, 6, "Thank you for choosing Traveling!", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 6, fmt.Sprintf("Generated at %s", time.Now().Format("02/01/2006 15:04:05")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatNumber(n int64) string {
	in := fmt.Sprintf("%d", n)
	numOfDigits := len(in)
	if numOfDigits <= 3 {
		return in
	}
	out := make([]byte, 0, numOfDigits+(numOfDigits-1)/3)
	for i, c := range in {
		if i > 0 && (numOfDigits-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}
