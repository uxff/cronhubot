# 分布式定时任务服务

分布式定时任务系统需要至少三个节点，数据库使用Mysql，自动选举master，职责如下：

|职责\节点|主节点|从节点|
|:----|:----|:----|
|ping nodes|Y|Y|
|检查master超时|-|Y|
|提供接口操作CRON(CRUD)|Y|Y|
|运行CRON|Y|Y|
|增删改通知MASTER|Y|Y|
|分配CRON|Y|-|
|检查CRON分配状态|Y|-|

