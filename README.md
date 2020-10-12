## HomeProxy
HTTP/HTTPS proxy over SSH.

### QQ群
782391570

## Future
- [x] 支持树莓派3b+, 4b+
- [x] 支持域名定向转发
- [x] 支持替换轮训器,现在随机
- [ ] 支持http Basic auth
- [ ] 支持https 证书
- [ ] 支持控制台需要账号访问

## 更新 2020年10月12日
1. 转发实例的删除和移除
2. 域名定向转发

## 更新 2020年10月11日
1. 转发实例新增优化
2. 转发实例编辑优化
3. 负载均衡器修改
4. 新增完全代理转发模式
5. 前台展示优化，使布局不那么乱
6. 针对转发实例进行存活探测并增加延时展示

## 更新 2020年10月9日
1. 去除sqlite3依赖 
2. 去除其他文件或目录依赖，改成命令行模式

## 手动构建
详情见build.sh

# 使用帮助
## 下载
根据自己的版本下载相应文件[下载](https://gitee.com/Dukeshi/home-proxy/releases)

## 启动进程
```bash
./proxy -h 0.0.0.0:80 -path db -name default
```
## Linux（树莓派） 开机自启 （Centos7以上版本）
修改systemd/system/proxy.service 中 /path/to/proxy 为你部署的具体路径

```bash
cp systemd/system/proxy.service /lib/systemd/system/
systemctl enable proxy
systemct start proxy
```

## 启动代理
![](/pic/start.jpg)

## 新增规则
![](/pic/add_role.jpg)

## 新增转发代理
![](/pic/add_instance.jpg)

## 配置http代理
[帮助](https://jingyan.baidu.com/article/72ee561a053a87e16138dfed.html)