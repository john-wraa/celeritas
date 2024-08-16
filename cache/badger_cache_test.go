package cache

import (
	"testing"
)

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}
	inCache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("foo should not be in cache")
	}

	_ = testBadgerCache.Set("foo", "bar")
	inCache, err = testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("foo should be in cache")
	}
	err = testBadgerCache.Forget("foo")
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	x, err := testBadgerCache.Get("foo")
	if err != nil {
		t.Error(err)
	}
	if x != "bar" {
		t.Error("expected bar, got ", x)
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("foo", "foo")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("foo should not be in cache")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("alpha should not be in cache")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}
	err = testBadgerCache.Set("alpha2", "beta2")
	if err != nil {
		t.Error(err)
	}
	err = testBadgerCache.Set("beta", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("a")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("alpha should not be in cache")
	}

	inCache, err = testBadgerCache.Has("alpha2")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("alpha2 should not be in cache")
	}

	inCache, err = testBadgerCache.Has("beta")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("beta should be in cache")
	}

}
