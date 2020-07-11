# 双人五子棋

## 运行环境

mysql+redis

## 加分项

websocket实现客户端显示信息，redis有序集合存储棋局，分数1和2是两个玩家。redis集合储存房间玩家信息，若找不到玩家则是观战。默认是通过AddRoom加入房间，也可以用url加入，认输通过关键字识别，（求和太复杂没写）

## 接口描述

user/registeer：注册

./login：登录并给username的cookie

chess/addroom：读取username并尝试加入对局，成功给roomid的cookie

./room：使用websocket下棋发言规则是1+...或者2+...（1+发言，2+落子坐标）

## 其他事项

落子限时本来想用数据库存储时间戳实现，但这样操作人数多了数据就很大就没做了