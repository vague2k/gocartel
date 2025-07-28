package gocartel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"
)

type account struct {
	ID               int64           // The unique ID of the account.
	Subdomain        string          // The unique subdomain for the account on .bigcartel.com.
	StoreName        string          // The name of the shop.
	Description      string          // A brief description of the shop.
	ContactEmail     string          // The email address the shop can be contacted at.
	FirstName        string          // The shop owner's first name
	LastName         string          // The shop owner's last name
	URL              string          // The URL of the online store, either using a custom domain or the default .bigcartel.com one.
	Website          string          // The URL of a external website the account may have, separate from Big Cartel.
	CreatedAt        string          // When the account was created.
	UpdatedAt        string          // When the account was last updated.
	UnderMaintenance bool            // Whether or not the account is in maintenance mode.
	InventoryEnabled bool            // Whether or not the account has inventory tracking enabled.
	ArtistsEnabled   bool            // Whether or not the account has artist categories enabled.
	Timezone         string          // The timezone of the account.
	Currency         accountCurrency // The account's currency details.
	Country          accountCountry  // The account's country details.
	Plan             accountPlan     // The account's plan details.
	Image            accountImage    // The account image's details.
	Links            accountLinks
}

type accountCurrency struct {
	ID     string // The unique 3-letter code for the account’s currency.
	Name   string // The name of the account’s currency.
	Sign   string // The symbol used for the account’s currency.
	Locale string // The locale associated with the account’s currency.
}

type accountCountry struct {
	ID   string // The unique 2-letter code for the account’s country.
	Name string // The name of the account's country.
}

type accountPlan struct {
	ID                  string  // The unique identifier for the account’s plan. Possible values are gold, platinum, diamond, or titanium.
	Name                string  // The name of the account’s plan.
	MaxProducts         int64   // The maximum number of products the account can have on their plan.
	MaxImagesPerProduct int64   // The maximum number of images per product the account can have on their plan.
	MonthlyRate         float64 // The rate the account is billed each month.
}

type accountImage struct {
	ID  string // The unique identifier for the account’s image.
	URL string // The customizable URL of the account’s image.
}

type accountLinks struct {
	Self       string // .../v1/accounts/{id}
	Orders     string // .../v1/accounts/{id}/orders
	Categories string // .../v1/accounts/{id}/categories
	Products   string // .../v1/accounts/{id}/products
}

func (c BigCartelClient) AccountWithContext(ctx context.Context) (*account, error) {
	resp, err := c.get(ctx, "/accounts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var JSON []byte
	err = json.Unmarshal(body, &JSON)
	if err != nil {
		return nil, err
	}

	acc, err := getAccountValues(JSON, false)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (c BigCartelClient) AccountByIDWithContext(ctx context.Context, id string) (*account, error) {
	resp, err := c.get(ctx, "/accounts/"+id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var JSON []byte
	err = json.Unmarshal(body, &JSON)
	if err != nil {
		return nil, err
	}

	acc, err := getAccountValues(JSON, true)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (c BigCartelClient) Account() (*account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Client.Timeout)
	defer cancel()
	return c.AccountWithContext(ctx)
}

func (c BigCartelClient) AccountByID(id string) (*account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Client.Timeout)
	defer cancel()
	return c.AccountByIDWithContext(ctx, id)
}

// im not very interested in writing a bunch of structs to marshal json into,
// the gjson package works flawlessly for my use cases.
func getAccountValues(b []byte, withID bool) (*account, error) {
	var dataQuery string
	if withID {
		dataQuery = "data"
	} else {
		dataQuery = "data.0"
	}

	attributes := fmt.Sprintf("%s.attributes", dataQuery)

	if !gjson.GetBytes(b, dataQuery+".id").Exists() || !gjson.GetBytes(b, attributes+".url").Exists() {
		return nil, fmt.Errorf("no account data found")
	}

	acc := &account{
		ID:               gjson.GetBytes(b, dataQuery+".id").Int(),
		Subdomain:        gjson.GetBytes(b, attributes+".subdomain").String(),
		StoreName:        gjson.GetBytes(b, attributes+".store_name").String(),
		Description:      gjson.GetBytes(b, attributes+".description").String(),
		ContactEmail:     gjson.GetBytes(b, attributes+".contact_email").String(),
		FirstName:        gjson.GetBytes(b, attributes+".first_name").String(),
		LastName:         gjson.GetBytes(b, attributes+".last_name").String(),
		URL:              gjson.GetBytes(b, attributes+".url").String(),
		Website:          gjson.GetBytes(b, attributes+".website").String(),
		CreatedAt:        gjson.GetBytes(b, attributes+".created_at").String(),
		UpdatedAt:        gjson.GetBytes(b, attributes+".updated_at").String(),
		UnderMaintenance: gjson.GetBytes(b, attributes+".under_maintenance").Bool(),
		InventoryEnabled: gjson.GetBytes(b, attributes+".inventory_enabled").Bool(),
		ArtistsEnabled:   gjson.GetBytes(b, attributes+".artists_enabled").Bool(),
		Timezone:         gjson.GetBytes(b, attributes+".time_zone").String(),
		Currency: accountCurrency{
			ID:     gjson.GetBytes(b, `included.#(type="currencies").id`).String(),
			Name:   gjson.GetBytes(b, `included.#(type="currencies").attributes.name`).String(),
			Sign:   gjson.GetBytes(b, `included.#(type="currencies").attributes.sign`).String(),
			Locale: gjson.GetBytes(b, `included.#(type="currencies").attributes.locale`).String(),
		},
		Country: accountCountry{
			ID:   gjson.GetBytes(b, `included.#(type="countries").id`).String(),
			Name: gjson.GetBytes(b, `included.#(type="countries").attributes.name`).String(),
		},
		Plan: accountPlan{
			ID:                  gjson.GetBytes(b, `included.#(type="plans").id`).String(),
			Name:                gjson.GetBytes(b, `included.#(type="plans").attributes.name`).String(),
			MaxProducts:         gjson.GetBytes(b, `included.#(type="plans").attributes.max_products`).Int(),
			MaxImagesPerProduct: gjson.GetBytes(b, `included.#(type="plans").attributes.max_images_per_product`).Int(),
			MonthlyRate:         gjson.GetBytes(b, `included.#(type="plans").attributes.monthly_rate`).Float(),
		},
		Image: accountImage{
			ID:  gjson.GetBytes(b, `included.#(type="account_images").id`).String(),
			URL: gjson.GetBytes(b, `included.#(type="account_images").attributes.url`).String(),
		},
		Links: accountLinks{
			Self:       gjson.GetBytes(b, dataQuery+".links.self").String(),
			Orders:     gjson.GetBytes(b, dataQuery+".relationships.orders.links.related").String(),
			Categories: gjson.GetBytes(b, dataQuery+".relationships.categories.links.related").String(),
			Products:   gjson.GetBytes(b, dataQuery+".relationships.products.links.related").String(),
		},
	}

	return acc, nil
}
