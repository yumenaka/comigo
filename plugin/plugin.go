package plugin

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"gopkg.in/ini.v1"

	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
)

func StartFrpC(configPath string) error {
	//借助ini库，保存一个ini文件
	cfg := ini.Empty()
	//配置文件类似：
	//[common]
	//server_addr = frp.example.net
	//server_port = 7000
	//token = Nscffaass
	//[comi]
	//type = tcp
	//local_ip = 127.0.0.1
	//local_port = 1234
	//remote_port = 23456
	_, err := cfg.NewSection("common")
	_, err = cfg.Section("common").NewKey("server_addr", common.Config.FrpConfig.ServerAddr)
	_, err = cfg.Section("common").NewKey("server_port", strconv.Itoa(common.Config.FrpConfig.ServerPort))
	_, err = cfg.Section("common").NewKey("token", common.Config.FrpConfig.Token)
	FrpConfigName := common.ReadingBook.Name + "(" + "comi " + common.Version + " " + time.Now().Format("2006-01-02 15:04:05") + ")"
	_, err = cfg.NewSection(FrpConfigName)
	_, err = cfg.Section(FrpConfigName).NewKey("type", common.Config.FrpConfig.FrpType)
	_, err = cfg.Section(FrpConfigName).NewKey("local_ip", "127.0.0.1")
	_, err = cfg.Section(FrpConfigName).NewKey("local_port", strconv.Itoa(common.Config.Port))
	_, err = cfg.Section(FrpConfigName).NewKey("remote_port", strconv.Itoa(common.Config.FrpConfig.RemotePort))
	//保存文件
	err = cfg.SaveToIndent(configPath+"/frpc.ini", "\t")
	if err != nil {
		fmt.Println(locale.GetString("frpc_ini_error"))
		return err
	} else {
		fmt.Println(locale.GetString("frpc_setting_save_completed"), configPath, cfg)
	}
	//实际执行
	var cmd *exec.Cmd
	cmd = exec.Command(common.Config.FrpConfig.FrpcCommand, "-c", configPath+"/frpc.ini")
	fmt.Println(cmd)
	if err = cmd.Start(); err != nil {
		return err
	}
	return err
}

//func StartWebPServer(webpConfigFile string, imgPath string, exhaustPath string, port int) error {
//	//Config.WebpCommand = wepBinaryPath
//	Config.WebpConfig.ImgPath = imgPath
//	Config.WebpConfig.ExhaustPath = exhaustPath
//	Config.WebpConfig.PORT = strconv.Itoa(port)
//	//Config.WebpConfig.QUALITY = quality
//	if Config.WebpConfig.WebpCommand == "" || Config.WebpConfig.ImgPath == "" || Config.WebpConfig.ExhaustPath == "" {
//		return errors.New(locale.GetString("webp_setting_error"))
//	}
//	jsonObject, err := os.OpenFile(webpConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
//	if err != nil {
//		return err
//	}
//	defer jsonObject.Close()
//	content, err := json.Marshal(Config.WebpConfig)
//	if err != nil {
//		return err
//	}
//	if _, err := jsonObject.Write(content); err == nil {
//		fmt.Println(locale.GetString("webp_setting_save_completed"), webpConfigFile, content)
//	}
//	//err = webpCMD(webpConfigFile, Config.WebpCommand)
//	var cmd *exec.Cmd
//	cmd = exec.Command(Config.WebpConfig.WebpCommand, "--config", webpConfigFile)
//	fmt.Println(cmd)
//	if err = cmd.Start(); err != nil {
//		return err
//	}
//	return err
//}
