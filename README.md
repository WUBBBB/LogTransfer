# LogTransfer
Use together with LogAgent



# V3.99

为了方便与LogAgent配合使用，二者采用相同版本号表示彼此的对应版本



但是该程序仍有较大问题：

首先，是根据配置文件获得kafka的ip地址以及希望读取的topic

从main函数的连接kafka那个part的kafka.Init函数参数来看，kafka的ip地址是用一个string切片来存的，故**貌似**是可以连接到分布在多台主机上的kafka获取分布存储的信息，这个能力还有待验证。但是，topic只有一个，构思上来讲并不合理，用户可能会使用多个不同的topic来存放不同类型的日志，当前程序的解决办法是多开，针对该问题构想的两种思路：

1、可以依据不同的topic分类，LogTransfer可以根据不同的topic获取到存储有该topic的对应ip地址，想必是极好的

2、也可以设置一个topic切片用于存储所有需要用到的topic，对addr切片的每个ip分别试图读取topic切片的所有topic，这样有好有坏，一方面减少了重复连接到同一主机的次数，但另一方面会造成一些无效查找

<img src="C:\Users\沐\AppData\Roaming\Typora\typora-user-images\image-20220312235320589.png" alt="image-20220312235320589" style="float:left" />

ES的连接及数据传入方面到底有哪些问题还不明了，因为当前对ES的了解还不够，比如Index之下的Type的一些概念和限制

看一下kafka.go的40行注释，似乎这里有个问题没解决？

此外，研究一下es.go代码中对interface{}的使用，这个利用interface实现让chan存储map而不是string从而减少chan的容量和流量压力的技巧需要学习（因为map可以传指针，而传string的话太大了）

最后在kafka.go中还发现了一个可选的改进点，或许可以做一个manager来决定是否从kafka的某个topic往外读数据，实现热修改，即通过更改配置文件，实现立刻新增、删除从kafka中读取数据到chan的某个特定进程

***

### 以上提到的问题的解决与优化，敬请期待v4.0
