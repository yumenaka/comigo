/// 参考文档：
/// 在Flutter中发起HTTP网络请求 https://doc.flutterchina.club/networking/
import 'package:dio/dio.dart';
// debugPrint 是 Flutter 提供的一个用于调试的打印方法，可以在控制台中查看
import 'package:flutter/foundation.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'page_info.dart';

/// 预先配置 Dio 实例
Dio _createDio() {
  final baseUrl = dotenv.env['DEFAULT_HOST'];
  if (baseUrl == null || baseUrl.isEmpty) {
    throw Exception("环境变量 'DEFAULT_HOST' 未配置或为空，请检查 .env 文件。");
  }

  return Dio(
    BaseOptions(
      baseUrl: baseUrl,      // e.g. "https://example.com"
      connectTimeout: const Duration(milliseconds: 5000),  // 5s
      receiveTimeout: const Duration(milliseconds: 5000),
    ),
  );
}

Future<Book> getBook({required String bookID}) async {
  final dio = _createDio();
  try {
    // 使用 queryParameters 来构建参数，避免手动字符串拼接
    final response = await dio.get<Map<String, dynamic>>(
      '/api/get_book',
      queryParameters: {'id': bookID},
    );

    // 检查状态码和响应体
    if (response.statusCode == 200 && response.data != null) {
      debugPrint('getBook response data: ${response.data}'); // 调试用
      return Book.fromJson(response.data!);
    } else {
      throw Exception('Failed to load book with id=$bookID');
    }
  } catch (e, stackTrace) {
    debugPrint('getBook error: $e');
    debugPrint('Stack trace: $stackTrace');
    rethrow; // 将错误继续向上抛
  }
}

Future<List<Book>> getBookList() async {
  final dio = _createDio();
  try {
    // 同样使用 queryParameters
    final response = await dio.get<List<dynamic>>(
      '/api/top_shelf',
      queryParameters: {'sort_by': 'filename'},
    );

    if (response.statusCode == 200 && response.data != null) {
      // 将每个 List item 转换为 Book
      final books = response.data!
          .map((jsonItem) => Book.fromJson(jsonItem as Map<String, dynamic>))
          .toList();
      return books;
    } else {
      throw Exception('Failed to load books');
    }
  } catch (e, stackTrace) {
    debugPrint('getBookList error: $e');
    debugPrint('Stack trace: $stackTrace');
    rethrow;
  }
}

class Book {
  final String title;
  final String id;
  final String type;
  final int? pageCount;
  final int? childBookNum;
  final PageInfo? cover;
  final List<PageInfo>? pages;

  const Book({
    required this.title,
    required this.id,
    required this.type,
    this.pageCount = 0,
    this.childBookNum = 0,
    this.cover,
    this.pages,
  });

  /// 工厂构造函数，用于从 JSON 中解析 Book
  factory Book.fromJson(Map<String, dynamic> json) {
    // 安全地解析 pages
    final pagesJson = json['pages'] as Map<String, dynamic>?;
    final imagesList = pagesJson?['images'] as List<dynamic>? ?? [];
    final pagesParsed = imagesList.map((e) => PageInfo.fromJson(e)).toList();

    return Book(
      title: json['title'] as String? ?? '',
      id: json['id'] as String? ?? '',
      type: json['type'] as String? ?? '',
      cover: json['cover'] != null
          ? PageInfo.fromJson(json['cover'] as Map<String, dynamic>)
          : null,
      pages: pagesParsed.isEmpty ? null : pagesParsed,
      pageCount: (json['page_count'] as int?) ?? 0,
      childBookNum: (json['child_book_num'] as int?) ?? 0,
    );
  }

  @override
  String toString() {
    return 'Book(title: $title, id: $id, type: $type)';
  }
}
