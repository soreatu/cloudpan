南邮网络编程大作业

## 一、项目简介

本项目是一个私人加密云网盘平台，致力于为用户提供一个便捷、高效、安全的云网盘服务。

本项目注重个人隐私安全，不同用户之间的文件设有隔离，确保不会被非法查看、修改、删除。本项目设置了注册登录功能，此外，本项目的云网盘存储服务还默认提供高性能的加密功能，将会对用户上传到云端的文件加密存储，以保证即使在服务器被攻陷后，攻击者在没有拿到密钥的情况下无法对文件进行解密、获取有用信息。

项目采用前后端分离的开发方式。前端整体基于Vue框架搭建，并大量使用了简洁的Element UI组件；后端采用golang作为开发语言，通过gin框架来实现后端api请求处理，利用sessions框架来对用户进行身份识别和权限控制，并使用轻量级的sqlite3作为数据库和gorm作为与数据库交互的工具，可以满足中小范围内的云端存储需求。

## 二、项目实现

前端主要由三个页面构成：首页、用户注册登录页面、文件管理页面。通过Vue设置三个路由"/", "/auth", "/file"来实现这三个页面之间的跳转。

后端主要提供了若干个api，用于用户注册、用户登陆、用户注销、文件上传、文件下载、获取文件等功能。

### 2.1 首页

项目的首页是一个欢迎页面，并提供一个“登陆”按钮，可以跳转至用户注册登陆界面。

![image-20210607131829204](https://soreatu-1300077947.cos.ap-nanjing.myqcloud.com/uPic/image-20210607131829204.png)

### 2.2 用户注册登陆页面

用户登陆界面主要由一个表单和三个按钮组成，表单中可以填入用户名和密码信息。点击“登陆”按钮，会请求后端的登陆api进行用户登陆，如果登陆验证成功，则会调转至该用户的文件管理界面；点击“注册”按钮，会请求后端的注册api进行用户注册，如果注册成功，会提醒重新登陆；点击”返回“按钮，将会跳转至首页。

如果用户名或密码错误，登陆将会失败，会弹出“Wrong username or password!”消息框；如果用户名已经被注册，注册将会失败，会弹出“Username has been used!”消息框。

![image-20210607131854388](https://soreatu-1300077947.cos.ap-nanjing.myqcloud.com/uPic/image-20210607131854388.png)

### 2.3 文件管理页面

文件管理页面由一个表格、一个上传组件和一个“注销”按钮组成。在页面刚开始加载的时候，就会通过后端api去请求当前用户所有保存在云网盘中的文件，并显示在表格中。表格中还提供了“下载”和“删除”两个操作按钮，可以分别用来下载和从云网盘中删除对应的文件。上传组件支持拖拽上传方式，并且有上传进度显示。“注销”按钮将注销当前登陆的用户，并调转至首页。

![image-20210607131906144](https://soreatu-1300077947.cos.ap-nanjing.myqcloud.com/uPic/image-20210607131906144.png)



### 2.4 文件上传与下载

文件通过后端提供的api上传到后端服务器，后端服务器接收到上传的文件后，会使用当前登陆用户的密钥key对文件进行AES-CTR加密，并将加密后的结果通过文件形式保存在服务器的upload文件夹中。此外，后端服务器还会生成一个File记录，用来记录上传文件的文件名、大小和拥有者，并将该File记录存储在数据库中。

当用户点击“下载”按钮时，前端js会通过api发起下载请求时，后端服务器会先检查请求下载的文件是否属于当前登陆用户，若不属于，则会返回“Unauthorized access”错误；若检查通过，后端服务器随后会读取upload文件夹中对应的文件，并使用当前登陆用户的key对其进行AES-CTR解密，并将解密后的字节流（原文件）返回给前端。前端拿到数据后，会弹出一个文件下载框，提示用户保存下载好的文件。

### 2.5 密钥管理

每个用户在注册创建时，都会生成一串16-byte的密钥key，该key将用于加密该用户上传到服务器的所有文件，并作为User用户结构体的一个属性值存储在数据库中。注意，该key并不是明文形式存储在后端数据库中的，而是由一个主密钥master key进行加密后再存储的；在从数据库中取出用户信息时，会由主密钥对加密的key进行解密恢复，这个操作可以借助gorm的BeforeCreate和AfterFind钩子Hook函数来实现。

> 主密钥master key通过读取配置文件.env文件来获取

## 构建

```shell
$ go build -o cloudpan .
$ ./cloudpan
```

前端运行在localhost:8080

后端运行在localhost:8081

## 框架

后端：gin + sessions + sqlite3 + gorm

前端：Vue + Element UI

