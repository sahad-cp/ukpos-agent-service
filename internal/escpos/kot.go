package escpos

import (
	"fmt"
	"strings"
)


// -------- DATA MODELS --------

type KOTItem struct {
	Name      string   `json:"name"`
	Qty       int      `json:"qty"`
	Size      string   `json:"size"`
	Modifiers []string `json:"modifiers"`
	Notes     string   `json:"notes"`
}

type KOTData struct {
	Restaurant struct {
		Name string `json:"name"`
	} `json:"restaurant"`

	Order struct {
		OrderId      string `json:"orderId"`
		OrderType    string `json:"orderType"`
		Table        string `json:"table"`
		DateTime     string `json:"dateTime"`
		CustomerName string `json:"customerName"`
		Phone        string `json:"phone"`
	} `json:"order"`

	KOT struct {
		KotNo string `json:"kotNo"`
	} `json:"kot"`

	Items []KOTItem `json:"items"`
	Notes     string   `json:"notes"`
}

// -------- BUILDER --------

func BuildKOT(k KOTData) []byte {
	var b strings.Builder

	line := strings.Repeat("-", printerWidth)

	b.WriteString(cmdInit)
	b.WriteString(cmdAlignLeft)

	// ---------- HEADER ----------
	b.WriteString(cmdBoldOn)
	b.WriteString(k.Restaurant.Name + "\n")
	b.WriteString(cmdBoldOff)
	b.WriteString(line + "\n")

	b.WriteString("Order ID: " + k.Order.OrderId + "\n")
	b.WriteString("Order Type: " + k.Order.OrderType + "\n")
	b.WriteString("Table: " + k.Order.Table + "\n")
	b.WriteString("Date: " + k.Order.DateTime + "\n")
	b.WriteString("Customer Name: " + k.Order.CustomerName + "\n")
	b.WriteString("Phone: " + k.Order.Phone + "\n")

	b.WriteString(line + "\n")
b.WriteString("\n")
	// ---------- KOT META ----------
	b.WriteString(cmdBoldOn)
	b.WriteString("KOT : " + k.KOT.KotNo + "\n")
	b.WriteString(cmdBoldOff)
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("%-28s %10s\n", "Item", "Quantity"))
	b.WriteString(line + "\n")

	// ---------- ITEMS ----------
	for _, it := range k.Items {

		// Item name + qty
		b.WriteString(cmdBoldOn)
		b.WriteString(it.Name)
		b.WriteString(strings.Repeat(" ", printerWidth-len(it.Name)-len(fmt.Sprint(it.Qty))))
		b.WriteString(fmt.Sprint(it.Qty) + "\n")
		b.WriteString(cmdBoldOff)

		// Size
		if it.Size != "" {
			b.WriteString("Size: " + it.Size + "\n")
		}

		// Modifiers
		for _, m := range it.Modifiers {
			b.WriteString(" - " + m + "\n")
		}

		// Notes
		if it.Notes != "" {
			b.WriteString(it.Notes + "\n")
		}

		b.WriteString(line + "\n")
	}
	b.WriteString("note : " + k.Notes + "\n")
	b.WriteString("\n")
	b.WriteString(cmdAlignCenter)
	b.WriteString("KITCHEN COPY\n\n")


	b.WriteString(cmdCut)

	return []byte(b.String())
}