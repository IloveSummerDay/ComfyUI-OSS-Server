package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	region     string // 存储区域
	bucketName string // 存储空间名称
	objectName string // 对象名称
)

var logger *zap.Logger

func main() {
	// 初始化 zap 日志
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 初始化 Gin
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 加载默认配置并设置凭证提供者和区域
	// 创建OSS客户端
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	// 注册路由
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	})

	r.POST("/save-oss", func(c *gin.Context) {
		var req SaveOSSRequest

		// parse args
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "parse json args error",
				Error:   err.Error(),
			})
			return
		}

		// http get image files
		for _, file_name := range req.File_name_list {
			resp, err := http.Get("http://" + req.Ai_server_port + ":" + req.Ai_server_port + "?filename=" + file_name + "&type=output")
			if err != nil {
				c.JSON(http.StatusInternalServerError,
					ErrorResponse{
						Message: "Failed to fetch image from URL",
						Error:   err.Error(),
					})
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				c.JSON(resp.StatusCode, ErrorResponse{
					Message: "Failed to fetch image",
				})
				return
			}

			image_bytes, err := io.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Message: "Failed to read response body",
					Error:   err.Error(),
				})
				return
			}

			// 创建上传对象的请求
			objectName = strings.Replace(file_name, " ", "_", -1)
			body := bytes.NewReader(image_bytes)
			request := &oss.PutObjectRequest{
				Bucket: oss.Ptr(bucketName), // 存储空间名称
				Key:    oss.Ptr(objectName), // 对象名称
				Body:   body,                // 要上传的图片数据
			}

			// 发送上传对象的请求
			_, err = client.PutObject(context.TODO(), request)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Message: "failed to put object",
					Error:   err.Error(),
				})
			}
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
		})
	})

	// 启动服务
	logger.Info("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
