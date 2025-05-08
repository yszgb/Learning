// Path: dart/main.go

// dart 基础知识

void main() {
  print('Hello World!');
  dartVariables();

  dartClasses();

  dartCollections();
}

void dartVariables() {
  // Variables
  int age = 25;
  String name = 'John';

  print(name);
  print(age);

  var isStudent = true;
  print(isStudent);

  var marks = 95;
  if (marks > 90) {
    print('A+');
  } else if (marks > 80) {
    print('A');
  }

  var grade = 'A+';
  switch (grade) {
    case 'A+':
      print('Excellent');
      break;
    case 'A':
      print('Good');
      break;
    default:
      print('Fail');
  }
}

class Person {
  String name;
  int age;
  Person(this.name, this.age);
  // Person(required this.name, required this.age)
  // 使用 required，当传入参数时，必须指定 name 和 age 的标签

  String describe() => 'My Name is $name, and I am $age years old.';
}

void dartClasses() {
  var person = Person('John', 25);
  print(person.describe());
}

void dartCollections() {
  // List
  var fruits = ['Apple', 'Banana', 'Orange']; // dart '' "" 一样，字符用 rune
  print(fruits);
  print(fruits.length);
  for (var fruit in fruits) {
    print(fruit);
  }
  // 循环和直接打印的区别？
  // 循环会打印出索引，直接打印会打印出值 -> [Apple, Banana, Orange]

  // Set
  var countries = {'China', 'India', 'USA'};
  for (var country in countries) {
    print(country);
  }

  // Map
  var person = {'name': 'John', 'age': 25};
  print(person['name']);
  print(person['age']);
}
