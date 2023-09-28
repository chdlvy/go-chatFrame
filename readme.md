1. 克隆项目
```
git clone git@github.com:chdlvy/go-chatFrame.git
```

2. 根据自身情况配置config目录下的config.yaml
3. 试着运行example



### 路由接口文档

1. 创建用户

+ URL:"/createUser"

+ 请求方法：POST

+ 请求体：

  ```json
  {
      "NickName":"hhh",
      "FaceUL":"/hh.png"
  }
  ```
+ 响应：
  ```json
   {
      "status": "success/failed",
      "data": {user信息}
    }
  ```
2. 创建群聊

+ URL：“/createGroup”

+ 请求方法：POST

+ 请求体：

  ```
  {
      "gorupName":"newGroup",
      "introduction":"xxxx",
      "faceURL":"xxx.png",
      "creatorUserID":123
      "groupType":GroupType
  }
  ```

3. 加入群聊

+ URL："/joinGroup"

+ Method：POST

+ body：

  ```
  {
      UserID:     2,
      NickName:   "xxx",
      FaceURL:    "xxx.png",
  }
  ```

  





### API

#### config包
**InitConfig(configName, configFolderPath string) error**

初始化配置文件，程序开始之前必须先初始化配置文件，配置文件模板如下
```yaml
mysql:
  address: 127.0.0.1:3306
  username: root
  password: 123456
  database: chatIM
  maxOpenConn: 1000
  maxIdleConn: 100
  maxLifeTime: 60
  logLevel: 4
  slowThreshold: 500

redis:
  address:
    - 127.0.0.1:6379
  username:
  password:
```
