# 日志收集项目架构设计

## 项目背景
每个业务都有日志，当系统出现问题时，需要通过日志信息来定位和解决问题。当系统机器比较少时，登录到服务器上查看即可满足；当系统及其规模巨大，登录到机器上查看几乎不现实（分布式系统，一个系统部署在十几台机器上）
## 解决方案
把机器上的日志实时收集，统一存储到中心系统。再对这些日志建立索引，通过搜索即可快速找到对应的日志记录。通过提供一个界面友好的web页面实现日志展示与检索
## 业界解决方案-ELK


## ELK方案的问题
1.运维成本高，每增加一个日志收集项，都需要手动修改配置
2.监控缺失，无法准确获取logstash状态
3.无法做到定制化开发与维护

## 架构设计
<img src="/Users/mengxuanping/Documents/思维导图/LogCollection架构设计.png" alt="LogCollection架构设计" style="zoom:50%;"/>