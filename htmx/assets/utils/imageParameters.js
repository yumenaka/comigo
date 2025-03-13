//请求图片文件时，可添加的额外参数
const imageParameters = {
    resize_width: -1, // 缩放图片,指定宽度
    resize_height: -1, // 指定高度,缩放图片
    do_auto_resize: false,
    resize_max_width: 800, //图片宽度大于这个上限时缩小
    resize_max_height: -1, //图片高度大于这个上限时缩小
    do_auto_crop: false,
    auto_crop_num: 1, // 自动切白边阈值,范围是0~100,其实为1就够了
    gray: false, //黑白化
};

//添加各种字符串参数,不需要的话为空
const resize_width_str =
    imageParameters.resize_width > 0
        ? "&resize_width=" + imageParameters.resize_width
        : "";
const resize_height_str =
    imageParameters.resize_height > 0
        ? "&resize_height=" + imageParameters.resize_height
        : "";
const gray_str = imageParameters.gray ? "&gray=true" : "";
const do_auto_resize_str = imageParameters.do_auto_resize
    ? "&resize_max_width=" + imageParameters.resize_max_width
    : "";
const resize_max_height_str =
    imageParameters.resize_max_height > 0
        ? "&resize_max_height=" + imageParameters.resize_max_height
        : "";
const auto_crop_str = imageParameters.do_auto_crop
    ? "&auto_crop=" + imageParameters.auto_crop_num
    : "";

//所有附加的转换参数
let addStr =
    resize_width_str +
    resize_height_str +
    do_auto_resize_str +
    resize_max_height_str +
    auto_crop_str +
    gray_str;

if (addStr!=="") {
    addStr = "?" + addStr.substring(1);
    console.log("addStr:", addStr);
}

export { addStr, imageParameters }; 