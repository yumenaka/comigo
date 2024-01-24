class PageInfo {
  String url;
  String filename;

  PageInfo({required this.url, required this.filename});

  factory PageInfo.fromJson(Map<String, dynamic> json) {
    return PageInfo(
      url: json['url'],
      filename: json['filename'],
    );
  }
}
