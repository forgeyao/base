package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装响应写入器来捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LogAPIRequest 记录 API 请求日志。
func LogAPIRequest(c *gin.Context, requestBody interface{}) {
	headers := make(map[string]string, len(c.Request.Header))
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	apiLog := map[string]interface{}{
		"headers":      headers,
		"method":       c.Request.Method,
		"path":         c.Request.URL.Path,
		"client_ip":    c.ClientIP(),
		"request_body": requestBody,
	}

	if data, err := marshalLog(apiLog); err == nil {
		namedLogger("api").Infof("Req: %s", data)
	} else {
		namedLogger("api").Warnf("marshal request log failed: %v", err)
	}
}

// LogAPIResponse 记录 API 响应日志。
func LogAPIResponse(c *gin.Context, response interface{}, startTime time.Time, err error) {
	apiLog := map[string]interface{}{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"client_ip":   c.ClientIP(),
		"response":    response,
		"status_code": c.Writer.Status(),
		"duration":    time.Since(startTime).String(),
	}
	if err != nil {
		apiLog["error"] = err.Error()
	}

	if data, marshalErr := marshalLog(apiLog); marshalErr == nil {
		namedLogger("api").Infof("Resp: %s", data)
	} else {
		namedLogger("api").Warnf("marshal response log failed: %v", marshalErr)
	}
}

func marshalLog(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(buf.Bytes())), nil
}

// GinLog 在 Gin 请求链路上读取请求体、捕获响应体并写入 API 日志。
func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var requestBody interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 清除 form 解析缓存，确保 ParseForm 能重新从 body 读取（针对 form 格式请求）
			c.Request.Form = nil
			c.Request.PostForm = nil

			// 尝试解析JSON请求体
			if len(bodyBytes) > 0 {
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					requestBody = jsonBody
				} else {
					requestBody = string(bodyBytes)
				}
			}
		}

		// 记录请求日志
		LogAPIRequest(c, requestBody)

		// 包装响应写入器
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 解析响应体
		var responseBody interface{}
		if w.body.Len() > 0 {
			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(w.body.Bytes(), &jsonResponse); err == nil {
				responseBody = jsonResponse
			} else {
				responseBody = w.body.String()
			}
		}

		// 记录响应日志
		var err error
		if len(c.Errors) > 0 {
			err = c.Errors.Last()
		}
		LogAPIResponse(c, responseBody, startTime, err)
	}
}
