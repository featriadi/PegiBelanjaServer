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

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "What's Good!")
	})

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
	e.POST("/customer/register", controller.StoreCustomer)
	e.PUT("/customer/update/:id", controller.UpdateCustomer)
	e.DELETE("/customer/delete/:id", controller.DeleteCustomer)
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
	e.POST("/mitra/create", controller.StoreMitra)
	e.PUT("/mitra/update/:id", controller.UpdateMitra)
	e.DELETE("/mitra/delete/:id", controller.DeleteMitra)
	//End Mitra

	//Review
	e.GET("/review/get", controller.FetchAllReview)
	e.GET("/review/get/:product_id", controller.FetchReviewByProductId)
	e.POST("/review/create", controller.StoreReview)
	//End Review

	//PBCourier
	e.GET("/p-courier/get", controller.FetchAllPBCourier)
	e.POST("/p-courier/create", controller.StorePBCourier)
	e.PUT("/p-courier/update", controller.UpdatePBCourier)
	e.DELETE("/p-courier/delete/:id", controller.DeletePBCourier)
	//End PBCourier

	//Courier
	e.GET("/courier/get", controller.FetchAllCourier)
	e.POST("/courier/create", controller.StoreCourier)
	e.PUT("/courier/update/:id", controller.UpdateCourier)
	e.DELETE("/courier/delete/:id", controller.DeleteCourier)
	//End Courier

	//Order
	e.POST("/order/create", controller.CreateOrder)
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

	return e
}
