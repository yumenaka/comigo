package data_api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/comigo_remote"
)

// remoteComigoClientByKey 根据前端 remote_store 参数找到配置中的远端 Comigo，并创建请求客户端。
func remoteComigoClientByKey(remoteStoreKey string) (*comigo_remote.Client, string, error) {
	storeURL, err := comigo_remote.MatchStoreByKey(config.GetCfg().StoreUrls, remoteStoreKey)
	if err != nil {
		return nil, "", err
	}
	client, err := comigo_remote.NewClient(storeURL, config.GetCfg().TimeoutLimitForScan)
	if err != nil {
		return nil, "", err
	}
	return client, storeURL, nil
}

// remoteComigoBookFromRequest 判断本次请求是否要代理到远端 Comigo，并返回本地缓存的远端定位信息。
func remoteComigoBookFromRequest(c echo.Context, localBookID string) (*model.Book, *comigo_remote.Client, string, bool, error) {
	remoteStoreKey := c.QueryParam(comigo_remote.RemoteStoreQuery)
	if remoteStoreKey == "" {
		return nil, nil, "", false, nil
	}
	client, storeURL, err := remoteComigoClientByKey(remoteStoreKey)
	if err != nil {
		return nil, nil, "", true, err
	}
	book, err := model.IStore.GetBook(localBookID)
	if err != nil {
		return nil, nil, "", true, err
	}
	if book.RemoteBookID == "" {
		return nil, nil, "", true, fmt.Errorf("book is not a remote Comigo book: %s", localBookID)
	}
	return book, client, storeURL, true, nil
}

// remoteComigoQuery 复制当前查询参数并替换为远端原始 BookID，避免把本地 BookID 传给远端服务。
func remoteComigoQuery(c echo.Context, remoteBookID string) url.Values {
	query := url.Values{}
	for key, values := range c.Request().URL.Query() {
		if key == comigo_remote.RemoteStoreQuery {
			continue
		}
		for _, value := range values {
			query.Add(key, value)
		}
	}
	query.Set("id", remoteBookID)
	return query
}

// writeRemoteComigoError 统一把远端 Comigo 失败映射为网关错误，便于前端区分本地校验失败。
func writeRemoteComigoError(c echo.Context, err error) error {
	return apiresp.Error(c, http.StatusBadGateway, "remote_comigo_failed", err.Error(), nil)
}
