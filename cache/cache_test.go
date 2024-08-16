package cache

import "testing"

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("Foo found in cache, and it shouldn't be there")
	}

	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("Foo not found in cache, but it should be there")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}
	x, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}
	if x != "bar" {
		t.Error("Expected bar, got ", x)
	}
}

func TestRedisCache_Forget(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}
	err = testRedisCache.Forget("alpha")
	if err != nil {
		t.Error(err)
	}
	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("Alpha found in cache, and it shouldn't be there")
	}
}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("alpha", "beta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("Alpha found in cache, and it shouldn't be there")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("alpha", "theta")
	if err != nil {
		t.Error(err)
	}
	err = testRedisCache.Set("alpha2", "theta")
	if err != nil {
		t.Error(err)
	}
	err = testRedisCache.Set("beta", "theta")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("alpha")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("alpha")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("Alpha found in cache, and it shouldn't be there")
	}

	inCache, err = testRedisCache.Has("alpha2")
	if err != nil {
		t.Error(err)
	}
	if inCache {
		t.Error("Alpha2 found in cache, and it shouldn't be there")
	}

	inCache, err = testRedisCache.Has("beta")
	if err != nil {
		t.Error(err)
	}
	if !inCache {
		t.Error("Beta not found in cache, and it should be there")
	}
}

func TestEncodeDecode(t *testing.T) {
	entry := Entry{}
	entry["foo"] = "bar"
	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}
	_, err = decode(string(bytes))
	if err != nil {
		t.Error(err)
	}
}
