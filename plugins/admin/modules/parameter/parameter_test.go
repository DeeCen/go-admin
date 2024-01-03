package parameter

import (
    "fmt"
    "testing"
)

func TestGetParamFromUrl(t *testing.T) {
    fmt.Println(GetParamFromURL("/admin/info/user?_p=1&_ps=10&_srt=id&_st=desc",
        1, "asc", "id"))
}

func TestParameters_PKs(t *testing.T) {
    pks := BaseParam().PKs()
    fmt.Println("pks", pks, "len", len(pks))
}
