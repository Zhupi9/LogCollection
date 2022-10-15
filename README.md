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
![架构设计](Resource/LogCollection架构设计.png)

## 主要组件介绍
LogAgent： 日志收集客户端，用来收集服务器上的日志。
Kafka：高吞吐量的分布式队列（Link开发，apache顶级开源项目）。
ElasticSearch：开源的搜索引擎，提供基于HTTP RESTful的web接口。
Kibana：开源的ES数据分析和可视化工具。
Hadoop：分布式计算框架，能够对大量数据进行分布式处理的平台。
Storm：一个免费且开源的分布式实时计算系统

## 学习到的技能
1. 服务端agent开发
2. 后段服务组件开发
3. Kafka和zookeeper的使用
4. ES和Kibana的使用
5. etcd的使用

## Kafka介绍
基础结构：producer -> cluster -> consumer
### 1. Kafka集群架构
· broker
· topic
· partition：leader和follower
· replication
### 2. Kafka发送数据的流程（6步）
1. 生产者从Kafka集群获取分区leader信息
2. 生产者将消息发送给leader
3. leader将消息存储在本地磁盘
4. follower从leader拉取消息
5. follower向leader回复ack
6. leader收到所有ack后向producer发送ack
### 3. Kafka选择分区的模式（3种）
1. 指定往哪个分区写
2. 给定key，根据key做hash，决定写哪个分区
3. 轮询方式
### 4.生产者往kafka发送数据的模式（3种）
1. 0：把数据发送给leader就结束
2. 1: 把数据发送给leader等待leader的ack
3. 2: 把数据发送给leader，follower拉去数据回复ack给leader，leader再回复ack

