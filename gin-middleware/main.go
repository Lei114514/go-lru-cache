package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	lru "github.com/Lei114514/go-lru-cache"
)

type CacheResponse struct{
	Status int
	Headers http.Header
	Body []byte
}

type cacheWriter struct{
	gin.ResponseWriter
	status int
	body *bytes.Buffer
}

func newCacheWriter(w gin.ResponseWriter) *cacheWriter{
	return &cacheWriter{
		ResponseWriter: w,
		status: http.StatusOK,
		body: &bytes.Buffer{},
	}
}

func (w *cacheWriter) WriteHeader(code int){
	w.status=code
	w.WriteHeader(code)
}

func (w *cacheWriter) Write(data []byte) (int, error){
	w.body.Write(data)
	return w.Write(data)
}

func CacheMiddleware(cache lru.Cache) gin.HandlerFunc {
	return func(c *gin.Context){
		if c.Request.Method !="GET"{
			c.Next()
			return 
		}

		key := c.Request.URL.String()
 
		if value, hit := cache.Get(key) ; hit{
			cached := value.(*CacheResponse)
			for k,v := range cached.Headers{
				for _,vv := range v{
					c.Writer.Header().Add(k,vv)
				}
			}
			c.Data(cached.Status,cached.Headers.Get("Content-Type"),cached.Body)
			c.Abort()
			return 
		}

		writer := newCacheWriter(c.Writer)
		c.Writer = writer 

		c.Next()

		if writer.status == http.StatusOK && writer.body.Len() > 0{
			cache.Put(key,&CacheResponse{
				Status: writer.status,
				Headers: c.Writer.Header().Clone(),
				Body: writer.body.Bytes(),
			})
		}
	}
}


func main(){
	r := gin.Default()

	// 創建並發安全緩存（容量100）
	cache := lru.NewConcurrentCache(100)

	// 使用緩存中間件
	r.Use(CacheMiddleware(cache))

	// 示例 API 1：模擬耗時查詢
	r.GET("/api/data/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 模擬耗時操作（如數據庫查詢）
		time.Sleep(100 * time.Millisecond)

		c.JSON(200, gin.H{
			"id":      id,
			"data":    "這是從數據庫查詢的數據",
			"time":    time.Now().Format("15:04:05.000"),
			"cached":  false,
		})
	})

	// 示例 API 2：獲取當前緩存狀態
	r.GET("/api/cache/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size": cache.Len(),
			"cap":  cache.Cap(),
		})
	})

	// 示例 API 3：清空緩存
	r.POST("/api/cache/clear", func(c *gin.Context) {
		cache.Clear()
		c.JSON(200, gin.H{"message": "cache cleared"})
	})

	// 示例 API 4：POST 請求（不會被緩存）
	r.POST("/api/data", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"message": "created",
			"name":    req.Name,
			"time":    time.Now().Format("15:04:05.000"),
		})
	})

	r.Run(":8081")
}