package main

type Offer struct {
	Id string `json: "id"`
	Title string `json: "title"`
	Content string `json: "content"`
	CreatedAt string `json: "createdAt"`
}

type Company struct {
	Id string `json: "id"`
	Name string `json: "name"`
	Description string `json: "description"`
}
