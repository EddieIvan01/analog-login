# Analog_Login

CUMT教务系统模拟登录，代码对post数据进行rsa加密后提交到服务器端

import login后继承Loginer父类，依次调用函数获取cookie后即可使用子类self.sessions进行后续操作

脚本中含有成绩获取代码

截止到2018.3.28教务系统版本可用，未来教务系统改版会持续修改代码


## 2019/8/1

因项目需要实现了Go版本登录代码, 用法同Python库，登录后直接调用`Loginer.S`。因静态语言类型精度问题，Go的实现版本中用32 bytes一组的循环big integer运算生成RSA公钥

至于原来的Python代码，越看越觉得like piece of shit :) 其实自己实现的`hex_base64`和那个大数RSA库都不是必需的，只是我写该库的时候太菜不知如何解决而已

但没什么大BUG也不想重写了，唯一需要注意的是，当用该库进行多用户登录时需修改requests.Session为实例成员
