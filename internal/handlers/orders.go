package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"my_app/internal/models"
	"my_app/internal/pkg/mongodb"
	"net/http"
)

func CreateOrder(c *gin.Context) {
	userId := getUserInfo(c).UserId

	data := new(models.Order)
	err := c.ShouldBind(data)
	if err != nil {
		fail(c, err)
		return
	}

	data.OrderedBy = userId

	if len(data.Items) == 0 {
		failMsg(c, "order item(s) is required")
	}

	coll := mongodb.GetColl("orders")

	data.Id = primitive.ObjectID{}
	res, err := coll.InsertOne(c, data)
	if err != nil {
		fail(c, err)
		return
	}

	data.Id = res.InsertedID.(primitive.ObjectID)
	{
		itemColl := mongodb.GetColl("order_items")

		docs := []any{}
		for i := range data.Items {
			data.Items[i].OrderId = data.Id
			docs = append(docs, &data.Items[i])
		}

		res, err := itemColl.InsertMany(c, docs)
		if err != nil {
			fail(c, err)
			return
		}

		for i, id := range res.InsertedIDs {
			data.Items[i].OrderId = id.(primitive.ObjectID)
		}
	}

	c.JSON(http.StatusOK, data)

}

func GetOneOrder(c *gin.Context) {
	id, err := objIfFromParam(c, "id")
	if err != nil {
		fail(c, err)
		return
	}

	coll := mongodb.GetColl("orders")

	cur, err := coll.Aggregate(c, []any{
		Obj{
			"$match": Obj{"_id": id},
		},
		Obj{
			"$lookup": Obj{
				"from":         "order_items",
				"localField":   "_id",
				"foreignField": "order_id",
				"as":           "items",
			},
		},
	})

	if err != nil {
		fail(c, err)
		return
	}

	type order struct {
		models.Order `bson:",inline"`
		OrderItems   []models.OrderItem `bson:"items"`
	}
	model := new(order)
	err = cur.Decode(model)
	if err != nil {
		fail(c, err)
		return
	}

	model.Items = model.OrderItems

	c.JSON(http.StatusOK, model.Order)

}

func GetOrders(c *gin.Context) {
	data := new(orderFilter)
	err := c.ShouldBindQuery(data)
	if err != nil {
		fail(c, err)
		return
	}
	if data.Limit == 0 {
		data.Limit = 10
	}

	pipe := []any{}
	match := Obj{}
	{
		if data.OrderedBy != nil {
			match["ordered_by"] = data.OrderedBy
		}
		if data.Comment != nil {
			match["comment"] = data.Comment
		}
	}
	if len(match) > 0 {
		pipe = append(pipe, Obj{
			"$match": match,
		})
	}

	if data.Search != nil {
		pipe = append(pipe, Obj{
			"$text": Obj{
				"$search": *data.Search,
			},
		})
	}

	if data.SortField != nil {
		val := 1
		if data.Desc {
			val = -1
		}

		pipe = append(pipe, Obj{
			"$sort": Obj{
				*data.SortField: val,
			},
		})

	}

	if data.Offset > 0 {
		pipe = append(pipe, Obj{
			"$skip": data.Offset,
		})
	}

	pipe = append(pipe, Obj{
		"$limit": data.Limit,
	})

	coll := mongodb.GetColl("orders")

	cur, err := coll.Aggregate(c, pipe)
	if err != nil {
		fail(c, err)
		return
	}

	list := make([]models.Order, 0)
	for cur.Next(c) {
		var val models.Order
		err = cur.Decode(&val)
		if err != nil {
			fail(c, err)
			return
		}

		list = append(list, val)

	}

	c.JSON(http.StatusOK, list)

}

func DeleteOrder(c *gin.Context) {
	id, err := objIfFromParam(c, "id")
	if err != nil {
		fail(c, err)
		return
	}

	coll := mongodb.GetColl("orders")

	res, err := coll.DeleteOne(c, Obj{"_id": id})
	if err != nil {
		fail(c, err)
		return
	}

	if res.DeletedCount == 0 {
		failMsg(c, "order not found")
		return
	}

	res, err = mongodb.GetColl("order_items").DeleteMany(c, Obj{"order_id": id})
	if err != nil {
		fail(c, err)
		return
	}

	c.JSON(http.StatusOK, "done")

}

type orderFilter struct {
	OrderedBy *primitive.ObjectID `json:"ordered_by" form:"ordered_by"`
	Comment   *string             `json:"comment" form:"comment"`

	Search    *string `json:"search" form:"search"`
	SortField *string `json:"sorting_field" form:"sorting_field"`
	Desc      bool    `json:"desc" form:"desc"`
	Limit     int     `json:"limit" form:"limit"`
	Offset    int     `json:"offset" form:"offset"`
}
