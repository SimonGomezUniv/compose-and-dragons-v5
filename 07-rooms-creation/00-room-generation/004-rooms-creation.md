<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
</style>
# Room<span class="indianred">S</span> Creation: <span class="dodgerblue">Structured Output</span>

```golang
type Item struct {
	Type     ItemType `json:"type"`
	Quantity int      `json:"quantity"`
}

type Room struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Items            []Item `json:"items"`
}
```
> **`structured.NewAgent[[]Room]`**