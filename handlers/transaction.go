package handlers

import (
	"expense-tracker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type TransactionHandler struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewTransactionHandler(ctx context.Context, collection *mongo.Collection, usrCollection *mongo.Collection) *TransactionHandler {
	return &TransactionHandler{
		collection:     collection,
		userCollection: usrCollection,
		ctx:            ctx,
	}
}

func (handler *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	var user models.User
	email := c.MustGet("email")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userErr := handler.userCollection.FindOne(handler.ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if userErr != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"user": userErr.Error()})
		return
	}

	transaction.ID = primitive.NewObjectID()
	transaction.Owner = user.ID
	const shortForm = "2006-01-02"
	dt, _ := time.Parse(shortForm, transaction.Date)
	transaction.InvDt = primitive.NewDateTimeFromTime(dt)
	createdTransaction, createErr := handler.collection.InsertOne(handler.ctx, transaction)

	_, updateErr := handler.userCollection.UpdateOne(handler.ctx, bson.M{
		"email": email,
	}, bson.D{{"$push", bson.D{
		{"transactions", createdTransaction.InsertedID},
	}}})

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	if createErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating new category"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (handler *TransactionHandler) ListTransaction(c *gin.Context) {
	var user models.User
	email := c.MustGet("email")
	userErr := handler.userCollection.FindOne(handler.ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if userErr != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"user": userErr.Error()})
		return
	}

	// TODO: remove owner from response
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"owner", bson.D{
				{"$eq", user.ID},
			},
			},
		}}},
		{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "category"},
			{"foreignField", "_id"},
			{"as", "cat"},
		}}},
	}

	cur, err := handler.collection.Aggregate(handler.ctx, pipeline)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)
	transactions := make([]models.Transaction, 0)

	for cur.Next(handler.ctx) {
		var transaction models.Transaction
		cur.Decode(&transaction)
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}

func (handler *TransactionHandler) DeleteTransaction(c *gin.Context) {
	email := c.MustGet("email")
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.DeleteOne(handler.ctx, bson.D{{"_id", objectId}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	_, updateErr := handler.userCollection.UpdateOne(handler.ctx, bson.M{
		"email": email,
	}, bson.D{{"$pull", bson.D{
		{"transactions", objectId},
	}}})

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successfully removed"})
}

func (handler *TransactionHandler) UpdateTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const shortForm = "2006-01-02"
	dt, _ := time.Parse(shortForm, transaction.Date)
	transaction.InvDt = primitive.NewDateTimeFromTime(dt)

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"amount", transaction.Amount},
		{"category", transaction.Category},
		{"date", transaction.Date},
		{"InvDt", transaction.InvDt},
	}}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction was successfully updated"})
}

func (handler *TransactionHandler) GetTransactionsByCategory(c *gin.Context) {
	var user models.User
	email := c.MustGet("email")

	userErr := handler.userCollection.FindOne(handler.ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if userErr != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"user": userErr.Error()})
		return
	}

	matchStage := bson.D{{"$match", bson.D{{"owner", bson.D{{"$eq", user.ID}}}}}}
	group := bson.D{{
		"$group", bson.D{
			{
				"_id", "$category",
			},
			{"total", bson.D{{"$sum", "$amount"}}},
		},
	}}
	pipeline := mongo.Pipeline{
		matchStage,
		group,
		{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "_id"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
	}

	cur, err := handler.collection.Aggregate(handler.ctx, pipeline)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)
	transactions := make([]models.TransactionCategory, 0)

	for cur.Next(handler.ctx) {
		var transaction models.TransactionCategory
		cur.Decode(&transaction)
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
