# Seckill
e-commerce secKill application
the application is highly decoupled,there are three main modules that includes proxy module、layer module and web module.
the system uses redis and etcd，High availability of the system is achieved.

in proxy module, wo achive unitive url routing control, in layer module, we use bucket algorithm to achive fine-grained flow control,
in web module, we achive basic and kernel business include putting goods in sale,goods repertory managment, goods secKill and so on.
