package invoice_api

import (
	"encoding/json"
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

func newInvoiceApi() *InvoiceApi {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})

	return NewOfficialAccountApi(officialAccount)
}

func TestInvoiceUploadPdf(t *testing.T) {
	api := newInvoiceApi()

	file, err := os.Open(test.InvoicePdf)
	require.Empty(t, err)
	defer file.Close()

	fi, err := file.Stat()
	require.Empty(t, err)

	mediaID, err := api.PlatformSetpdf("fapiao.pdf", fi.Size(), file)
	require.Equal(t, nil, err)
	fmt.Printf("media id %s\n", mediaID)
}

func TestSetContact(t *testing.T) {
	api := newInvoiceApi()

	setbizattrObj := &SetbizattrObj{
		Phone:   test.InvoicePhone,
		TimeOut: 7200,
	}

	err := api.SetContact(setbizattrObj)
	require.Equal(t, nil, err)

	result, err := api.GetContact()
	require.Equal(t, nil, err)
	require.Equal(t, result.Phone, setbizattrObj.Phone)
	require.Equal(t, result.TimeOut, setbizattrObj.TimeOut)
}

func TestPlatformCreateCard(t *testing.T) {
	api := newInvoiceApi()

	cardID, err := api.PlatformCreateCard(&CreateCardObj{
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

func TestInvoiceInsert(t *testing.T) {
	api := newInvoiceApi()

	param := &InvoiceInsertObj{
		OrderID: "1623934539819863381",
		CardID:  "pWVnF6HIt6_NNMm4NvlNhrUWBWZc",
		Appid:   test.OfficialAccountAppid,
		CardExt: &InvoiceInsertCardExt{
			NonceStr: fmt.Sprintf("%d", time.Now().UnixNano()),
			UserCard: struct {
				InvoiceUserData *InvoiceInsertCardExtUser `json:"invoice_user_data"`
			}{
				InvoiceUserData: &InvoiceInsertCardExtUser{
					Fee:           209300,
					Title:         test.InvoicePayee,
					BillingTime:   int(time.Now().Unix()),
					BillingNO:     "47516482",
					BillingCode:   "012002000511",
					CheckCode:     "58626678892273709512",
					FeeWithoutTax: 197453,
					Tax:           11847,
					SPdfMediaID:   "2164102706656625645",
					Info: []InvoiceInsertCardExtItem{
						{
							Name:  "洗车",
							Price: 1,
							Num:   1,
							Unit:  "次",
						},
					},
				},
			},
		},
	}

	b, _ := json.Marshal(param)
	fmt.Println(string(b))

	result, err := api.Insert(param)
	require.Equal(t, nil, err)
	fmt.Printf("code : %s, openid: %s, unionid: %s\n", result.Code, result.OpenID, result.UnionID)

}

func TestInvoice(t *testing.T) {
	api := newInvoiceApi()

	// {
	// 	result, err := api.GetAuthData(&AuthDataObj{
	// 		OrderID: "1623930786748654309",
	// 		SPappID: "d3g1OTg2NGE5ZTU3ODIyOWVhX_oiY7-5OuzNHme3fHyMQQWjstgWqHfPcktQ40c-H73D",
	// 	})
	// 	require.Equal(t, nil, err)
	// 	fmt.Println(result.InvoiceStatus)
	// }
	// {
	// 	err := api.RejectInsert(&RejectInsertObj{
	// 		OrderID: "1623894036798956962",
	// 		SPappID: "d3g1OTg2NGE5ZTU3ODIyOWVhX_oiY7-5OuzNHme3fHyMQQWjstgWqHfPcktQ40c-H73D",
	// 		Reason:  "some reason",
	// 	})
	// 	require.Equal(t, nil, err)
	// }
	///////////////////////////// debug

	spappID := ""
	{
		result, err := api.SetUrl()
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
	fmt.Printf("order id %s\n", orderID)
	{
		ticket, _, err := api.OfficialAccount.GetWxCardApiTicket()
		require.Equal(t, nil, err)

		result, err := api.GetAuthUrl(&AuthUrlObj{
			SPappID:   spappID,
			Money:     209300,
			Source:    "web",
			OrderID:   orderID,
			Timestamp: time.Now().Unix(),
			Type:      1,
			Ticket:    ticket,
		})
		require.Equal(t, nil, err)
		fmt.Printf("url : %s , appid %s\n", result.AuthURL, result.AppID)
	}

	{
		result, err := api.GetAuthData(&AuthDataObj{
			OrderID: orderID,
			SPappID: spappID,
		})
		require.Equal(t, nil, err)
		fmt.Println(result.InvoiceStatus)
	}
}
