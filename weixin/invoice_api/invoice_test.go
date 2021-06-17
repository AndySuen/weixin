package invoice_api

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func TestInvoice(t *testing.T) {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})

	invoiceApi := NewOfficialAccountApi(officialAccount)
	setbizattrObj := &SetbizattrObj{
		Phone:   "4006280808",
		TimeOut: 7200,
	}

	// {
	// 	result, err := invoiceApi.GetAuthData(&AuthDataObj{
	// 		OrderID: "1623894036798956962",
	// 		SPappID: "d3g1OTg2NGE5ZTU3ODIyOWVhX_oiY7-5OuzNHme3fHyMQQWjstgWqHfPcktQ40c-H73D",
	// 	})
	// 	require.Equal(t, nil, err)
	// 	fmt.Println(result)
	// }
	// {
	// 	err := invoiceApi.RejectInsert(&RejectInsertObj{
	// 		OrderID: "1623894036798956962",
	// 		SPappID: "d3g1OTg2NGE5ZTU3ODIyOWVhX_oiY7-5OuzNHme3fHyMQQWjstgWqHfPcktQ40c-H73D",
	// 		Reason:  "some reason",
	// 	})
	// 	require.Equal(t, nil, err)
	// }
	///////////////////////////// debug

	{
		file, err := os.Open(test.InvoicePdf)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		mediaID, err := invoiceApi.PlatformSetpdf("fapiao.pdf", fi.Size(), file)
		require.Equal(t, nil, err)
		fmt.Printf("media id %s\n", mediaID)
	}

	{
		cardID, err := invoiceApi.PlatformCreateCard(&CreateCardObj{
			Payee: test.InvoicePayee,
			Type:  test.InvoiceType,
			BaseInfo: &CreateCardBaseInfo{
				Title:                "汽车通用服务发票",
				CustomUrlName:        test.InvoiceCustomUrlName,
				CustomURL:            test.InvoiceCustomURL,
				CustomUrlSubTitle:    test.InvoiceCustomUrlSubTitle,
				PromotionUrlName:     "查看其他",
				PromotionURL:         "https://www.baidu.com",
				PromotionUrlSubTitle: "详情",
				LogoUrl:              "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png",
			},
		})
		require.Equal(t, nil, err)
		fmt.Printf("card id %s\n", cardID)
	}

	{
		err := invoiceApi.SetContact(setbizattrObj)
		require.Equal(t, nil, err)
	}
	{
		result, err := invoiceApi.GetContact()
		require.Equal(t, nil, err)
		require.Equal(t, result.Phone, setbizattrObj.Phone)
		require.Equal(t, result.TimeOut, setbizattrObj.TimeOut)
	}
	spappID := ""
	{
		result, err := invoiceApi.SetUrl()
		require.Equal(t, nil, err)
		fmt.Println(result)

		u, err := url.Parse(result)
		require.Equal(t, nil, err)
		m, err := url.ParseQuery(u.RawQuery)
		require.Equal(t, nil, err)
		pappid, ok := m["s_pappid"]
		require.Equal(t, true, ok)
		require.NotEmpty(t, pappid)
		spappID = pappid[0]
		fmt.Printf("s_pappid : %s\n", spappID)
	}

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	{
		ticket, _, err := invoiceApi.OfficialAccount.GetWxCardApiTicket()
		require.Equal(t, nil, err)

		result, err := invoiceApi.GetAuthUrl(&AuthUrlObj{
			SPappID:   spappID,
			Money:     1,
			Source:    "web",
			OrderID:   orderID,
			Timestamp: time.Now().Unix(),
			Type:      1,
			Ticket:    ticket,
		})
		require.Equal(t, nil, err)

		fmt.Println(result.AuthURL, result.AppID)
	}

	{
		result, err := invoiceApi.GetAuthData(&AuthDataObj{
			OrderID: orderID,
			SPappID: spappID,
		})
		require.Equal(t, nil, err)
		fmt.Println(result)
	}
}
