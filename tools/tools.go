package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/bbrks/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/mandykoh/autocrop"

	"github.com/yumenaka/comi/locale"
)

// GetImageDataBlurHash  获取图片的BlurHash
func GetImageDataBlurHash(loadedImage []byte, components int) string {
	// Generate the BlurHash for a given image
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return "error blurhash!"
	}
	str, err := blurhash.Encode(components, components, imageData)
	if err != nil {
		// Handle errors
		fmt.Println(err)
		return "error blurhash!"
	}
	fmt.Printf("Hash: %s\n", str)
	return str
}

// GetImageDataBlurHashImage 获取图片的BlurHash图
func GetImageDataBlurHashImage(loadedImage []byte, components int) []byte {
	// Generate the BlurHash for a given image
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
	}
	str, err := blurhash.Encode(components, components, imageData)
	if err != nil {
		fmt.Println(err)
	}
	// Generate an imageData for a given BlurHash
	// Punch specifies the contrasts and defaults to 1
	img, err := blurhash.Decode(str, imageData.Bounds().Dx(), imageData.Bounds().Dy(), 1)
	if err != nil {
		fmt.Println(err)
	}
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, img, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResizeByWidth 根据一个固定宽度缩放图片
func ImageResizeByWidth(loadedImage []byte, width int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	sourceWidth := image.Bounds().Dx()
	scalingRatio := float64(width) / float64(sourceWidth)
	height := int(float64(image.Bounds().Dy()) * scalingRatio)
	//生成缩略图
	image = imaging.Resize(image, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, image, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResizeByMaxWidth  设定一个图片宽度上限，大于这个宽度就缩放
func ImageResizeByMaxWidth(loadedImage []byte, maxWidth int) ([]byte, error) {
	buf := bytes.NewBuffer(loadedImage)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("imaging.Decode() Error")
	}
	sourceWidth := image.Bounds().Dx()
	if maxWidth > sourceWidth {
		return nil, errors.New("ImageResizeByMaxWidth Error maxWidth(" + strconv.Itoa(maxWidth) + ")> sourceWidth(" + strconv.Itoa(sourceWidth) + ")")
	}
	scalingRatio := float64(maxWidth) / float64(sourceWidth)
	height := int(float64(image.Bounds().Dy()) * scalingRatio)
	//生成缩略图
	image = imaging.Resize(image, maxWidth, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, image, imaging.JPEG)
	if err != nil {
		return nil, errors.New("imaging.Encode() Error")
	}
	return buf2.Bytes(), nil
}

// ImageResizeByMaxHeight  设定一个图片高度上限，大于这个高度就缩放
func ImageResizeByMaxHeight(loadedImage []byte, maxHeight int) ([]byte, error) {
	buf := bytes.NewBuffer(loadedImage)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("imaging.Decode() Error")
	}
	sourceHeight := image.Bounds().Dy()
	if maxHeight > sourceHeight {
		return nil, errors.New("ImageResizeByMaxHeight Error maxWidth(" + strconv.Itoa(maxHeight) + ")> sourceWidth(" + strconv.Itoa(sourceHeight) + ")")
	}
	scalingRatio := float64(maxHeight) / float64(sourceHeight)
	width := int(float64(image.Bounds().Dx()) * scalingRatio)
	image = imaging.Resize(image, width, maxHeight, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, image, imaging.JPEG)
	if err != nil {
		return nil, errors.New("imaging.Encode() Error")
	}
	return buf2.Bytes(), nil
}

