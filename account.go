package gocartel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"
)

type account struct {
	ID               string          // The unique ID of the account.
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
	var JSON json.RawMessage
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
	var JSON json.RawMessage
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

	parsed := gjson.ParseBytes(b)
	if !parsed.Get(dataQuery+".id").Exists() || !parsed.Get(attributes+".url").Exists() {
		return nil, fmt.Errorf("no account data found")
	}

	acc := &account{
		ID:               parsed.Get(dataQuery + ".id").String(),
		Subdomain:        parsed.Get(attributes + ".subdomain").String(),
		StoreName:        parsed.Get(attributes + ".store_name").String(),
		Description:      parsed.Get(attributes + ".description").String(),
		ContactEmail:     parsed.Get(attributes + ".contact_email").String(),
		FirstName:        parsed.Get(attributes + ".first_name").String(),
		LastName:         parsed.Get(attributes + ".last_name").String(),
		URL:              parsed.Get(attributes + ".url").String(),
		Website:          parsed.Get(attributes + ".website").String(),
		CreatedAt:        parsed.Get(attributes + ".created_at").String(),
		UpdatedAt:        parsed.Get(attributes + ".updated_at").String(),
		UnderMaintenance: parsed.Get(attributes + ".under_maintenance").Bool(),
		InventoryEnabled: parsed.Get(attributes + ".inventory_enabled").Bool(),
		ArtistsEnabled:   parsed.Get(attributes + ".artists_enabled").Bool(),
		Timezone:         parsed.Get(attributes + ".time_zone").String(),
		Currency: accountCurrency{
			ID:     parsed.Get(`included.#(type="currencies").id`).String(),
			Name:   parsed.Get(`included.#(type="currencies").attributes.name`).String(),
			Sign:   parsed.Get(`included.#(type="currencies").attributes.sign`).String(),
			Locale: parsed.Get(`included.#(type="currencies").attributes.locale`).String(),
		},
		Country: accountCountry{
			ID:   parsed.Get(`included.#(type="countries").id`).String(),
			Name: parsed.Get(`included.#(type="countries").attributes.name`).String(),
		},
		Plan: accountPlan{
			ID:                  parsed.Get(`included.#(type="plans").id`).String(),
			Name:                parsed.Get(`included.#(type="plans").attributes.name`).String(),
			MaxProducts:         parsed.Get(`included.#(type="plans").attributes.max_products`).Int(),
			MaxImagesPerProduct: parsed.Get(`included.#(type="plans").attributes.max_images_per_product`).Int(),
			MonthlyRate:         parsed.Get(`included.#(type="plans").attributes.monthly_rate`).Float(),
		},
		Image: accountImage{
			ID:  parsed.Get(`included.#(type="account_images").id`).String(),
			URL: parsed.Get(`included.#(type="account_images").attributes.url`).String(),
		},
		Links: accountLinks{
			Self:       parsed.Get(dataQuery + ".links.self").String(),
			Orders:     parsed.Get(dataQuery + ".relationships.orders.links.related").String(),
			Categories: parsed.Get(dataQuery + ".relationships.categories.links.related").String(),
			Products:   parsed.Get(dataQuery + ".relationships.products.links.related").String(),
		},
	}

	return acc, nil
}
