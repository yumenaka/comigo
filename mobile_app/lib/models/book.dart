import 'image_info.dart';

// 在Flutter中发起HTTP网络请求 https://doc.flutterchina.club/networking/
import 'dart:convert';
import 'package:http/http.dart' as http;

Future<List<Book>> fetchBooks() async {
  final response = await http.get(Uri.parse('http://192.168.3.8:1234/api/book_infos?depth=0&sort_by=name'));

  if (response.statusCode == 200) {
    List<Book> bookshelf = [];
    List<dynamic> booksJson = jsonDecode(response.body);
    for (var bookJson in booksJson) {
      bookshelf.add(Book.fromJson(bookJson));
    }
    return bookshelf;
  } else {
    throw Exception('Failed to load books');
  }
}

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