// ImageResizeByHeight 根据一个固定 Height 缩放图片
func ImageResizeByHeight(loadedImage []byte, height int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	sourceHeight := image.Bounds().Dy()
	scalingRatio := float64(height) / float64(sourceHeight)
	width := int(float64(image.Bounds().Dx()) * scalingRatio)
	image = imaging.Resize(image, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, image, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageResize 重设图片分辨率
func ImageResize(loadedImage []byte, width int, height int) []byte {
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	//生成缩略图，尺寸width*height
	image = imaging.Resize(image, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, image, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageThumbnail 根据设定的图片大小,剪裁图片
func ImageThumbnail(loadedImage []byte, width int, height int) []byte {
	buf := bytes.NewBuffer(loadedImage)
	imageData, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	//生成缩略图，尺寸width*height
	imageData = imaging.Thumbnail(imageData, width, height, imaging.Lanczos)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, imageData, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageAutoCrop  自动裁白边
func ImageAutoCrop(loadedImage []byte, energyThreshold float32) []byte {
	////读取本地文件，本地文件尺寸300*400
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	img, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	//使用 BoundsForThreshold 查找图像的自动裁剪边界
	//croppedBounds := autocrop.BoundsForThreshold(image, energyThreshold/100)

	nRGBAImg := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(nRGBAImg, nRGBAImg.Bounds(), img, img.Bounds().Min, draw.Src)
	result := autocrop.ToThreshold(nRGBAImg, energyThreshold/100)
	//如果不需要边界，可以使用ToThreshold函数方便地获得裁剪图像
	//croppedImg := autocrop.ToThreshold(image, energyThreshold)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, result, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// ImageGray 转换为黑白图片
func ImageGray(loadedImage []byte) []byte {
	////读取本地文件，本地文件尺寸300*400
	//loadedImage, _ := ioutil.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(loadedImage)
	img, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return loadedImage
	}
	result := imaging.Grayscale(img)
	//如果不需要边界，可以使用ToThreshold函数方便地获得裁剪图像
	//croppedImg := autocrop.ToThreshold(image, energyThreshold)
	buf2 := &bytes.Buffer{}
	//将图片编码成jpeg
	err = imaging.Encode(buf2, result, imaging.JPEG)
	if err != nil {
		return loadedImage
	}
	return buf2.Bytes()
}

// GetContentTypeByFileName https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
func GetContentTypeByFileName(fileName string) (contentType string) {
	ext := strings.ToLower(path.Ext(fileName))
	switch {
	case ext == ".png":
		contentType = "image/png"
	case ext == ".jpg" || ext == ".jpeg":
		contentType = "image/jpeg"
	case ext == ".webp":
		contentType = "image/webp"
	case ext == ".gif":
		contentType = "image/gif"
	case ext == ".bmp":
		contentType = "image/bmp"
	case ext == ".heif":
		contentType = "image/heif"
	case ext == ".ico":
		contentType = "image/image/vnd.microsoft.icon"
	case ext == ".zip":
		contentType = "application/zip"
	case ext == ".rar":
		contentType = "application/x-rar-compressed"
	case ext == ".pdf":
		contentType = "application/pdf"
	case ext == ".txt":
		contentType = "text/plain"
	case ext == ".tar":
		contentType = "application/x-tar"
	case ext == ".epub":
		contentType = "application/epub+zip"
	default:
		contentType = "application/octet-stream"
	}
	return contentType
}

// DetectUTF8 检测 s 是否为有效的 UTF-8 字符串，以及该字符串是否必须被视为 UTF-8 编码（即，不兼容CP-437、ASCII 或任何其他常见编码）。
//来自： go\src\archive\zip\reader.go
func DetectUTF8(s string) (valid, require bool) {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		// Officially, ZIP uses CP-437, but many readers use the system's
		// local character encoding. Most encoding are compatible with a large
		// subset of CP-437, which itself is ASCII-like.
		//
		// Forbid 0x7e and 0x5c since EUC-KR and Shift-JIS replace those
		// characters with localized currency and overline characters.
		if r < 0x20 || r > 0x7d || r == 0x5c {
			if !utf8.ValidRune(r) || (r == utf8.RuneError && size == 1) {
				return false, false
			}
			require = true
		}
	}
	return true, require
}

// PrintAllReaderURL 打印阅读链接
func PrintAllReaderURL(Port int, OpenBrowserFlag bool, EnableFrpcServer bool, PrintAllIP bool, ServerHost string, ServerAddr string, FrpRemotePort int, DisableLAN bool, enableTls bool, etcStr string) {
	protocol := "http://"
	if enableTls {
		protocol = "https://"
	}
	localURL := protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr
	fmt.Println(locale.GetString("local_reading") + localURL + etcStr)
	//PrintQRCode(localURL)
	//打开浏览器
	if OpenBrowserFlag {
		OpenBrowser(protocol + "127.0.0.1:" + strconv.Itoa(Port) + etcStr)
		if EnableFrpcServer {
			OpenBrowser(protocol + ServerAddr + ":" + strconv.Itoa(FrpRemotePort) + etcStr)
		}
	}
	if !DisableLAN {
		printURLAndQRCode(Port, EnableFrpcServer, PrintAllIP, ServerHost, ServerAddr, FrpRemotePort, protocol, etcStr)
	}
}

func printURLAndQRCode(port int, EnableFrpcServer bool, PrintAllIP bool, ServerHost string, ServerAddr string, FrpRemotePort int, protocol string, etcStr string) {
	//启用Frp的时候
	if EnableFrpcServer {
		readURL := protocol + ServerAddr + ":" + strconv.Itoa(FrpRemotePort) + etcStr
		fmt.Println(locale.GetString("frp_reading_url_is") + readURL)
		PrintQRCode(readURL)
	}
	if ServerHost != "" {
		readURL := protocol + ServerHost + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
		return
	}
	//打印所有可用网卡IP
	if PrintAllIP {
		IPList, err := GetIPList()
		if err != nil {
			fmt.Printf(locale.GetString("get_ip_error")+" %v", err)
		}
		for _, IP := range IPList {
			readURL := protocol + IP + ":" + strconv.Itoa(port) + etcStr
			fmt.Println(locale.GetString("reading_url_maybe") + readURL)
			PrintQRCode(readURL)
		}
	} else {
		//只打印本机的首选出站IP
		OutIP := GetOutboundIP().String()
		readURL := protocol + OutIP + ":" + strconv.Itoa(port) + etcStr
		fmt.Println(locale.GetString("reading_url_maybe") + readURL)
		PrintQRCode(readURL)
	}
}

func PrintQRCode(text string) {
	obj := qrcodeTerminal.New()
	obj.Get(text).Print()
}

// CheckPort 检测端口是否可用
func CheckPort(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, locale.GetString("cannot_listen")+"%q: %s", port, err)
		if err != nil {
			return false
		}
		return false
	}
	err = ln.Close()
	if err != nil {
		fmt.Println(locale.GetString("check_pork_error") + strconv.Itoa(port))
	}
	//fmt.Printf("TCP Port %q is available", port)
	return true
}

// GetIPList 获取本机IP列表
func GetIPList() (IPList []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if i.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Printf(locale.GetString("get_ip_error")+"%v", err)
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			IPList = append(IPList, ip.String())
		}
	}
	return IPList, err
}

