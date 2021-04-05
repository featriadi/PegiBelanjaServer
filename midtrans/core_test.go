package midtrans_test

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/cheekybits/is"
	midtrans "github.com/veritrans/go-midtrans"
)

var orderId1 string

func TestCoreCharge(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)
	orderId1 = "order-id-go-" + timestamp

	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-KwXI1e8Nm35VtlkXbasdE8DT"
	midclient.ClientKey = "SB-Mid-client-n7FPWIDQI9ukXcPg"
	midclient.APIEnvType = midtrans.Sandbox
	midclient.LogLevel = 3

	coreGateway := midtrans.CoreGateway{
		Client: midclient,
	}

	chargeReq := &midtrans.ChargeReq{
		PaymentType: midtrans.SourceGopay,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId1,
			GrossAmt: 200000,
		},
		Gopay: &midtrans.GopayDetail{
			EnableCallback: true,
			CallbackUrl:    "https://example.org",
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Some item",
			},
		},
	}

	log.Println("Charge:")
	chargeResp, err := coreGateway.Charge(chargeReq)
	if err != nil {
		log.Println("Fail w/ err:")
		log.Fatal(err)
		log.Println(err)
	} else {
		log.Println("Success w/ res:")
		log.Println(chargeResp)
		is.OK(chargeResp)
		is.OK(chargeResp.Actions)
	}
}

func TestCoreChargeWithMap(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)
	orderId1 = "order-id-go-" + timestamp

	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-KwXI1e8Nm35VtlkXbasdE8DT"
	midclient.ClientKey = "SB-Mid-client-n7FPWIDQI9ukXcPg"
	midclient.APIEnvType = midtrans.Sandbox
	midclient.LogLevel = 3

	coreGateway := midtrans.CoreGateway{
		Client: midclient,
	}

	chargeReq := &midtrans.ChargeReqWithMap{
		"payment_type": midtrans.SourceGopay,
		"transaction_details": midtrans.TransactionDetails{
			OrderID:  orderId1,
			GrossAmt: 200000,
		},
		"gopay": &midtrans.GopayDetail{
			EnableCallback: true,
			CallbackUrl:    "https://example.org",
		},
		"item_details": &[]midtrans.ItemDetail{
			midtrans.ItemDetail{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Map Interface",
			},
		},
	}

	log.Println("Charge:")
	chargeResp, err := coreGateway.ChargeWithMap(chargeReq)
	if err != nil {
		log.Println("Fail w/ err:")
		log.Fatal(err)
		log.Println(err)
	} else {
		log.Println("Success w/ res:")
		log.Println(chargeResp)
		is.OK(chargeResp["actions"])
	}
}

func TestCoreStatus(t *testing.T) {
	is := is.New(t)

	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-client-n7FPWIDQI9ukXcPg"
	midclient.ClientKey = "SB-Mid-server-KwXI1e8Nm35VtlkXbasdE8DT"
	midclient.APIEnvType = midtrans.Sandbox
	midclient.LogLevel = 3

	coreGateway := midtrans.CoreGateway{
		Client: midclient,
	}

	log.Println("Status:")
	statusResp, err := coreGateway.Status(orderId1)
	if err != nil {
		log.Println("Fail w/ err:")
		log.Fatal(err)
		log.Println(err)
	} else {
		log.Println("Success w/ res:")
		log.Println(statusResp)
		is.OK(statusResp)
		is.Equal("pending", statusResp.TransactionStatus)
	}
}
