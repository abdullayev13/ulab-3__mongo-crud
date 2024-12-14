package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"my_app/internal/models"
	"my_app/internal/pkg/mongodb"
	"net/http"
)

func CreateProduct(c *gin.Context) {
	data := new(models.Product)
	err := c.ShouldBind(data)
	if err != nil {
		fail(c, err)
		return
	}

	if data.Name == "" {
		failMsg(c, "product name is required")
	}

	coll := mongodb.GetColl("products")

	data.Id = primitive.ObjectID{}
	res, err := coll.InsertOne(c, data)
	if err != nil {
		fail(c, err)
		return
	}

	data.Id = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, data)

}

func GetOneProduct(c *gin.Context) {
	id, err := objIfFromParam(c, "id")
	if err != nil {
		fail(c, err)
		return
	}

	coll := mongodb.GetColl("products")

	res := coll.FindOne(c, Obj{"_id": id})
	err = res.Err()
	if err != nil {
		fail(c, err)
		return
	}

	model := new(models.Product)
	err = res.Decode(model)
	if err != nil {
		fail(c, err)
		return
	}

	c.JSON(http.StatusOK, model)

}

func GetProducts(c *gin.Context) {
	data := new(productFilter)
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
		if data.Name != nil {
			match["name"] = data.Name
		}
		if data.Price != nil {
			match["price"] = data.Price
		}
		if data.CategoryName != nil {
			match["category_name"] = data.CategoryName
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

	coll := mongodb.GetColl("products")

	cur, err := coll.Aggregate(c, pipe)
	if err != nil {
		fail(c, err)
		return
	}

	list := make([]models.Product, 0)
	for cur.Next(c) {
		var val models.Product
		err = cur.Decode(&val)
		if err != nil {
			fail(c, err)
			return
		}

		list = append(list, val)

	}

	c.JSON(http.StatusOK, list)

}

func UpdateProduct(c *gin.Context) {
	id, err := objIfFromParam(c, "id")
	if err != nil {
		fail(c, err)
		return
	}

	data := new(models.Product)
	err = c.ShouldBind(data)
	if err != nil {
		fail(c, err)
		return
	}

	if data.Name == "" {
		failMsg(c, "product name is required")
	}

	data.Id = id

	coll := mongodb.GetColl("products")

	res, err := coll.UpdateOne(c, Obj{"_id": id}, data)
	if err != nil {
		fail(c, err)
		return
	}

	if res.MatchedCount == 0 {
		failMsg(c, "product not found")
		return
	}

	c.JSON(http.StatusOK, data)

}

func DeleteProduct(c *gin.Context) {
	id, err := objIfFromParam(c, "id")
	if err != nil {
		fail(c, err)
		return
	}

	coll := mongodb.GetColl("products")

	res, err := coll.DeleteOne(c, Obj{"_id": id})
	if err != nil {
		fail(c, err)
		return
	}

	if res.DeletedCount == 0 {
		failMsg(c, "product not found")
		return
	}

	c.JSON(http.StatusOK, "done")

}

type productFilter struct {
	Name         *string  `json:"name" form:"name"`
	Price        *float64 `json:"price" form:"price"`
	CategoryName *string  `json:"category_name" form:"category_name"`
	Search       *string  `json:"search" form:"search"`
	SortField    *string  `json:"sorting_field" form:"sorting_field"`
	Desc         bool     `json:"desc" form:"desc"`
	Limit        int      `json:"limit" form:"limit"`
	Offset       int      `json:"offset" form:"offset"`
}
