package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Container represents a container's metadata and security status
type Container struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContainerID string             `bson:"container_id" json:"container_id"`
	Name        string             `bson:"name" json:"name"`
	BaseImage   string             `bson:"image" json:"image"`
	Status      string             `bson:"status" json:"status"`
	//CreatedAt   time.Time          `bson:"created_at" json:"created_at"`

	UserID string `bson:"userId" json:"userId"`
}

type FalcoAlert struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        string             `json:"userId" bson:"userId"`
	Rule          string             `json:"rule" bson:"rule"`
	Priority      string             `json:"priority" bson:"priority"`
	Source        string             `json:"source" bson:"source"`
	Count         int                `json:"count" bson:"count"`
	ContainerID   string             `json:"container.id" bson:"containerId"`
	ContainerName string             `json:"container.name" bson:"containerName"`
	ProcessName   string             `json:"proc.name" bson:"processName"`
}
