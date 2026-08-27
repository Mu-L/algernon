package engine

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xyproto/algernon/lua/luastate"
	"github.com/xyproto/algernon/utils"
	"github.com/xyproto/datablock"
)

// Nested Lua tables should be available to Pongo2 templates, see issue #119
func TestLuaFunctionMapNestedTables(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	filename := "testdata/issue119/index.po2"
	luafilename := "testdata/issue119/data.lua"
	pongodata, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed reading file: %s", err)
	}
	luadata, err := os.ReadFile(luafilename)
	if err != nil {
		t.Fatalf("Failed reading file: %s", err)
	}

	ac := &Config{versionString: "test"}
	ac.SetFileStatCache(datablock.NewFileStat(true, time.Minute*1))
	ac.cache = datablock.NewFileCache(20000000, true, 64*utils.KiB, true, 0)
	ac.luapool = luastate.New()
	defer ac.luapool.Shutdown()

	funcs, err := ac.LuaFunctionMap(w, req, luadata, luafilename)
	if err != nil {
		t.Fatalf("Error with LuaFunctionMap: %s", err)
	}

	// A table that refers to itself must not cause endless recursion
	sl, ok := funcs["cyclic"].([]any)
	if !ok {
		t.Fatalf("cyclic = %T, want []any", funcs["cyclic"])
	}
	if len(sl) != 1 || sl[0] != nil {
		t.Errorf("cyclic = %v, want [<nil>]", sl)
	}

	ac.PongoPage(w, req, filename, pongodata, funcs)

	got := strings.TrimSpace(w.Body.String())
	const want = "numbers:1,2,3| objects:1,2,3| dicts:a=1,b=2| config:localhost:5432| matrix:3| cyclic:"
	if got != want {
		t.Errorf("rendered page = %q, want %q", got, want)
	}
}
