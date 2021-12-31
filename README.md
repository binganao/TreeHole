# 饼干树洞

## 开发背景

写这个应用的初衷是让用户有一个可以说悄悄话的地方。什么是悄悄话呢？我理解为秘密吧，就是那种想找个朋友倾诉但是又害怕朋友将你的秘密公之于众。这个程序就是为了解决这一难题而开发的。

## 程序概况

拥有基本的注册、登录功能以及一个聊天框。

### 精心设计的 UI，简约易用

饼干树洞的UI设计方向为简单，让用户更轻易的学会使用。为了保护用户隐私信息，仅使用邮箱作为注册和登录功能的前置。

![](https://cdn.bingbingzi.cn/blog/20211230215547.png)

### 不可逆加密，给你的秘密上锁

采用几乎无法破解的 MD5算法 进行信息加密，在提交给后端的信息中通过 MD5 和 AES 的巧妙组合，保证每个用户的密钥不可逆且唯一。

![](https://cdn.bingbingzi.cn/blog/20211231004937.png)

### 开源，才是看得见的安全

我认为，当无法证明自己的软件多么注重信息保护的时候，开源就是最好的选择。让程序曝光在大众之下，更加避免了用户信息被开发者窃取的风险。

代码仓库放置在：https://github.com/binganao/TreeHole 等待优化完代码结构后就会上传项目。

## 开发信息

这个程序分为前端和后端，前端便是 Android （JAVA）开发的，后端使用的是 Golang （GIN框架）开发的。这边我分为前端和后端分别讲解。

### 前端

其实我不太会 Android 开发，当时做大作业的时候我是想着使用 VUE + UNIAPP 进行开发的。这样就能一次开发，跨平台使用。但是原生 Android 性能会好很多，所以我还是放弃了这个想法。

#### UI设计

如你所见，饼干树洞的 UI 非常简洁，在设计的时候我参考了 https://dribbble.com/search/shots/popular/mobile?q=chat 的设计版图，这在很大程度上影响了该程序的 UI。

![](https://cdn.bingbingzi.cn/blog/20211231010744.png)

接着就是设置布局，我不是很习惯使用 XML 的形式进行布局，但是使用之后感觉 Android 在这方面可能会更灵活，虽然也变得复杂很多。

![](https://cdn.bingbingzi.cn/blog/20211231010922.png)

该程序中所有的图片资源均在 https://www.vecteezy.com 上面找到的，这些插画我感觉很好看。

![](https://cdn.bingbingzi.cn/blog/20211231011137.png)

#### 编写代码

首先是程序的结构，

![](https://cdn.bingbingzi.cn/blog/20211231011438.png)

受到 MVC 设计模式的影响，我将程序结构分为用来存放 **Activity** 的 **activities** 包，用来存放适配器的 **apater** 包，用来存放模型的 **model** 包，以及一个工具包 **util**。当时命名的时候比较随意忘记了复数形式，在修正后开源。

##### 整体上

我在设计该 App 的时候为了提高用户的沉浸感，将标题栏和状态栏都隐藏了起来，同时为安卓 9 及后续版本提供了适配（安卓 9 开始不允许在主线程中进行网络操作，必须使用异步或者配置策略，我这里使用的方法是配置策略）

```java
// 配置线程策略
if (android.os.Build.VERSION.SDK_INT > 9) {
    StrictMode.ThreadPolicy policy = new StrictMode.ThreadPolicy.Builder().permitAll().build();
    StrictMode.setThreadPolicy(policy);
}
getSupportActionBar().hide(); // 隐藏标题栏
getWindow().setFlags(WindowManager.LayoutParams.FLAG_FULLSCREEN
        , WindowManager.LayoutParams.FLAG_FULLSCREEN); // 隐藏状态栏
```

##### SignInActivity

程序运行后将会进入到这个 Activity 中，他负责的是展示程序用途、提供注册、登录的功能。也就是说需要两个按钮控件，一个负责直接页面跳转，一个负责判断是否登录成功后跳转。

###### 注册按钮

```java
binding.textCreateNewAccount.setOnClickListener(v ->
                startActivity(new Intent(getApplicationContext(), SignUpActivity.class)));
```

通过 **databinding**，将控件绑定后设置监听，当按下后启动新 **Activity**，即注册页面。

###### 登录按钮

首先进行的是输入判断

```java
if (!isEmail(binding.inputEmail.getText().toString())) {
    Toast.makeText(getApplicationContext(), "你的邮箱格式不正确?", Toast.LENGTH_LONG).show();
    return;

// 判断长度
if (binding.inputPassword.getText().toString().length() < 6) {
    Toast.makeText(getApplicationContext(), "密码长度应大于等于6?", Toast.LENGTH_LONG).show();
    return;
}
```

其中 **isEmail** 使用正则表达式进行判断，如果匹配到则返回 **true**，否则返回 **false**。

```java
public static boolean isEmail(String strEmail) {
    String strPattern = "^\\s*\\w+(?:\\.{0,1}[\\w-]+)*@[a-zA-Z0-9]+(?:[-.][a-zA-Z0-9]+)*\\.[a-zA-Z]+\\s*$";
    Pattern p = Pattern.compile(strPattern);
    Matcher m = p.matcher(strEmail);
    return m.matches();
}
```

完成判断后进行登录操作，也就是访问后端服务端提供的 **api**，如果返回的 **Code**，为 **200**，则登录成功，否则登录失败。这边的登录成功和失败我就使用了一个 **Toast** 来输出，比较方便点也符合我的程序 UI。

这边有个不足的点，那就是直接采用硬编码数据，我将服务器的 ip 直接写在了程序里面而不是使用 xml 储存，这样等到程序变大之后，只要修改了服务器 ip 信息，就得改很多地方这样不好。其他地方的话应该是可以将登录写成一个方法，这样可以提高代码的可读性。还有调试的数据使用了 **System.out.println**，可以使用更方便的 **Log** 包，可以更好的管理输出的调试信息。

##### SignUpActivity

这是注册页面的具体功能实现，首先依旧是判断输入信息

```java
// 判断填写完整
if (binding.inputEmail.getText().toString().isEmpty() ||
        binding.inputPassword.getText().toString().isEmpty() ||
        binding.inputConfirmPassword.getText().toString().isEmpty()) {
    Toast.makeText(getApplicationContext(), "你好像还有位置没有填写完整?", Toast.LENGTH_LONG).show();
    return;

// 判断格式
if (!isEmail(binding.inputEmail.getText().toString())) {
    Toast.makeText(getApplicationContext(), "你的邮箱格式不正确?", Toast.LENGTH_LONG).show();
    return;

// 判断长度
if (binding.inputPassword.getText().toString().length() < 6) {
    Toast.makeText(getApplicationContext(), "密码长度应大于等于6?", Toast.LENGTH_LONG).show();
    return;
}
  
// 判断确认密码和密码正确
if (!binding.inputPassword
        .getText()
        .toString()
        .equals(binding.inputConfirmPassword
                .getText()
                .toString()
        )) {
    Toast.makeText(getApplicationContext(), "两次密码输入不同", Toast.LENGTH_LONG).show();
    return;
}
```

其实注册的功能实现和登录基本没有差别，这里的不足也基本相同，不在叙述。

##### SecretListActivity

为了更好的定制样式，我选用的是 **RecyclerView**，这个算是一个比较高级的列表，能够从底部向上插入数据。为了给用户带来一种如同在对话的感觉，特地设计的聊天风格。

其中添加悄悄话的功能实现为

```java
if(s != null) {
        secretsList.add(s);
    binding.inputSecret.setText(null);
    if (secretsList.size() == 0) {
            secretAdapter.notifyDataSetChanged();
    } else {
            secretAdapter.notifyItemRangeInserted(secretsList.size(), secretsList.size());
        binding.secretRecyclerView.smoothScrollToPosition(secretsList.size() - 1);
    }
    binding.secretRecyclerView.setVisibility(View.VISIBLE);
}
```

这边的逻辑跟登录和注册有所不同，在页面初始化的时候会进行一个获取悄悄话的功能，并且通过一个 for 循环将内容输出到列表中。

##### AES

本加密参考 https://blog.csdn.net/baidu_27419681/article/details/61206149 不反复造轮子了。

##### 总结

对于前端代码，有很多可以优化的地方，比如 ip 地址硬编码问题，然后就是代码写的不是很优雅，我有一个这样的想法，就是创建一个 **Context** 类，用来存放这些数据信息，然后通过 **static** 来保存，之后每次需要调用就写类似下方的代码。

```java
Toast.makeText(getApplicationContext(), Context.PROCESS_ERROR, Toast.LENGTH_LONG).show();
```

还有使用了一个我不太安全的包 **FastJson**，他在我的程序中的作用是将 Json 字符串转化为对象，但是据说它在之前的版本中爆出过很多 **RCE** 漏洞，这可能会影响到我的程序。

### 后端

这个后端基本是实现了个功能，并没有非常正式的去写。比如没有将路由功能分包写，所有的函数都写在了一个 go 文件中。没有将配置信息放在一个 **config.yml** 或者 **config.json** 中，不利于修改服务端配置。

#### 路由部分代码

```go
r.POST("/api/login", func(ctx *gin.Context) {
	encUserName, _ := ctx.GetPostForm("encUserName")
	encPassWord, _ := ctx.GetPostForm("encPassWord")
	if verifyLogin(encUserName, encPassWord, db) {
		ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
	} else {
		ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
	}
})

r.POST("/api/register", func(ctx *gin.Context) {
	encUserName, _ := ctx.GetPostForm("encUserName")
	encPassWord, _ := ctx.GetPostForm("encPassWord")
	if doRegister(encUserName, encPassWord, db) {
		ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
	} else {
		ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
	}
})

r.POST("/api/secretlist", func(ctx *gin.Context) {
	encUserName, _ := ctx.GetPostForm("encUserName")
	encPassWord, _ := ctx.GetPostForm("encPassWord")
	if b, data := getEncSecretList(encUserName, encPassWord, db); b {
		ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}, User: data})
	} else {
		ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
	}
})

r.POST("/api/addsecret", func(ctx *gin.Context) {
	encUserName, _ := ctx.GetPostForm("encUserName")
	encPassWord, _ := ctx.GetPostForm("encPassWord")
	encSecret, _ := ctx.GetPostForm("encSecret")
	if addSecret(encUserName, encPassWord, encSecret, db) {
		ctx.JSON(http.StatusOK, Resp{Code: "200", Data: Datas{Msg: "success"}})
	} else {
		ctx.JSON(http.StatusInternalServerError, Resp{Code: "500", Data: Datas{Msg: "fail"}})
	}
})
```

服务端提供了四个路由，分别为

> /api/login 提供登录
>
> /api/register 提供注册
>
> /api/secretlist 返回当前用户的所有悄悄话
>
> /api/addsecret 增加一段悄悄话

#### 具体功能代码

```go
func verifyLogin(encUserName, encPassWord string, db *gorm.DB) bool {
	var user User

	db.Where("enc_user_name = ? and enc_pass_word = ?", encUserName, encPassWord).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

func doRegister(encUserName, encPassWord string, db *gorm.DB) bool {
	fmt.Println(encUserName)

	user := User{
		EncUserName: encUserName,
		EncPassWord: encPassWord,
		EncSecret: "",
	}

	fmt.Println(user)

	if err := db.Create(&user).Error; err == nil{
		return true
	}

	return false
}

func getEncSecretList(encUserName, encPassWord string, db *gorm.DB) (bool, []User) {
	var data []User

	db.Where("enc_user_name = ? and enc_pass_word = ?", encUserName, encPassWord).Not("enc_secret = ?", "").Find(&data)
	if len(data) >= 0 {
		return true, data
	}


	return false, []User{}
}

func addSecret(encUserName, encPassWord, encSecret string, db *gorm.DB) bool {
	user := User{
		EncUserName: encUserName,
		EncPassWord: encPassWord,
		EncSecret: encSecret,
	}

	if err := db.Create(&user).Error; err == nil{
		return true
	}

	return false
}
```

##### 登录方面

其实是对数据库进行了查询，如果存在用户名和密码的数据，便返回 **true** 代表登录成功

##### 注册方面

实现了数据库的增加操作，需要两个参数，分别是用户名和密码

##### 获取悄悄话列表

对数据库进行了查询，并且返回所有不为空的悄悄话

##### 增加悄悄话

实现了数据库的增加操作，与注册不同的是需要将悄悄话作为参数传递进来

#### 安全性

为了防止黑客直接对后端服务器中的 **api** 进行攻击，在执行数据库请求的时候会对数据进行预编译，从而从根本上阻断了 **SQL注入** 的可能性。

与此同时，采用了 **Golang 1.17** 版本，也能解决 Go 语言在 1.16 及以下版本中的拒绝服务漏洞

#### 部署简便

后端语言有很多种选择，Java、Python、PHP、ASP、NODE、GO 等等，但是 Java、PHP、ASP 后端需要配置环境，Python部署时需要安装第三方类库，而 NODE 的后端采用了 JavaScript 作为解释语言性能并不出众且部署不便。出于方便部署的角度，我选择了 Go 语言，它具有跨平台，交叉编译的优势，非常满足我的需求。

同时使用了 Mysql 作为数据库程序，原先我想用 Redis 作为缓存数据库的，但是发现这个程序好像并不需要有缓存的内容（点赞、浏览量等信息），所以最后只使用了 Mysql 作为数据库程序。

##### 实际部署

##### 后端程序部署

有两种方式推荐，第一种是开启一个 **screen** 保持程序运行，第二种是作为一个 **service** 开机自启，两种方式都比较方便，我习惯使用 **screen** 便以此作为演示。

![](https://cdn.bingbingzi.cn/blog/20211231105805.png)

##### 数据库部署

一共需要执行两部操作，第一是创建数据库，第二是导入 SQL 文件，其中 SQL 文件包含在程序包中。

![](https://cdn.bingbingzi.cn/blog/20211231110109.png)

## 未来

- 添加一个 **Context** 类，将数据信息写入到这里面，并使用统一的命名规范来命名变量。

- 加入记住用户密码功能，其实挺简单的使用 **SharedPreferences** 写个就行

- 将 IP 地址信息写入到 xml 文件中，方便未来维护

- 添加删除悄悄话功能，目前只有添加和查询，有些单调

- 加入分享功能

- 优化代码使程序更快运行

- 加入 Docker 快速后端部署
