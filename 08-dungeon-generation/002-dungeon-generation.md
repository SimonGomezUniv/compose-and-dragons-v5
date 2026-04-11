---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
</style>
# 🏰 Dungeon Generation | **With AI**

```golang
dungeonGeneratorAgent, err := structured.NewAgent[[]GeneratedRoom](...)

dungeonGeneratorAgent.AddMessage(
	roles.System, `DUNGEON MAP:\n`+dungeonMap+`\n`,
)
dungeonGeneratorAgent.AddMessage(
	roles.System, `DUNGEON ROOMS CONNECTIONS LIST:\n`+dungeonRoomsConnections+`\n`,
)
```


```golang
generatedRooms, finishReason, err := dungeonGeneratorAgent.GenerateStructuredData(...)
```



