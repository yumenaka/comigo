package routers

//// 4、setWebpServer TODO：新的webp模式：https://docs.webp.sh/usage/remote-backend/
//func setWebpServer(engine *gin.Engine) {
//	//webp反向代理
//	if common.Config.EnableWebpServer {
//		webpError := common.StartWebPServer(common.CachePath+"/webp_config.json", common.ReadingBook.ExtractPath, common.CachePath+"/webp", common.Config.Port+1)
//		if webpError != nil {
//			fmt.Println(locale.GetString("webp_server_error"), webpError.Error())
//			//engine.Static("/cache", common.CachePath)
//
//		} else {
//			fmt.Println(locale.GetString("webp_server_start"))
//			engine.Use(reverse_proxy.ReverseProxyHandle("/cache", reverse_proxy.ReverseProxyOptions{
//				TargetHost:  "http://localhost",
//				TargetPort:  strconv.Itoa(common.Config.Port + 1),
//				RewritePath: "/cache",
//			}))
//		}
//	} else {
//		if common.ReadingBook.IsDir {
//			engine.Static("/cache/"+common.ReadingBook.BookID, common.ReadingBook.GetFilePath())
//		} else {
//			engine.Static("/cache", common.CachePath)
//		}
//	}
//}
