package escpos

import (
	"fmt"
	"strings"
)

/*
════════════════════════════════════
 ESC/POS COMMANDS
════════════════════════════════════
*/

const (
	cmdInit = "\x1b@"

	// Alignment
	cmdAlignLeft   = "\x1ba\x00"
	cmdAlignCenter = "\x1ba\x01"
	// cmdAlignRight  = "\x1ba\x02"

	// Emphasis
	cmdBoldOn  = "\x1bE\x01"
	cmdBoldOff = "\x1bE\x00"

	// Reverse (Black Background)
	cmdReverseOn  = "\x1dB\x01"
	cmdReverseOff = "\x1dB\x00"

	// Size
	// cmdDoubleHW = "\x1d!\x30"
	cmdNormal = "\x1d!\x00"

	// Cut
	cmdCut = "\x1dV\x42\x00"

	// Underline (Used for Thick Separator Lines)
	cmdUnderlineThick = "\x1b-\x02"
	cmdUnderlineOff   = "\x1b-\x00"

	// ✅ LEFT MARGIN CORRECTION
	// Changed from \x30 (48) to \x40 (64 dots).
	// This adds approx 16 more dots of spacing on the left.
	// It pushes the content to the center, removing the extra space on the right.
	cmdMarginLeft = "\x1dL\x22\x00"

	// ✅ PRINTER WIDTH
	// Kept at 42 to prevent "ED" wrapping.
	printerWidth = 47
)

/*
════════════════════════════════════
 DATA MODELS
════════════════════════════════════
*/

