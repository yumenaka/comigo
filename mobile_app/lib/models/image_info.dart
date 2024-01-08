class ImageInfo {
  String url;
  String filename;

  ImageInfo({required this.url, required this.filename});

  factory ImageInfo.fromJson(Map<String, dynamic> json) {
    return ImageInfo(
      url: json['url'],
      filename: json['filename'],
    );
  }
}
