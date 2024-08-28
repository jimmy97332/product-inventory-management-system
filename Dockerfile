# Golang 基礎映像檔
from golang:1.23.0

# 設置工作目錄，WORKDIR 是 Dockerfile 的指令，用來設置在容器內部用於存儲應用程式檔案的目錄。
WORKDIR myapp

# 複製 依賴文件到工作目錄
COPY go.mod go.sum ./
# 安裝
RUN go mod download

# 複製應用程式程式碼到工作目錄
COPY . .
COPY .env .

# 建立
RUN go build -o myapp

# 暴露應用程式所使用的端口
EXPOSE 8080

# 容器內啟動應用程式
CMD ["./myapp"]