import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'pages.dart';

// 在Flutter中发起HTTP网络请求 https://doc.flutterchina.club/networking/

Future<List<Book>> fetchBooks() async {
  final dio = Dio();
  final prefs = await SharedPreferences.getInstance();
  final comigoHost = prefs.getString('comigo_host') ?? "http://192.168.3.15:1234";
  var url = '$comigoHost/api/book_infos?depth=1&sort_by=name';
  final response = await dio.get(url);
  if (response.statusCode == 200) {
    try {
      final data = response.data as List<dynamic>;
      final books = data.map((e) => Book.fromJson(e)).toList();
      return books;
    } catch (e) {
      debugPrint(e.toString());
      rethrow;
    }
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
  PageInfo? cover;
  List<PageInfo>? pages;

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
      cover: json['cover'] != null ? PageInfo.fromJson(json['cover']) : null,
      pages: json['pages'] != null ? (json['pages'] as List).map((i) => PageInfo.fromJson(i)).toList() : null,
      pageCount: json['page_count'] ?? 0,
      childBookNum: json['child_book_num'] ?? 0,
    );
  }
}