// GetOutboundIP 获取本机的首选出站IP
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

//// 获取mac地址列表,暂时用不着
//func GetMacAddrList() (macAddrList []string) {
//	netInterfaces, err := net.Interfaces()
//	if err != nil {
//		fmt.Printf(locale.GetString("check_mac_error")+": %v", err)
//		return macAddrList
//	}
//	//for _, netInterface := range netInterfaces {
//	//	macAddr := netInterface.HardwareAddr.String()
//	//	if len(macAddr) == 0 {
//	//		continue
//	//	}
//	//	macAddrList = append(macAddrList, macAddr)
//	//}
//	for _, netInterface := range netInterfaces {
//		flags := netInterface.Flags.String()
//		if strings.Contains(flags, "up") && strings.Contains(flags, "broadcast") {
//			macAddrList = append(macAddrList, netInterface.HardwareAddr.String())
//		}
//	}
//	return macAddrList
//}

// ChickExists 判断所给路径文件或文件夹是否存在
func ChickExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// ChickIsDir 判断所给路径是否为文件夹
func ChickIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// OpenBrowser 打开浏览器
func OpenBrowser(uri string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("CMD", "/C", "start", uri)
		if err := cmd.Start(); err != nil {
			fmt.Println(locale.GetString("open_browser_error"))
			fmt.Println(err.Error())
		}
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", uri)
		if err := cmd.Start(); err != nil {
			fmt.Println(locale.GetString("open_browser_error"))
			fmt.Println(err.Error())
		}
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("xdg-open", uri)
	}
}

// MD5file 计算文件MD5
func MD5file(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		log.Fatal(e)
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}

//从字符串中提取数字,如果有几个数字，就简单地加起来
func getNumberFromString(s string) (int, error) {
	var err error
	num := 0
	//同时设定倒计时秒数
	valid := regexp.MustCompile("[0-9]+")
	numbers := valid.FindAllStringSubmatch(s, -1)
	if len(numbers) > 0 {
		//循环取出多维数组
		for _, value := range numbers {
			for _, v := range value {
				temp, errTemp := strconv.Atoi(v)
				if errTemp != nil {
					fmt.Println("error num value:" + v)
				} else {
					num = num + temp
				}
			}
		}
		//fmt.Println("get Number:",num," form string:",s,"numbers[]=",numbers)
	} else {
		err = errors.New("number not found")
		return 0, err
	}
	return num, err
}

//检测字符串中是否有关键字
func haveKeyWord(checkString string, list []string) bool {
	//转换为小写，使Sketch、DOUBLE也生效
	checkString = strings.ToLower(checkString)
	for _, key := range list {
		if strings.Contains(checkString, key) {
			return true
		}
	}
	return false
}
