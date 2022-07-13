# 企业微信会话存档服务

## 实现功能：

- [x] 不同版本私钥管理
- [x] 消息存储
  - [x] 数据库
    - [x] mysql
    - [ ] sqlserver
  - [x] 消息队列
    - [x] redis
    - [ ] rabbitmq
- [x] 附件存储
  - [x] 腾讯云
  - [x] 七牛云
- [x] 事件推送触发拉取
- [x] API接口
  - [x] /callback 事件回调(GET|POST)
  - [x] /api/audit/groupchat 获取会话内容存档内部群信息(POST)
  - [x] /api/audit/checkagree 获取会话同意情况(POST)
  - [x] /api/audit/permituser 获取会话内容存档开启成员列表(POST)


## TODO:
- [ ] sqlserver支持
- [ ] 附件存储阿里云、又拍云...
- [ ] Mixed 等类型消息解析
- [ ] 前端界面

## 使用方法

1、修改`.env`文件里的配置后，在任意支持docker的环境里执行
```
docker build . -t msg_audit:latest
```
2、从编译好的镜像启动
```
docker pull golaoji/wecom-dev-audit:latest
docker run -itd --env-file=.env golaoji/wecom-dev-audit:latest
```
