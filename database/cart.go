package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kaykobadhossain/e-commerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this product from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userColletion *mongo.Collection, productID primitive.ObjectID, userID string) error {

	searchfromdb, err := prodCollection.Find(ctx, primitive.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	if err = searchfromdb.All(ctx, &productCart); err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	update := bson.D{primitive.E{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{primitive.E{Key: "$each", Value: productCart}}}}}}

	_, err = userColletion.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userColletion *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	//update := bson.M{{Key: "$pull", Value: bson.M{"usercart": bson.M{{Key:"_id", Value: productID}}}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = userColletion.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil
}
func BuyItemFromCart(ctx context.Context, userColletion *mongo.Collection, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{primitive.E{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{primitive.E{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, primitive.E{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentresults, err := userColletion.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	ctx.Done()

	if err != nil {
		panic(err)
	}

	var getusercart []bson.M

	if err = currentresults.All(ctx, &getusercart); err != nil {
		log.Println(err)
		return ErrCantGetItem
	}

	var total_price int32

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	update := bson.D{primitive.E{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}

	_, err = userColletion.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userColletion.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"each": getcartitems.UserCart}}}

	_, err = userColletion.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser, 0)

	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	//update3 := bson.D{"$set": bson.D{"usercart": usercart_empty}}

	_, err = userColletion.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InsatantBuyer(ctx context.Context, prodCollection, userColletion *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)

	if err != nil {
		log.Println(err)
	}
	orders_detail.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{primitive.E{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	_, err = userColletion.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userColletion.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	return nil
}
