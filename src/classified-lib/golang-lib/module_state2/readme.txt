module_state2是module_state的升级版。增加了以下功能：
1. 除了支持counter外，还支持设置string类型的state
2. 支持counter slice，也就是获得两个counter的差值（用于获得
  一段时间内counter的变化情况）
  - 推荐的使用方式：
    使用一个独立的go routine, 周期性的调用CounterSlice的Set()
3. 支持扁平计数结构到对应的层次化技术结构json输出转换


文件列表
- counter.go: 实现Counters
  * Counters定义为map[string]int64
- counter_slice.go：实现CounterSlice
  * CounterSlice用于获得在一定时间内、Counters的两个快照间的变化
- interval.go
  * 实现了NextInterval()，用于判断聚下一周期剩余的时间
- module_state2.go: 实现State
  * 包括counter, 字符串类型的state, 数字类型的state
- counter_convert.go
    * 提供由Counters结构构造多叉树的功能
- counter_hier.go :
    * 提供层次化的counter结构及扁平counter到层次counter的转换操作
- counter_slice_hier.go :
    * 提供扁平的counter slice结构到层次化counter slice结构的json输出转换
- module_state2_hier.go :
    * 提供扁平的module状态统计结构到层次化module状态统计结构的json输出转换
