# Introduce

---

简单介绍下功能，设计思路，部分API说明。

## War

- 抽象为战争
- 根据War_Type 维护所有Battle
- `API: Join` 处理加入战斗请求
- `API: Do` 处理操作指令

## Battle

- 抽象为战役
- map + list 维护Conn，对每次使用的Conn，会移到list 末尾用于提高Conn清理效率
- `API: Upgrade` 清理Conn，清除超时的Conn

## Fight

- 抽象为战斗

## Gamer

- 玩家信息
- ``