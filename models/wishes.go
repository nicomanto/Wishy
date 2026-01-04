package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/signintech/gopdf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseWish struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Link       string             `json:"link" bson:"link"`
	Preference PreferenceType     `json:"preference" bson:"preference"`
	IsRecent   bool               `json:"recent" bson:"recent"`
	Ts         time.Time          `json:"ts" bson:"ts"`
}

type Wish struct {
	BaseWish `json:",inline" bson:",inline"`
	Category BaseCategory       `json:"cat" bson:"cat"`
	UserId   primitive.ObjectID `json:"uid" bson:"uid"`
	Active   bool               `json:"active" bson:"active"`
}

func (w *BaseWish) SetRecent(month int) {
	fmt.Print(time.Now().String())
	monthsAgo := time.Now().AddDate(0, -month, 0)
	w.IsRecent = w.Ts.After(monthsAgo)
	fmt.Print(w.IsRecent)
}

func (c BaseWish) DBCollectionName() string {
	return "wishes"
}

type PreferenceType int

const (
	Low PreferenceType = iota + 1
	Medium
	High
)

// Convert preference level to string
func (p PreferenceType) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return ""
	}
}

type WishByCategory struct {
	Cat    string     `json:"cat" bson:"_id"`
	Wishes []BaseWish `json:"wishes" bson:"wishes"`
}

type UserWishes struct {
	Wishes     []WishByCategory `json:"wishes" bson:"wishes"`
	Username   string           `json:"username" bson:"username"`
	LastUpdate string           `json:"last_update" bson:"last_update"`
}

// GenerateWishlistPDF generates the PDF based on user wishes
func (user UserWishes) GenerateWishlistPDF() ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	// Load Font
	err := pdf.AddTTFFont("Go-Italic", "./templates/fonts/Go-Italic.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}
	// Define Page Constraints
	const MaxY = 800 // Bottom margin for A4
	// Title (Left-Aligned)
	pdf.SetFont("Go-Italic", "", 18)
	pdf.SetY(30)
	pdf.Cell(nil, fmt.Sprintf("%s's Wishlist", user.Username))
	// Extra Space
	pdf.Br(30)
	// Last Update (Left-Aligned, Smaller Font)
	pdf.SetFont("Go-Italic", "", 10)
	pdf.Cell(nil, fmt.Sprintf("Last update: %s", user.LastUpdate))
	// Extra Space
	pdf.Br(40)
	// Wishlist Items
	pdf.SetFont("Go-Italic", "", 12)
	for _, category := range user.Wishes {
		// Check if there's space before adding a new category
		if pdf.GetY() > MaxY {
			pdf.AddPage()
			pdf.SetY(50) // Reset position on new page
		}
		// Print category title
		pdf.SetFont("Go-Italic", "B", 12)
		pdf.Cell(nil, fmt.Sprintf("- %s:", category.Cat))
		pdf.Br(20)
		// Reset font for items
		pdf.SetFont("Go-Italic", "", 12)
		for _, wish := range category.Wishes {
			// Check if there's space before adding a new wish
			if pdf.GetY() > MaxY {
				pdf.AddPage()
				pdf.SetY(50) // Reset position on new page
			}
			// Print wish item
			wishText := fmt.Sprintf("  • %s: %s", wish.Name, wish.Preference.String())
			pdf.Cell(nil, wishText)
			// Add clickable link
			if wish.Link != "" {
				pdf.SetTextColor(0, 0, 255)
				linkX, linkY := pdf.GetX(), pdf.GetY()
				pdf.Cell(nil, " (Click here)")
				pdf.AddExternalLink(wish.Link, linkX, linkY, 60, 10)
				pdf.SetTextColor(0, 0, 0)
			}
			pdf.Br(15) // Space between wishes
		}
		pdf.Br(20) // Space between categories
	}
	// Convert PDF to []byte
	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %v", err)
	}
	return buf.Bytes(), nil
}
