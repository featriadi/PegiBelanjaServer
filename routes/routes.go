package routes

import (
	"net/http"
	"pb-dev-be/controller"
	"pb-dev-be/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.DefaultStaticConfig))

	// CORS restricted
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: middlewares.ValidateKey,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "What's Good!")
	})

	//Image
	e.GET("/image/helper/:type/get-base64/:file_name", controller.GetEncodedImage)
	//End Image

	//Category
	// e.GET("/category/get", controller.FetchAllCategoryData, middlewares.IsAuthtenticated)
	e.GET("/category/get", controller.FetchAllCategoryData)
	e.POST("/category/create", controller.StoreCategory)
	e.PUT("/category/update/:id", controller.UpdateCategory)
	e.DELETE("/category/delete", controller.DeleteCategory)
	// End Category

	//Customer
	e.GET("/customer/get", controller.FetchAllCustomerData)
	e.GET("/customer/:id", controller.ShowCustomerDataById)
	e.POST("/customer/create", controller.StoreCustomer)
	e.POST("/customer/register", controller.StoreCustomerWithUser)
	e.PUT("/customer/update/:id", controller.UpdateCustomer)
	e.DELETE("/customer/delete/:id", controller.DeleteCustomer)

	e.GET("/customer/validate/:email", controller.CheckCustomerByEmail)
	//End Customer

	//Bank
	e.GET("/bank/get", controller.FetchAllBankData)
	e.POST("/bank/create", controller.StoreBank)
	e.PUT("/bank/update/:id", controller.UpdateBank)
	e.DELETE("/bank/delete", controller.DeleteBank)
	//End Bank

	//Coupon
	e.GET("/coupon/get", controller.FetchAllCouponData)
	e.POST("/coupon/create", controller.StoreCouponData)
	e.PUT("/coupon/update/:id", controller.UpdateCouponData)
	e.DELETE("/coupon/delete/:id", controller.DeleteCoupon)
	//END Coupon

	//Variant
	e.GET("/variant/get", controller.FetchAllVariantData)
	e.POST("/variant/create", controller.StoreVariant)
	e.PUT("/variant/update/:id", controller.UpdateVariant)
	e.DELETE("/variant/delete/:id", controller.DeleteVariant)
	//End Variant

	//DiscountType
	e.GET("/discount-type/get", controller.FetchAllDiscountData)
	e.POST("/discount-type/create", controller.StoreDiscountData)
	e.PUT("/discount-type/update/:id", controller.UpdateDiscountData)
	e.DELETE("/discount-type/delete/:id", controller.DeleteDiscountData)
	//End Discount Type

	//Banner
	e.GET("/banner/get", controller.FetchAllBannerData)
	e.POST("/banner/create", controller.StoreBanner)
	e.DELETE("/banner/delete/:id", controller.DeleteBanner)
	e.PUT("/banner/update/:id", controller.UpdateBanner)
	//End Banner

	//Product
	e.GET("/product/get", controller.FetchAllProductData)
	e.GET("/product/get/:id", controller.ShowProductDataById)
	e.POST("/product/create", controller.StoreProduct)
	e.PUT("/product/update/:id", controller.UpdateProduct)
	e.DELETE("/product/delete/:id", controller.DeleteProduct)
	e.POST("/product/file/upload", controller.UploadFileProduct)
	e.GET("/product/count", controller.GetTotalProducts)
	//End Product

	//Mitra
	e.GET("/mitra/get", controller.FetchAllMitraData)
	e.GET("/mitra/get/:id", controller.ShowMitraDataById)
	e.POST("/mitra/create", controller.StoreMitra)
	e.PUT("/mitra/update/:id", controller.UpdateMitra)
	e.DELETE("/mitra/delete/:id", controller.DeleteMitra)
	//End Mitra

	//Review
	e.GET("/review/get", controller.FetchAllReview)
	e.GET("/review/get/:product_id", controller.FetchReviewByProductId)
	e.POST("/review/create", controller.StoreReview)
	//End Review

	//Email
	// e.POST("/email/send", controller.SendEmail)
	//End Email

	//Verification
	e.POST("/verification/create", controller.CreateVerification)
	//End Verification

	//Time Delivery
	e.GET("/timedelivery/get", controller.FetchAllTimeDeliveryData)
	// e.GET("/review/get/:product_id", controller.FetchReviewByProductId)
	// e.POST("/review/create", controller.StoreReview)
	//End Time Delivery

	//Cart
	e.GET("/cart/get", controller.GetCartByCustomerId)
	e.POST("/cart/create", controller.CreateCart)
	e.PUT("/cart/update", controller.UpdateCart)
	e.DELETE("/cart/delete/:id", controller.DeleteCart)
	//End Cart

	//PBCourier
	e.GET("/p-courier/get", controller.FetchAllPBCourier)
	e.POST("/p-courier/create", controller.StorePBCourier)
	e.PUT("/p-courier/update", controller.UpdatePBCourier)
	e.DELETE("/p-courier/delete/:id", controller.DeletePBCourier)
	//End PBCourier

	//Courier
	e.GET("/courier/get", controller.FetchAllCourier)
	e.POST("/courier/create", controller.StoreCourier)
	e.POST("/courier/third-pty/cost", controller.GetThirdPartyCourier)
	e.PUT("/courier/update/:id", controller.UpdateCourier)
	e.DELETE("/courier/delete/:id", controller.DeleteCourier)
	//End Courier

	//Order
	e.POST("/order/create", controller.CreateOrder)
	e.POST("/order/createcore", controller.CreateOrderCore)
	// e.POST("/order/post", controller.ChargeDirect)
	e.GET("/order/get/transaction/status/:order_id", controller.GetTransactionStatus)
	e.GET("/order/get", controller.GetOrderData)
	e.GET("/order/stats/get", controller.GetOrderStats)
	e.POST("/order/tracking/create", controller.CreateOrderTracking)
	e.GET("/order/tracking/get/:order_id", controller.GetOrderTracking)
	e.PUT("/order/update/waybill/:order_id", controller.UpdateWaybillOrder)
	//End Order

	//Stock
	e.POST("/stock-h/create", controller.CreateStockHInAndOut)
	//End Stock

	//User
	e.POST("/user/add", controller.StoreUser)
	e.GET("/user/get", controller.FetchAllUserData)
	// e.GET("/test", controller.Test)
	//ENd User

	//TP
	e.POST("/province/get-and-restore", controller.StoreProvince)
	e.POST("/city/get-and-restore", controller.StoreCity)
	e.POST("/subdistrict/get-and-restore", controller.StoreSubdistrict)
	e.GET("/province/get", controller.FetchAllProvince)
	e.GET("/city/get", controller.FetchAllCity)
	e.GET("/subdistrict/get", controller.FetchAllSubDistrict)
	//END TP

	//Auth
	e.GET("/generate-hash/:password", controller.GenerateHashPassword)
	e.GET("/is-loggedin", controller.Restricted, middlewares.IsAuthtenticated)
	e.GET("/is-mitra", controller.Restricted, middlewares.IsAuthtenticated, controller.IsMitra)
	e.POST("/login", controller.CheckLogin)
	e.POST("/refresh-token", controller.Token)
	e.POST("/mitra/register", controller.StoreMitraWithUser)
	e.GET("/srv-key/get", controller.GetAuthForMidtrans)
	//End Auth

	//Acurate Auth
	// e.POST("https://account.accurate.id/oauth", controller.StoreAccurate)
	e.GET("/accurateauth", controller.Accurate)

	return e
}
