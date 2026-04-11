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
type GeneratedRoom struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Items            []Item `json:"items"`
}
```


```golang
dng.Generate(nbDungeonRooms)
dungeonMap := dng.GetBWDetailedGridCompact()
dungeonRoomsConnections := dng.GetDungeonDescriptionText()
rooms := dng.GetRooms()
```



