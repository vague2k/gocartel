package gocartel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/tidwall/gjson"
)

type Category struct {
	ID        string
	Name      string
	Permalink string
	Position  string
}

func (a account) CategoriesWithContext(ctx context.Context) ([]Category, error) {
	resp, err := a.parentClient.get(ctx, "/accounts/"+a.ID+"/categories")
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

	var categories []Category
	result := gjson.GetBytes(JSON, "data")
	result.ForEach(func(key, value gjson.Result) bool {
		c := Category{
			ID:        value.Get("id").String(),
			Name:      value.Get("attributes.name").String(),
			Permalink: value.Get("attributes.permalink").String(),
			Position:  value.Get("attributes.position").String(),
		}
		categories = append(categories, c)
		return true
	})

	return categories, nil
}

func (a account) CategoriesByIDWithContext(ctx context.Context, id string) (*Category, error) {
	resp, err := a.parentClient.get(ctx, "/accounts/"+a.ID+"/categories/"+id)
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

	result := gjson.GetBytes(JSON, "data")
	c := &Category{
		ID:        result.Get("id").String(),
		Name:      result.Get("attributes.name").String(),
		Permalink: result.Get("attributes.permalink").String(),
		Position:  result.Get("attributes.position").String(),
	}
	return c, nil
}

func (a account) CategoriesCreateWithContext(ctx context.Context, name string) (*Category, error) {
	// Because of BigCartel's API expectations, there's no way around nesting regardless
	// of wether it's a struct or maps.
	//
	// I personally didn't feel like defining a struct specifically for this payload
	// was neccessary but perhaps this oculd change. For now this works just fine
	payload := map[string]any{
		"data": map[string]any{
			"type": "categories", // "categories" is the only accepted type for this payload currently.
			"attributes": map[string]any{
				"name": name,
			},
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := a.parentClient.post(ctx, "/accounts/"+a.ID+"/categories/", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("the request to create the category '%s' was unsuccessful. the status code is '%d'. expected '201' ", name, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var JSON json.RawMessage
	err = json.Unmarshal(body, &JSON)
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(JSON, "data")
	c := &Category{
		ID:        result.Get("id").String(),
		Name:      result.Get("attributes.name").String(),
		Permalink: result.Get("attributes.permalink").String(),
		Position:  result.Get("attributes.position").String(),
	}

	return c, nil
}

func (a account) Categories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.parentClient.Client.Timeout)
	defer cancel()
	return a.CategoriesWithContext(ctx)
}

func (a account) CategoriesByID(id string) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.parentClient.Client.Timeout)
	defer cancel()
	return a.CategoriesByIDWithContext(ctx, id)
}

func (a account) CreateCategories(name string) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.parentClient.Client.Timeout)
	defer cancel()
	return a.CategoriesCreateWithContext(ctx, name)
}
