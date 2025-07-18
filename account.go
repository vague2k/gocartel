package gocartel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type accounts struct {
	Data []accountData `json:"data"`
	// Links struct{}      `json:"links"`
}

type accountData struct {
	ID            string               `json:"id"`
	Type          string               `json:"type"`
	Attributes    accountAttributes    `json:"attributes"`
	Links         accountLinks         `json:"links"`
	Relationships accountRelationships `json:"relationships"`
	Meta          struct{}             `json:"meta"`
}

type accountAttributes struct {
	Subdomain                    string    `json:"subdomain"`
	StoreName                    string    `json:"store_name"`
	Description                  string    `json:"description"`
	FirstName                    string    `json:"first_name"`
	LastName                     string    `json:"last_name"`
	ContactEmail                 string    `json:"contact_email"`
	URL                          string    `json:"url"`
	Website                      string    `json:"website"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
	UnderMaintenance             bool      `json:"under_maintenance"`
	InventoryEnabled             bool      `json:"inventory_enabled"`
	ArtistsEnabled               bool      `json:"artists_enabled"`
	StripePublishableKey         any       `json:"stripe_publishable_key"`
	Launched                     bool      `json:"launched"`
	HasActiveAdvancedTaxSettings bool      `json:"has_active_advanced_tax_settings"`
	TimeZone                     string    `json:"time_zone"`
}

type accountLinks struct {
	Self string `json:"self"`
}

type accountRelationships struct {
	Currency       relationshipData `json:"currency"`
	Country        relationshipData `json:"country"`
	Plan           relationshipData `json:"plan"`
	AccountImage   relationshipData `json:"account_image"`
	AccountFavicon struct {
		Data any `json:"data"`
	} `json:"account_favicon"`
	Orders     relatedLinks     `json:"orders"`
	Categories selfRelatedLinks `json:"categories"`
	Products   selfRelatedLinks `json:"products"`
	Discounts  selfRelatedLinks `json:"discounts"`
}

type relationshipData struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}

type relatedLinks struct {
	Links struct {
		Related string `json:"related"`
	} `json:"links"`
}

type selfRelatedLinks struct {
	Links struct {
		Self    string `json:"self"`
		Related string `json:"related"`
	} `json:"links"`
}

func (c bigCartelClient) AccountWithContext(ctx context.Context) (*accountData, error) {
	resp, err := c.get(ctx, "/accounts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &accounts{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("no account data found")
	}

	return &result.Data[0], nil
}

func (c bigCartelClient) Account() (*accountData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeoutDuration)
	defer cancel()
	return c.AccountWithContext(ctx)
}
