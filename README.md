# 图书馆管理系统

## **项目背景![Go](https://github.com/zhangzhengpu/gin_database/workflows/Go/badge.svg)**

这个项目是我选修了数据库原理这一门课，老师向我们提出的实验作业，本来想的就是随便玩玩就可以了，但是我还是要认真做这个，虽说自己还暂时不会写前端，后端技术也不是太好，但是这个项目算是我对后端迈出的第一步，也是我第一个完全自己写的项目

## 使用技术

- **Golang：一门新的程序设计语言**
- **Gin：Golang的一种框架**
- **GORM：对于Golang语言友好的一种开发人员ORM库**
- **ORM：一种对象关系映射，用来将对象和数据库之间的映射的元数据，将面向对象语言程序中的对象自动持久化到关系数据库中。 本质上就是将数据从一种形式转换到另外一种形式。**
- **JSON WEB Token**（**JWT**，读作 [/dʒɒt/]），是一种基于JSON的、用于在网络上声明某种主张的令牌（token）。JWT通常由三部分组成: 头信息（header）, 消息体（payload）和签名（signature）。

## 接口说明

### 用户认证接口

`/api/auth`

`POST` `/register`

输入示例：

```json
{
    "telephone":"00000000000",
    "username":"admin",
    "password":"password",
    "power":x,
    "sex":x,
    "age":xx
}
```

相关示例：

`telephone`必须要求十一位

`username`可以自由输入

`password`必须要求大于六位小于十八位，必须要有大写字母、小写字母、数字和标点符号

`power` 1代表学生，只有借还书和查阅图书功能，2代表图书管理员，可以对图书进行增删改查，3代表是超级管理员，可以对整个系统进行完全控制

`sex`：性别，1为男性，2为女性

`age`年龄，有一定限制，必须为数字

`PUT` `/login`

输入示例：

```json
{
    "username":"xxxxxx",
    "password":"xxxxxx"
}
```

返回示例：

```json
{
    	"code":    200
		"data":    token
		"message": "登录成功"
}		
```

`/api/book` 

`POST` 功能：增加一本书

输入样例：

```json
{
        "isbn":"9897240130102",
		"bookname":"Vue.js从入门到放弃",
		"author":"XXX",
		"press":"清华大学出版社",
		"category" :"科技",
		"getbooknum":"TP193.1.1",
        "position":"三楼科技类书库",
        "bookcode":"0000000000004"
}
```

相关用例如上代码所示

返回示例：

```json
{
    "code": 201,
    "msg": "Add book Successful"
}
```

`PUT` 功能：更新一本书部分信息

```json
{
        "isbn":"9897240130101",
		"bookname":"浑元形意太极门功法大全",
		"author":"xxxxx",
		"press":"清华大学出版社",
		"category" :"人文社科",
		"getbooknum":"TN191-100.4",
        "position":"七楼人文社科类书库",
        "bookcode":"0000000000003"
}
```

返回示例：

```json
{
    "code": 200,
    "id": x,
    "msg": "Book Update OK "
}
```

`DELETE` 功能：删除一本书

```json
{
        "bookcode":"0000000000002"
}
```

返回示例：

```json
{
    "code": 200,
    "msg": "Delete book Successful!"
}
```

`/api/lend` 

`POST` 借阅书籍

输入示例：

```json
{
    	"username":"admin",
        "isbn":"9897240130101",
		"bookname":"浑元形意太极门功法大全",
        "bookcode":"0000000000003"
        
}
```

输出示例：

```json
{
    "code": 200,
    "msg": "borrow book successful!"
}
```

`DELETE` 归还书籍

输入示例

```json
{
    	"username":"admin",
        "isbn":"9897240130101",
		"bookname":"浑元形意太极门功法大全",
        "bookcode":"0000000000003"
        
}
```

输出示例：

```json
{
    "code": 200,
    "id": x,
    "msg": "You Successfully back this book"
}
```

## 安装使用

首先我们在文件当中设置好数据库的位置，用户名密码，然后建立好一个数据库

接下来，我们将我们的文件下载下来

```shell
git clone https://github.com/zhangzhengpu/gin_database
go mod tidy
go build 
默认端口：8080
```

运行生成的文件就可以了（千万注意数据库连接）

感谢**E99p1ant**小哥哥的大力支持

## LICENSE

MIT License

Copyright (c) 2020 zhangzhengpu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.