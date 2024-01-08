import 'image_info.dart';

class Book {
  String title;
  String id;
  String type;
  int pageCount;
  int childBookNum;
  ImageInfo? cover;
  List<ImageInfo>? pages;

  Book(
      {required this.title,
        required this.id,
        required this.type,
        this.pageCount = 0,
        this.childBookNum = 0,
        this.cover,
        this.pages});

  factory Book.fromJson(Map<String, dynamic> json) {
    return Book(
      title: json['title'] as String,
      type: json['author'] as String,
      id: json['id'] as String,
      cover: json['cover'] != null ? ImageInfo.fromJson(json['cover']) : null,
      pages: json['pages'] != null ? (json['pages'] as List).map((i) => ImageInfo.fromJson(i)).toList() : null,
      pageCount: json['page_count'] ?? 0,
      childBookNum: json['child_book_num'] ?? 0,
    );
  }
}
