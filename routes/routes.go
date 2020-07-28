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
	//End Product

	//Mitra
	e.GET("/mitra/get", controller.FetchAllMitraData)
	e.POST("/mitra/create", controller.StoreMitra)
	e.PUT("/mitra/update/:id", controller.UpdateMitra)
	e.DELETE("/mitra/delete/:id", controller.DeleteMitra)
	//End Mitra

	//User
	e.POST("/user/add", controller.StoreUser)
	e.GET("/user/get", controller.FetchAllUserData)
	// e.POST("/")
	//ENd User

	//Test
	e.POST("/test", controller.StoreProduct)
	//End

	//Auth
	e.GET("/generate-hash/:password", controller.GenerateHashPassword)
	e.GET("/is-loggedin", controller.Restricted, middlewares.IsAuthtenticated)
	e.GET("/is-mitra", controller.Restricted, middlewares.IsAuthtenticated, controller.IsMitra)
	e.POST("/login", controller.CheckLogin)
	e.POST("/refresh-token", controller.Token)
	e.POST("/mitra/register", controller.StoreMitraWithUser)
	//End Auth

	return e
}
