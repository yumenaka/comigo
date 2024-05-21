import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'page_info.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

// 在Flutter中发起HTTP网络请求 https://doc.flutterchina.club/networking/
Future<Book> getBook({required String bookID}) async {
  final dio = Dio();
  final comigoHost = dotenv.env['DEFAULT_HOST']!;
  final fakeBookID = dotenv.env['FAKE_BOOK_ID']!;
  print('$comigoHost/api/get_book?id=$bookID');
  var url = '$comigoHost/api/get_book?id=$fakeBookID';
  final response = await dio.get(url);
  if (response.statusCode == 200) {
    try {
      print(response.data); // 添加这行来调试查看数据结构
      final data = response.data as Map<String, dynamic>;
      var book = Book.fromJson(data);
      print(book);
      return book;
    } catch (e) {
      debugPrint(e.toString());
      rethrow;
    }
  } else {
    throw Exception('Failed to load books');
  }
}

Future<List<Book>> getBookList() async {
  final dio = Dio();
  final prefs = await SharedPreferences.getInstance();
  final comigoHost = dotenv.env['DEFAULT_HOST']!;
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

// Book.fromJson就是一个命名构造函数。可以使用该构造函数从Map中生成一个Student对象，有点像是java中的工厂方法。
// 命名构造函数的格式是ClassName.identifier
// 普通构造函数是没有返回值，而factory构造函数需要一个返回值。
  factory Book.fromJson(Map<String, dynamic> json) {
    // 解析pages字段
    List<PageInfo> pagesList;
    if(json['pages'] != null){
      pagesList = (json['pages']['images'] as List)
          .map((i) => PageInfo.fromJson(i))
          .toList();
    }else{
      pagesList = [];
    }
    return Book(
      title: json['title'] as String,
      type: json['author'] as String,
      id: json['id'] as String,
      cover: json['cover'] != null ? PageInfo.fromJson(json['cover']) : null,
      pages: json['pages'] != null ? pagesList : null,
      pageCount: json['page_count'] ?? 0,
      childBookNum: json['child_book_num'] ?? 0,
    );
  }
}
