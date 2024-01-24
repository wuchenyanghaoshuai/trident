```
这个前端代码是为了测试hosts节点的ssh公功能,虽然页有pod的ssh功能但是只能ssh进去不能输入命令，有待解决
如果要使用pod的ssh功能，可以直接使用templete目录下的原生html代码，那个测试完毕，可以使用，根据启动的路由就可以访问
需要修改的地方就是ws的连接地址路径以及端口号,因为我在添加节点的时候已经把密钥拷贝过去了，所以我这边只需要ip+port+username就可以登录了
具体参考的代码为git@gitee.com:KubeSec/webssh.git  版本为v0.1
使用方法
cd webssh
npm install
npm run dev
```