class book {
  String title;
  String author;
  String description;
  String image;
  String price;
  String rating;
  String url;

  book({required this.title, this.author, this.description, this.image, this.price, this.rating, this.url});

  factory book.fromJson(Map<String, dynamic> json) {
    return book(
      title: json['title'],
      author: json['author'],
      description: json['description'],
      image: json['image'],
      price: json['price'],
      rating: json['rating'],
      url: json['url'],
    );
  }


}