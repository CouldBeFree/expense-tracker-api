package handlers

import (
	"expense-tracker-api/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type CategoryHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCategoryHandler(ctx context.Context, collection *mongo.Collection) *CategoryHandler {
	return &CategoryHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *CategoryHandler) ListCategory(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)
	categories := make([]models.Category, 0)

	for cur.Next(handler.ctx) {
		var category models.Category
		cur.Decode(&category)
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

func (handler *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caregoryErr := handler.collection.FindOne(handler.ctx, bson.M{
		"name": category.Name,
	}).Decode(&category)

	if caregoryErr != mongo.ErrNoDocuments && caregoryErr == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Category": "Category alredy exists"})
		return
	}

	category.ID = primitive.NewObjectID()
	createdCategory, err := handler.collection.InsertOne(handler.ctx, category)

	fmt.Println(createdCategory.InsertedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating new category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (handler *CategoryHandler) GetCategory(c *gin.Context) {
	var category models.Category

	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	err := handler.collection.FindOne(handler.ctx, bson.D{{"_id", objectId}}).Decode(&category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (handler *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	res, err := handler.collection.DeleteOne(handler.ctx, bson.D{{"_id", objectId}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	c.JSON(http.StatusOK, gin.H{"message": "Category successfully removed"})
}

func (handler *CategoryHandler) UpdateCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id")

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", category.Name},
		{"type", category.Type},
	}}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category was successfully updated"})
}
