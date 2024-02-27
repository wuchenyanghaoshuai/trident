```yaml
这个是直接从victoriaMetrics中使用promesql来获取数据


postnam直接发送post请求就可以了，


  {
    "sql":100-"cpu_usage_idle{ident=~\"dev-k8s-master-192.168.3.140\",cpu=\"cpu-total\"}"
  }
这个demo是获取cpu的使用率
由于不会前端，所以没办法展示出来，但是可以直接在postman中发送请求，就可以看到返回的数据了
```

