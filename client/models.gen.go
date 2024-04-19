// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package client

type SearchAllResponse struct { // struct
	Media *SearchResponse `bson:"media" json:"media"`
	Tmdb  *SearchResponse `bson:"tmdb" json:"tmdb"`
	Tvdb  *SearchResponse `bson:"tvdb" json:"tvdb"`
}

type SearchResponse struct { // struct
	Results []*SearchResult `bson:"results" json:"results"`
	Error   string          `bson:"error" json:"error"`
}

type SearchResult struct { // struct
	ID          string `bson:"id" json:"id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
	Kind        string `bson:"kind" json:"kind"`
	Date        string `bson:"date" json:"date"`
	Source      string `bson:"source" json:"source"`
	Image       string `bson:"image" json:"image"`
	Completed   bool   `bson:"completed" json:"completed"`
}
