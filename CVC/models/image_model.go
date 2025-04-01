package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// DockerImage represents a Docker image
type DockerImage struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Label      string             `bson:"label" json:"label"`
	Repository string             `bson:"repository" json:"repository"`
	Tag        string             `bson:"tag" json:"tag"`
	CreateddAt int64              `bson:"pulled_at" json:"created_at"`
}
