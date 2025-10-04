package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lin-snow/ech0/internal/config"
)

// TrimURL 去除 URL 前后的空格和斜杠
func TrimURL(url string) string {
	if url == "" {
		return ""
	}

	// 去除连接地址前后的空格和斜杠
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "/")
	url = strings.TrimSuffix(url, "/")
	return url
}

// ExtractDomain 从 URL 中提取域名
func ExtractDomain(url string) string {
	// 去除协议部分
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	// 提取域名部分 (到第一个斜杠为止)
	slashIndex := strings.Index(url, "/")
	if slashIndex != -1 {
		url = url[:slashIndex]
	}

	return url
}

// Header 自定义请求头结构体
type Header struct {
	Header  string
	Content string
}

// SendRequest 发送 HTTP 请求
func SendRequest(url, method string, customHeader Header, timeout ...time.Duration) ([]byte, error) {
	// 默认超时时间，如果有传入参数则使用传入的
	clientTimeout := 2 * time.Second
	if len(timeout) > 0 {
		clientTimeout = timeout[0]
	}

	// 自定义 HTTP 客户端，忽略 TLS 证书验证
	client := &http.Client{
		Timeout: clientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略证书验证
			},
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加自定义请求头
	if customHeader.Header != "" && customHeader.Content != "" {
		req.Header.Set(customHeader.Header, customHeader.Content)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Println("Failed to close response body:", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

// GetMIMETypeFromFilenameOrURL 根据文件名或 URL 获取 MIME 类型
func GetMIMETypeFromFilenameOrURL(filenameOrURL string) string {
	lowerFilename := strings.ToLower(filenameOrURL)
	switch {
	case strings.HasSuffix(lowerFilename, ".jpg"), strings.HasSuffix(lowerFilename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lowerFilename, ".png"):
		return "image/png"
	case strings.HasSuffix(lowerFilename, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lowerFilename, ".bmp"):
		return "image/bmp"
	case strings.HasSuffix(lowerFilename, ".webp"):
		return "image/webp"
	case strings.HasSuffix(lowerFilename, ".mp4"):
		return "video/mp4"
	case strings.HasSuffix(lowerFilename, ".mov"):
		return "video/quicktime"
	case strings.HasSuffix(lowerFilename, ".mp3"):
		return "audio/mpeg"
	case strings.HasSuffix(lowerFilename, ".wav"):
		return "audio/wav"
	case strings.HasSuffix(lowerFilename, ".ogg"):
		return "audio/ogg"
	case strings.HasSuffix(lowerFilename, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lowerFilename, ".doc"):
		return "application/msword"
	case strings.HasSuffix(lowerFilename, ".docx"):
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case strings.HasSuffix(lowerFilename, ".xls"):
		return "application/vnd.ms-excel"
	case strings.HasSuffix(lowerFilename, ".xlsx"):
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case strings.HasSuffix(lowerFilename, ".ppt"):
		return "application/vnd.ms-powerpoint"
	case strings.HasSuffix(lowerFilename, ".pptx"):
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case strings.HasSuffix(lowerFilename, ".txt"):
		return "text/plain"
	case strings.HasSuffix(lowerFilename, ".html"), strings.HasSuffix(lowerFilename, ".htm"):
		return "text/html"
	case strings.HasSuffix(lowerFilename, ".csv"):
		return "text/csv"
	default:
		return "application/octet-stream" // 默认二进制流
	}
}

// PostActivity 发送 POST 请求
func PostActivity(activity []byte, inboxURL string, actorID string) error {
	priv := config.RSA_PRIVATE
	if priv == nil {
		return fmt.Errorf("private key is not initialized")
	}

	// 构造请求
	req, err := http.NewRequest("POST", inboxURL, strings.NewReader(string(activity)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/activity+json")

	// 调用 SignRequest 进行签名
	if err := SignRequest(req, priv, actorID+"#main-key", activity); err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	// 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-2xx response: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// SignRequest 对请求进行签名
func SignRequest(req *http.Request, priv *rsa.PrivateKey, keyID string, body []byte) error {
	// 1. Digest
	digest := sha256.Sum256(body)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])
	req.Header.Set("Digest", "SHA-256="+digestBase64)

	// 2. Date
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set("Date", date)

	// 3. Signature
	signingString := fmt.Sprintf("(request-target): %s %s\ndate: %s\ndigest: SHA-256=%s",
		strings.ToLower(req.Method),
		req.URL.RequestURI(),
		date,
		digestBase64,
	)

	// 4. 签名
	hashed := sha256.Sum256([]byte(signingString))
	sig, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}
	sigBase64 := base64.StdEncoding.EncodeToString(sig)

	// 5. 构建 Signature 头
	signatureHeader := fmt.Sprintf(`keyId="%s",algorithm="rsa-sha256",headers="(request-target) date digest",signature="%s"`,
		keyID,
		sigBase64,
	)
	req.Header.Set("Signature", signatureHeader)

	return nil
}
