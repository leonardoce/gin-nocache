package nocache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	nocache "github.com/leonardoce/gin-nocache"
	"github.com/stretchr/testify/require"
)

type noCacheHeaders struct {
	header string
	value  string
}

func TestNoCache(t *testing.T) {
	t.Parallel()

	epoch := time.Unix(0, 0).Format(time.RFC1123)
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "test", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.Header.Set("ETag", "test")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(nocache.NoCache())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"test": "test",
		})
	})

	router.ServeHTTP(w, r)

	for _, tst := range [...]noCacheHeaders{
		{
			header: "Expires",
			value:  epoch,
		},
		{
			header: "Cache-Control",
			value:  "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		},
		{
			header: "Pragma",
			value:  "no-cache",
		},
		{
			header: "X-Accel-Expires",
			value:  "0",
		},
	} {
		tst := tst

		t.Run(tst.header, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, w.Header().Get(tst.header), tst.value)
		})

		t.Run(tst.header, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, r.Header.Get("ETag"), "")
		})
	}
}
