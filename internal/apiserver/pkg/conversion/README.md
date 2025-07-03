### Conversion包

conversion实现不同层的数据转换

miniblog采用了三层架构: Handler, Biz, Store

Handler层的数据在pkg/api/apiserver/v1/中, 结构体类型的命名格式为CreateXyzRequest, CreateXyzResponse

Store层的数据定义在internal/apiserver/model/中, 命名格式为XyzM, 例如PostM

Biz层的数据类型和Handler层的数据类型保持一致

不同层之间通信时, 需要进行数据类型转换, 将Biz层的数据转换为Store层的数据类型

miniblog项目的数据类型转换统一在internal/apiserver/pkg/conversion目录中实现, 转换方法名遵守: 

<资源名>ModelTo<资源名><版本号> <资源名><版本号>To<资源名>Model

不同层之间的数据类型转换都在同一个coversion包中实现, 需要避免出现循环依赖

可以将不同层之间的数据类型转换函数都定义在独立的包中来避免循环依赖