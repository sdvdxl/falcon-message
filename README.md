# falcon-message
open-falcon alarm 发送消息组件

现在实现了钉钉群消息和微信消息

# 微信
需要在通讯录的人员名单中的IM配置处配置微信名字。
配置说明请参考 https://github.com/Yanjunhui/chat 这里，代码也是从这里复制粘贴，进行了适当修改，以适应当前程序。

# 钉钉群
1. 钉钉消息是发送到某个群，而不是针对单个人发送，所以需要在这个群中设置一个机器人，定义的时候，选择自定义机器人，然后将webhook链接中access_token的值拷贝出来，以备待用。
2. 在falcon dashaboar 用户管理中心新建一个用户，填写email，然后在 IM 处填写 `[ding]:access_token`，这里的access_token就是上面的access_token，保存用户信息。
3. 在dashboard的群组管理中心新建一个群组，把上面的这个用户加入到这个群组。
4. 在要告警的地方把上面的这个群组加入即可。