type ReceiptItem struct {
	Name   string  `json:"name"`
	Qty    int     `json:"qty"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type ReceiptData struct {
	Outlet struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Address    string `json:"address"`
		GST        string `json:"gst"`
		LogoBase64 string `json:"logoBase64"`
	} `json:"outlet"`

	InvoiceNo string `json:"invoiceNo"`
	TRN       string `json:"trn"`

	Customer  string `json:"customer"`
	OrderType string `json:"orderType"`
	Table     string `json:"table"`
	DateTime  string `json:"dateTime"`

	Items []ReceiptItem `json:"items"`

	Subtotal       float64 `json:"subtotal"`
	VATPercent     float64 `json:"vatPercent"`
	VATAmount      float64 `json:"vatAmount"`
	DeliveryCharge float64 `json:"deliveryCharge"`
	Coupon         float64 `json:"coupon"`
	Total          float64 `json:"total"`

	Payments map[string]float64 `json:"payments"`
	Notes    string             `json:"notes"`
}

/*
════════════════════════════════════
 RECEIPT BUILDER
════════════════════════════════════
*/

func setLeftMargin(dots int) string {
	low := byte(dots % 256)
	high := byte(dots / 256)
	return string([]byte{0x1D, 0x4C, low, high})
}

func BuildReceipt(r ReceiptData) []byte {
	var b strings.Builder

	// ---------- Helpers ----------

	// ✅ FIXED: Bold Solid Separator Line
	// Using thick underline + spaces is cleaner than dashes ("-")
	// separator := func() {
	// 	b.WriteString(cmdAlignLeft) // Reset align to ensure line starts at margin
	// 	b.WriteString(cmdUnderlineThick)
	// 	b.WriteString(strings.Repeat("-", printerWidth))
	// 	b.WriteString(cmdUnderlineOff + "\n")
	// }
	separator := func() {
		b.WriteString(cmdAlignLeft)
		b.WriteString(strings.Repeat("-", printerWidth))
		b.WriteString("\n")
	}

	// Print Key-Value (Left ...... Right)
	printKV := func(key, val string) {
		b.WriteString(cmdAlignLeft)
		space := printerWidth - len(key) - len(val)
		if space < 1 {
			space = 1
		}
		b.WriteString(key + strings.Repeat(" ", space) + val + "\n")
	}

	// ---------- 1. INIT & MARGIN ----------
	// b.WriteString(cmdInit)
	// b.WriteString(cmdMarginLeft) // Sets the corrected global margin

	// ---------- 2. HEADER ----------
	// b.WriteString(cmdAlignCenter)

	// Top greeting
	// b.WriteString("Thank you for dining with us!\n\n")

	// [LOGO PLACEHOLDER]
	// b.Write(logoBytes)
	if r.Outlet.LogoBase64 != "" {
		logoBytes := ProcessBase64Logo(r.Outlet.LogoBase64)
		if logoBytes != nil {
			// b.WriteString("\x1dL\x00\x00")

			// b.WriteString("\x1dL\x8a\x00")
			b.WriteString(cmdAlignCenter)
			// 2. Print the image
			b.Write(logoBytes)
			b.WriteString("\n")

			// 3. Restore the text margin settings immediately after
			b.WriteString(cmdMarginLeft)

			// 4. Re-assert Center alignment for the text following the logo
			// b.WriteString(cmdAlignCenter)
		}
	}
	b.WriteString("\n")
	b.WriteString(cmdInit)

	// b.WriteString(cmdMarginLeft) 
	// Sets the corrected global margin
	b.WriteString(cmdAlignCenter)

	// Outlet Name: Normal Size, Bold
	b.WriteString(cmdNormal + cmdBoldOn)
	b.WriteString(r.Outlet.Name + "\n")
	b.WriteString(cmdBoldOff)

	// Contact Info
	if r.Outlet.Email != "" {
		b.WriteString(r.Outlet.Email + "\n")
	}
	b.WriteString("Phone: " + r.Outlet.Phone + "\n")
	if r.Outlet.GST != "" {
		b.WriteString("GST: " + r.Outlet.GST + "\n")
	}
	b.WriteString("\n")

	separator()

	// ---------- 3. INVOICE META ----------
	b.WriteString(cmdAlignCenter)
	b.WriteString("TAX INVOICE\n")
	b.WriteString("INV # " + r.InvoiceNo + "\n")
	if r.TRN != "" {
		b.WriteString("TRN: " + r.TRN + "\n")
	}

	separator()

	// ---------- 4. CUSTOMER INFO ----------
	b.WriteString(cmdAlignCenter)
	b.WriteString("Customer Info\n")
	b.WriteString(r.Customer + "\n")

	separator()

	// ---------- 5. ORDER DETAILS ----------
	// printKV relies on printerWidth. Since we increased the margin,
	// this block will now visually shift to the center.
	printKV("Order Type", r.OrderType)

	if r.Table != "" {
		printKV("Table", r.Table)
	}
	printKV("Date", r.DateTime)

	separator()

	// ---------- 6. ITEMS HEADER ----------
	b.WriteString(cmdBoldOn + cmdAlignLeft)
	b.WriteString(fmt.Sprintf("%-3s %-25s %12s\n", "Qty", "Item", "Amount"))
	b.WriteString(cmdBoldOff)

	// ---------- 7. ITEMS LOOP ----------
	for _, it := range r.Items {
		amt := fmt.Sprintf("%.2f AED", it.Amount)
		qty := fmt.Sprintf("%d", it.Qty)

		// Calculate Max Name Width
		maxName := printerWidth - 4 - len(amt) - 1

		wrappedName := wrapText(it.Name, maxName)
		lines := strings.Split(wrappedName, "\n")

		// Line 1: Qty | Name Part 1 | Amount
		b.WriteString(fmt.Sprintf("%-3s %s", qty, lines[0]))

		// Align Amount to Right
		padding := printerWidth - 3 - len(lines[0]) - len(amt)
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString(amt + "\n")

		// Print wrapped name lines
		for i := 1; i < len(lines); i++ {
			b.WriteString("    " + lines[i] + "\n")
		}

		// Unit Price (Indented)
		b.WriteString(fmt.Sprintf("    (%.2f AED)\n", it.Price))
	}

	separator()

	// ---------- 8. TOTALS ----------
	printKV("Sub Total", fmt.Sprintf("%.2f AED", r.Subtotal))
	printKV(fmt.Sprintf("VAT (%.0f%%)", r.VATPercent), fmt.Sprintf("%.2f AED", r.VATAmount))

	if r.DeliveryCharge > 0 {
		printKV("Delivery", fmt.Sprintf("%.2f AED", r.DeliveryCharge))
	}
	if r.Coupon > 0 {
		printKV("Coupon", fmt.Sprintf("%.2f AED", r.Coupon))
	}

	separator()

	// TOTAL (Bold)
	b.WriteString(cmdBoldOn)
	printKV("TOTAL", fmt.Sprintf("%.2f AED", r.Total))
	b.WriteString(cmdBoldOff)

	separator()

	// ---------- 9. PAYMENTS ----------
	for k, v := range r.Payments {
		printKV(strings.ToUpper(k), fmt.Sprintf("%.2f AED", v))
	}

	separator()

	// ---------- 10. FOOTER ----------
	b.WriteString(cmdAlignCenter)
	if r.Notes != "" {
		b.WriteString(r.Notes + "\n")
	}
	b.WriteString("Have a great meal from\n")
	b.WriteString(r.Outlet.Name + "\n\n")

	// PAID TAG
	b.WriteString(cmdReverseOn + "  PAID  " + cmdReverseOff + "\n\n")

	// Cut
	b.WriteString(cmdCut)

	return []byte(b.String())
}

/*
════════════════════════════════════
 UTIL
════════════════════════════════════
*/

func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}
	var out strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for _, w := range words {
		if lineLen+len(w) > width {
			out.WriteString("\n")
			lineLen = 0
		}
		if lineLen > 0 {
			out.WriteString(" ")
			lineLen++
		}
		out.WriteString(w)
		lineLen += len(w)
	}
	return out.String()
}



